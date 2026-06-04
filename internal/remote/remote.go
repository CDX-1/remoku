package remote

import (
	"fmt"
	"os"
	"github.com/CDX-1/remoku/internal/discovery"
	"github.com/CDX-1/remoku/internal/ecp"
	"github.com/CDX-1/remoku/internal/macro"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	rokuIP  string
	timeout time.Duration
)

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "remoku",
		Short: "Remoku is a CLI tool for controling Roku TVs",
		Long:  "A local network utility tool for controlling Roku TVs over HTTP.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Name() == "scan" {
				return
			}

			// automatically scan for IP if on auto-scan mode
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
		},
	}

	rootCmd.PersistentFlags().StringVarP(&rokuIP, "ip", "i", "auto", "IP address of your Roku TV")
	rootCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 2*time.Second, "Network request timeout limit")

	var scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "Scan the local network for Roku devices",
		Long:  "Broadcast an SSDP M-SEARCH packet over UDP to find the IP addresses of any Roku devices.",
		Run: func(cmd *cobra.Command, args []string) {
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
		Use:   "press",
		Short: "Send a keypress to the TV",
		Long:  "Simulate pressing a physical button on a remote controller.",
		Run: func(cmd *cobra.Command, args []string) {
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

	var interactiveCmd = &cobra.Command{
		Use:   "interactive",
		Short: "Start an interactive remote session",
		Long:  "Connects to the Roku TV and allows you to control it using your keyboard in real-time.",
		Run: func(cmd *cobra.Command, args []string) {
			err := RunInteractiveMode(rokuIP, timeout)
			if err != nil {
				fmt.Printf("❌ Interactive mode failed: %v\n", err)
				os.Exit(1)
			}
		},
	}

	var appsCmd = &cobra.Command{
		Use:   "apps",
		Short: "Lists installed applications",
		Long:  "Retrieves and displays a list of all installed applications on the Roku TV.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("📋 Fetching installed applications...")

			apps, err := ecp.GetApps(rokuIP, timeout)
			if err != nil {
				fmt.Printf("❌ Failed to fetch apps: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("✅ Apps list retrieved successfully!")
			fmt.Printf("=========================\n")
			
			for _, app := range apps {
				fmt.Printf("App: %s\n", app.Name)
				fmt.Printf("ID: %s\n", app.ID)
				fmt.Printf("Type: %s\n", app.Type)
				fmt.Printf("Version: %s\n", app.Version)
				fmt.Printf("\n")
			}

			fmt.Printf("=========================\n")

			fmt.Println("🚀 Launch any of these apps with 'remoku launch <appId>'")
		},
	}

	var launchCmd = &cobra.Command{
		Use:   "launch",
		Short: "Launch an application",
		Long:  "Launches the specified application on the Roku TV.",
		Run: func(cmd *cobra.Command, args []string) {
			appId := args[0]
			fmt.Printf("🚀 Launching application: %s...\n", appId)

			err := ecp.LaunchApp(rokuIP, timeout, appId)
			if err != nil {
				fmt.Printf("❌ Failed to launch app: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("✅ App launched successfully!")
		},
	}

	var macroCmd = &cobra.Command{
		Use:	"macro",
		Short:	"Execute a macro JSON file",
		Long:	"Executes a specified JSON file as a macro",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("❌ No macro file specified")
				os.Exit(1)
			}
			filePath := args[0]
			
			actions, err := macro.LoadMacroFromFile(filePath)
			if err != nil {
				fmt.Printf("❌ Failed to load macro: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("⌛ Executing macro...")
			err = macro.ExecuteActions(rokuIP, timeout, actions, 100 * time.Millisecond)
			if err != nil {
				fmt.Printf("❌ Failed to execute macro: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("✅ Macro executed successfully!")
		},
	}

	rootCmd.AddCommand(scanCmd, pressCmd, interactiveCmd, appsCmd, launchCmd, macroCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
