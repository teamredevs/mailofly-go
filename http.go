package mailofly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultBaseURL = "https://api.mailofly.com"
	APIPrefix      = "/api/v1"
)

func normalizeBaseURL(base string) string {
	if base == "" {
		base = DefaultBaseURL
	}
	return strings.TrimRight(base, "/")
}

func doRequest(httpClient *http.Client, baseURL, path, method, apiKey string, body any, query map[string]string) (any, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if len(query) > 0 {
		q := url.Values{}
		for k, v := range query {
			if v != "" {
				q.Set(k, v)
			}
		}
		if enc := q.Encode(); enc != "" {
			if strings.Contains(path, "?") {
				path += "&" + enc
			} else {
				path += "?" + enc
			}
		}
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, baseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Mailofly-Client", "sdk/go")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var parsed any
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &parsed); err != nil {
			parsed = string(raw)
		}
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		errCode := res.Status
		var detail string
		if m, ok := parsed.(map[string]any); ok {
			if e, ok := m["error"].(string); ok {
				errCode = e
			}
			if msg, ok := m["message"].(string); ok {
				detail = msg
			}
		} else if s, ok := parsed.(string); ok {
			detail = s
		}
		return nil, &Error{Status: res.StatusCode, Err: errCode, DetailMessage: detail, Body: parsed}
	}

	return parsed, nil
}

func enc(s string) string {
	return url.PathEscape(s)
}

func mustAPIKey(apiKey string) (string, error) {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return "", fmt.Errorf("mailofly: apiKey is required")
	}
	return apiKey, nil
}
