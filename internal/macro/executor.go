package macro

import (
	"fmt"
	"remoku/internal/ecp"
	"time"
)

type Action struct {
	Type     string `json:"type"`
	Value    string `json:"value,omitempty"`
	Duration int    `json:"duration,omitempty"`
}

func ExecuteActions(rokuIP string, timeout time.Duration, actions []Action, delay time.Duration) error {
	for _, action := range actions {
		switch action.Type {
		case "sleep":
			time.Sleep(time.Duration(action.Duration) * time.Millisecond)
			continue // explicitly skip global delay

		case "keypress":
			err := ecp.SendKeyPress(rokuIP, timeout, action.Value)
			if err != nil {
				return fmt.Errorf("keypress failed, key: %s, error: %w", action.Value, err)
			}

		case "launch":
			err := ecp.LaunchApp(rokuIP, timeout, action.Value)
			if err != nil {
				return fmt.Errorf("launch failed, app: %s, error: %w", action.Value, err)
			}

		default:
			return fmt.Errorf("unknown action %s", action.Type)
		}

		time.Sleep(delay)
	}

	return nil
}
