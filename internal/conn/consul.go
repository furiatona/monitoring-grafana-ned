package consul

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	dcpatterns "ned/internal/discovery"
)

var (
	ConsulAddress = defaultConsulAddress
	ConsulToken   = defaultConsulToken
	ConsulKVKey   = defaultConsulKVKey
)

// These will be overridden at build time
var (
	defaultConsulAddress string
	defaultConsulToken   string
	defaultConsulKVKey   string
)

func FetchDCPatterns() (map[string]dcpatterns.DCPattern, error) {
	if ConsulAddress == "" || ConsulKVKey == "" {
		return nil, fmt.Errorf("ConsulAddress and ConsulKVKey must be set")
	}

	url := fmt.Sprintf("%s/v1/kv/%s?raw", ConsulAddress, ConsulKVKey)
	client := &http.Client{Timeout: 10 * time.Second}

	const maxRetries = 3
	const retryDelay = 3 * time.Second

	var body []byte

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		if ConsulToken != "" {
			req.Header.Set("X-Consul-Token", ConsulToken)
		}

		resp, err := client.Do(req)
		if err != nil {
			if attempt == maxRetries {
				return nil, fmt.Errorf("failed to fetch from Consul after %d attempts: %v", maxRetries, err)
			}
			fmt.Printf("\033[1;33mAttempt %d failed. Retrying in %v...\033[0m\n", attempt, retryDelay)
			time.Sleep(retryDelay)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			if attempt == maxRetries {
				return nil, fmt.Errorf("consul returned status: %s", resp.Status)
			}
			fmt.Printf("\033[1;33mAttempt %d failed with status %s. Retrying in %v...\033[0m\n", attempt, resp.Status, retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		break
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("no data received from Consul after %d attempts", maxRetries)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse Consul response: %v", err)
	}

	patterns := make(map[string]dcpatterns.DCPattern)
	for k, v := range raw {
		val, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		key, _ := val["key"].(string)
		label, _ := val["label"].(string)
		patterns[k] = dcpatterns.DCPattern{Key: key, Label: label}
	}

	if len(patterns) == 0 {
		return nil, fmt.Errorf("no valid DC patterns found in Consul data")
	}

	return patterns, nil
}
