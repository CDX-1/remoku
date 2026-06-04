package remote

import (
	"fmt"
	"os"
	"remoku/internal/ecp"
	"time"

	"golang.org/x/term"
)

func RunInteractiveMode(rokuIP string, timeout time.Duration) (error) {
	fmt.Printf("🔗 Connected to Roku TV at %s\r\n", rokuIP)
	fmt.Println("🚀 Real-time control active! (Press 'q' or Ctrl+C to exit)\r")
	fmt.Println("Commands: W = Up, A = Left, S = Down, D = Right\r")
	fmt.Println("          [Enter] = Select, [Backspace] = Back, [Spacebar] = Home\r")
	fmt.Println("====================================================\r")

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("error setting terminal to raw mode: %v\r\n", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	buf := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}

		char := buf[0]
		var rokuKey string

		switch char {
		case 'w', 'W':
			rokuKey = "Up"
		case 'a', 'A':
			rokuKey = "Left"
		case 's', 'S':
			rokuKey = "Down"
		case 'd', 'D':
			rokuKey = "Right"
		case 13, 10: // Enter
			rokuKey = "Select"
		case 8, 127: // Backspace
			rokuKey = "Back"
		case ' ':
			rokuKey = "Home"
		case 'q', 'Q', 3:
			fmt.Printf("\n🚪 Exiting interactive mode.\r\n")
			return nil
		default:
			continue
		}

		fmt.Printf("\r📦 Sending keypress: %-10s", rokuKey)
		err = ecp.PostECP(rokuIP, timeout, "keypress/" + rokuKey)
		if err != nil {
			fmt.Printf("\nfailed to send keypress: %v\r\n", err)
		}
	}

	return nil
}
