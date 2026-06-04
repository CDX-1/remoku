package ecp

import (
	"encoding/xml"
	"fmt"
	"time"
)

type Apps struct {
	XMLName xml.Name `xml:"apps"`
	List    []App    `xml:"app"`
}

type App struct {
	XMLName xml.Name `xml:"app"`
	ID      string   `xml:"id,attr"`
	Type    string   `xml:"type,attr"`
	Version string   `xml:"version,attr"`
	Name    string   `xml:",chardata"`
}

func GetApps(rokuIP string, timeout time.Duration) ([]App, error) {
	body, err := GetECP(rokuIP, timeout, "query/apps")
	if err != nil {
		return nil, fmt.Errorf("failed to get apps: %w", err)
	}
	var apps Apps
	err = xml.Unmarshal([]byte(body), &apps)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal apps: %w", err)
	}
	return apps.List, nil
}