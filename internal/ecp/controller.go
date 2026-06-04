package ecp

import (
	"fmt"
	"time"
)

var keyMap = map[string]string{
	"home":   "Home",
	"back":   "Back",
	"up":     "Up",
	"down":   "Down",
	"left":   "Left",
	"right":  "Right",
	"ok":     "Select",
	"select": "Select",
	"play":   "Play",
	"pause":  "Pause",
	"vup":    "VolumeUp",
	"vdown":  "VolumeDown",
	"mute":   "VolumeMute",
	"power":  "PowerOff",
	"off":    "PowerOff",
	"on":     "PowerOn",
}

func SendKeyPress(rokuIP string, timeout time.Duration, key string) error {
	mappedKey, ok := keyMap[key]
	if !ok {
		return fmt.Errorf("unknown key: %s", key)
	}

	return PostECP(rokuIP, timeout, "keypress/"+mappedKey)
}

func LaunchApp(rokuIP string, timeout time.Duration, appId string) error {
	err := PostECP(rokuIP, timeout, "launch/"+appId)
	if err != nil {
		fmt.Errorf("failed to launch app: %v\r\n", err)
	}
	return nil
}