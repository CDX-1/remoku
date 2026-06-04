package remote

import (
	"fmt"
	"os"
	"remoku/internal/discovery"

	"github.com/spf13/cobra"
)

func Execute() {

	var rootCmd = &cobra.Command{
		Use: "remoku",
		Short: "Remoku is a CLI tool for controling Roku TVs",
		Long: "A local network utility tool for controlling Roku TVs over HTTP",
	}

	var scanCmd = &cobra.Command{
		Use:	"scan",
		Short:	"Scan the local network for Roku devices",
		Long:	"Broadcast an SSDP M-SEARCH packet over UDP to find the IP addresses of any Roku devices",
		Run:	func(cmd *cobra.Command, args []string) {
			fmt.Println("Broadcasting SSDP discovery packet over local subnet...")

			foundIP, err := discovery.ScanForRoku()
			if err != nil {
				fmt.Printf("Scan failed: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("\nDiscovery Successful")
			fmt.Printf("=========================\n")
			fmt.Printf("Device Type:	Roku TV / ECP Target\n")
			fmt.Printf("IP Address:		%s\n", foundIP)
			fmt.Printf("Target URL:		http://%s:8060\n", foundIP)
			fmt.Printf("=========================\n")
		},
	}

	rootCmd.AddCommand(scanCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	
}