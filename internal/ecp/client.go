package ecp

import (
	"fmt"
	"net/http"
	"time"
)

func PostECP(rokuIP string, timeout time.Duration, endpoint string) error {
	url := fmt.Sprintf("http://%s:8060/%s", rokuIP, endpoint)

	client := &http.Client{Timeout: timeout}

	res, err := client.Post(url, "text/plain", nil)
	if err != nil {
		return fmt.Errorf("failed to send ECP command: %w", err)
	}
	defer res.Body.Close()
	
	if res.StatusCode == http.StatusOK {
		return nil
	} else if res.StatusCode == http.StatusForbidden {
		return fmt.Errorf("failed to send ECP command: 403 Forbidden - ensure 'Control by mobile apps' is enabled in Settings > System > Advanced system settings")
	} else {
		return fmt.Errorf("failed to send ECP command: unexpected status code from TV (%d)", res.StatusCode)
	}
}