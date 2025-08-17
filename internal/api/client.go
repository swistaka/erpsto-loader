package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type APIClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewAPIClient(baseURL, apiKey string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *APIClient) DownloadToTempFile() (string, error) {
	data, err := c.getEntities()
	if err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "erpsto-api-*.json")
	if err != nil {
		return "", fmt.Errorf("temp file creation failed: %w", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(data); err != nil {
		return "", fmt.Errorf("temp file write failed: %w", err)
	}

	return tmpFile.Name(), nil
}

func (c *APIClient) getEntities() ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
