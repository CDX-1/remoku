package discovery

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	ssdpAddr = "239.255.255.250:1900"
	scanTime = 5 * time.Second
)

func ScanForRoku() (string, error) {
	localIP, err := getPhysicalLocalIP()
	if err != nil {
		return "", fmt.Errorf("could not find an active local network interface: %w", err)
	}

	udpRemoteAddr, err := net.ResolveUDPAddr("udp", ssdpAddr)
	if err != nil {
		return "", fmt.Errorf("failed to resolve SSDP address: %w", err)
	}

	udpLocalAddr, err := net.ResolveUDPAddr("udp", localIP+":0")
	if err != nil {
		return "", fmt.Errorf("failed to resolve local bind address: %w", err)
	}

	conn, err := net.ListenUDP("udp", udpLocalAddr)
	if err != nil {
		return "", fmt.Errorf("failed to listen on active interface: %w", err)
	}
	defer conn.Close()

	payload := strings.Join([]string{
		"M-SEARCH * HTTP/1.1",
		"HOST: " + ssdpAddr,
		"MAN: \"ssdp:discover\"",
		"ST: roku:ecp",
		"MX: 2",
		"",
		"",
	}, "\r\n")

	_, err = conn.WriteToUDP([]byte(payload), udpRemoteAddr)
	if err != nil {
		return "", fmt.Errorf("failed to broadcast SSDP payload: %w", err)
	}

	conn.SetReadDeadline(time.Now().Add(scanTime))

	buf := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			return "", fmt.Errorf("no Roku TV discovered on your network: %w", err)
		}

		response := string(buf[:n])
		if strings.Contains(response, "roku:ecp") {
			return extractIP(response)
		}
	}
}

func getPhysicalLocalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve network interfaces: %w", err)
	}

	for _, iface := range ifaces {
		// filter out down or loopback network interfaces
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// filter out common virtual network adapters
		name := strings.ToLower(iface.Name)
		if strings.Contains(name, "vbox") || strings.Contains(name, "docker") || strings.Contains(name, "wsl") || strings.Contains(name, "virtual") {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// filter out non IPv4 subnets
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}

			ipStr := ip.String()
			if strings.HasPrefix(ipStr, "10.") || strings.HasPrefix(ipStr, "192.168.") || strings.HasPrefix(ipStr, "172.") {
				return ipStr, nil
			}
		}
	}

	return "", fmt.Errorf("no active physical IPv4 address detected")
}

func extractIP(response string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(response))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.ToUpper(line), "LOCATION:") {
			parts := strings.Split(line, "//")
			if len(parts) > 1 {
				ipAndPort := strings.Split(parts[1], ":")
				if len(ipAndPort) > 0 {
					return ipAndPort[0], nil
				}
			}
		}
	}

	return "", fmt.Errorf("found a Roku signature, but location parsing failed")
}