package remote

import (
	"fmt"
	"os"
	"remoku/internal/discovery"
	"remoku/internal/ecp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	rokuIP	string
	timeout	time.Duration
)

func Execute() {

	var rootCmd = &cobra.Command{
		Use: "remoku",
		Short: "Remoku is a CLI tool for controling Roku TVs",
		Long: "A local network utility tool for controlling Roku TVs over HTTP.",
	}

	rootCmd.PersistentFlags().StringVarP(&rokuIP, "ip", "i", "auto", "IP address of your Roku TV")
	rootCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 2*time.Second, "Network request timeout limit")

	var scanCmd = &cobra.Command{
		Use:	"scan",
		Short:	"Scan the local network for Roku devices",
		Long:	"Broadcast an SSDP M-SEARCH packet over UDP to find the IP addresses of any Roku devices.",
		Run:	func(cmd *cobra.Command, args []string) {
			fmt.Println("📡 Broadcasting SSDP discovery packet over local subnet...")

			foundIP, err := discovery.ScanForRoku()
			if err != nil {
				fmt.Printf("❌ Scan failed: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("\n✅ Discovery Successful")
			fmt.Printf("=========================\n")
			fmt.Printf("Roku TV / ECP Target\n")
			fmt.Printf("🌐 IP Address:	%s\n", foundIP)
			fmt.Printf("🔗 Target URL:	http://%s:8060\n", foundIP)
			fmt.Printf("💡 Append \"--ip %s\" to skip auto-scan\n", foundIP)
			fmt.Printf("=========================\n")
		},
	}

	var pressCmd = &cobra.Command{
		Use:	"press",
		Short:	"Send a keypress to the TV",
		Long:	"Simulate pressing a physical button on a remote controller.",
		Run:	func(cmd *cobra.Command, args []string) {
			// no ip provided, auto-scan
			if rokuIP == "auto" {
				fmt.Println("🔍 No IP provided. Auto-scanning network...")
				foundIP, err := discovery.ScanForRoku()
				if err != nil {
					fmt.Printf("❌ Discovery Failed: %v\n", err)
					os.Exit(1)
				}
				rokuIP = foundIP
				fmt.Printf("🎯 Using Roku TV at %s\n", rokuIP)
			}

			input := strings.ToLower(args[0])
			fmt.Printf("👆 Sending \"%s\" keypress...\n", input)

			err := ecp.SendKeyPress(rokuIP, timeout, input)
			if err != nil {
				fmt.Printf("❌ Keypress failed: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("✅ Keypress sent successfully!")
		},
	}

	rootCmd.AddCommand(scanCmd, pressCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	
}