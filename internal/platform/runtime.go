package platform

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// RuntimeStatus is the bounded response exposed by runtime inspection.
type RuntimeStatus struct {
	Status string `json:"status"`
}

// RuntimeClient reads documented runtime health endpoints.
type RuntimeClient interface {
	Ready() (RuntimeStatus, error)
}

type HTTPRuntimeClient struct {
	BaseURL string
	Client  *http.Client
}

func (c HTTPRuntimeClient) Ready() (RuntimeStatus, error) {
	client := c.Client
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	url := strings.TrimRight(c.BaseURL, "/") + "/health/ready"
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return RuntimeStatus{}, fmt.Errorf("runtime request: %w", err)
	}
	response, err := client.Do(request)
	if err != nil {
		return RuntimeStatus{}, fmt.Errorf("runtime unavailable")
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return RuntimeStatus{}, fmt.Errorf("runtime returned status %d", response.StatusCode)
	}
	var status RuntimeStatus
	if err := json.NewDecoder(response.Body).Decode(&status); err != nil {
		return RuntimeStatus{}, fmt.Errorf("runtime returned invalid status")
	}
	return status, nil
}
