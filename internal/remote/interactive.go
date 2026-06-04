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
	fmt.Println("          J = Volume Up, K = Volume Down, M = Mute\r")
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
			rokuKey = "up"
		case 'a', 'A':
			rokuKey = "left"
		case 's', 'S':
			rokuKey = "down"
		case 'd', 'D':
			rokuKey = "right"
		case 'j', 'J':
			rokuKey = "vup"
		case 'k', 'K':
			rokuKey = "vdown"
		case 'm', 'M':
			rokuKey = "mute"
		case 13, 10: // Enter
			rokuKey = "select"
		case 8, 127: // Backspace
			rokuKey = "back"
		case ' ':
			rokuKey = "home"
		case 'q', 'Q', 3:
			fmt.Printf("\n🚪 Exiting interactive mode.\r\n")
			return nil
		default:
			continue
		}

		fmt.Printf("\r📦 Sending keypress: %-10s", rokuKey)

		start := time.Now()
		err = ecp.SendKeyPress(rokuIP, timeout, rokuKey)
		duration := time.Since(start)
		if err != nil {
			fmt.Printf("❌ failed to send keypress: %v\r\n", err)
		}
		fmt.Printf("\r🟢 Sent keypress: %-10s (duration: %s)\r\n", rokuKey, duration.Truncate(time.Millisecond))
	}

	return nil
}
