package googleapi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// Config structure for storing API details
type Config struct {
	CX     string `json:"cx"`
	APIKey string `json:"api_key"`
}

// Response structure for Google Custom Search
type GoogleSearchResponse struct {
	Items []struct {
		Link string `json:"link"`
	} `json:"items"`
	SearchInformation struct {
		TotalResults string `json:"totalResults"`
	} `json:"searchInformation"`
}

// GoogleSearch queries Google Programmable Search API and returns URLs
func GoogleSearch(config Config, query string, maxResults int) ([]string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	var urls []string
	start := 1
	for start < maxResults {
		endpoint := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s&start=%d",
			config.APIKey,
			config.CX,
			url.QueryEscape(query),
			start,
		)

		resp, err := client.Get(endpoint)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 429 || resp.StatusCode == 403 {
			return nil, fmt.Errorf("API limit reached or forbidden: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var searchResponse GoogleSearchResponse
		if err := json.Unmarshal(body, &searchResponse); err != nil {
			return nil, err
		}

		for _, item := range searchResponse.Items {
			urls = append(urls, item.Link)
		}

		totalResults, err := strconv.Atoi(searchResponse.SearchInformation.TotalResults)
		if err != nil || totalResults <= start {
			break
		}

		start += 10
	}

	return urls, nil
}
