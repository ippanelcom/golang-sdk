package ippanel

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultBaseURL = "https://edge.ippanel.com/v1/api"

type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(apiKey string, optionalBaseURL ...string) *Client {
	baseURL := defaultBaseURL
	if len(optionalBaseURL) > 0 && optionalBaseURL[0] != "" {
		baseURL = optionalBaseURL[0]
	}
	return &Client{
		APIKey:     apiKey,
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// WebserviceRequest --- Send Webservice ---
type webserviceRequest struct {
	Sender  string                 `json:"from_number"`
	Message string                 `json:"message"`
	Type    string                 `json:"sending_type"` // e.g. "sms"
	Params  map[string]interface{} `json:"params"`       // Optional parameters for the message
}

type SendResponse struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

// PatternRequest --- Send Pattern ---
type patternRequest struct {
	Sender      string                 `json:"from_number"`
	Recipient   []string               `json:"recipients"`
	PatternCode string                 `json:"code"`
	Params      map[string]interface{} `json:"params"`
	Type        string                 `json:"sending_type"`
}

// VOTPRequest --- Send VOTP ---
type vOTPRequest struct {
	Code   string                 `json:"message"`
	Type   string                 `json:"sending_type"` // e.g. "sms"
	Params map[string]interface{} `json:"params"`       // in seconds
}

// Meta --- Meta info returned from API ---
type Meta struct {
	Status            bool          `json:"status"`
	Message           string        `json:"message"`
	MessageParameters []interface{} `json:"message_parameters"`
	MessageCode       string        `json:"message_code"`
}

func (c *Client) SendWebservice(message string, sender string, recipients []string) (*SendResponse, error) {
	payload := webserviceRequest{
		Sender: sender,
		Params: map[string]interface{}{
			"recipients": recipients,
		},
		Message: message,
		Type:    "webservice",
	}
	return c.post("/send", payload)
}

func (c *Client) SendPattern(patternCode string, sender string, recipient string, params map[string]interface{}) (*SendResponse, error) {
	payload := patternRequest{
		Sender:      sender,
		Recipient:   []string{recipient},
		PatternCode: patternCode,
		Params:      params,
		Type:        "pattern",
	}
	return c.post("/send", payload)
}

func (c *Client) SendVOTP(code int32, recipient string) (*SendResponse, error) {
	payload := vOTPRequest{
		Code: fmt.Sprintf("%d", code),
		Type: "votp",
		Params: map[string]interface{}{
			"recipients": []string{recipient},
		},
	}
	return c.post("/send", payload)
}

func (c *Client) post(path string, payload interface{}) (*SendResponse, error) {
	url := c.BaseURL + path

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}(resp.Body)

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(string(resBody))
	}

	var sendResp SendResponse
	if err := json.Unmarshal(resBody, &sendResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &sendResp, nil
}
