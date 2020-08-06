package sqlpad

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Client ...
type Client struct {
	httpClient *http.Client
	baseURL    string
	username   string
	password   string
}

// NewClient ...
func NewClient(baseURL, username, password string, httpClient *http.Client) *Client {
	if strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL[:len(baseURL)-1]
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		username:   username,
		password:   password,
	}
}

func (c *Client) request(ctx context.Context, method, urlMethod string, payload interface{}) (*http.Request, error) {
	if !strings.HasPrefix(urlMethod, "/") {
		urlMethod = "/" + urlMethod
	}
	u, err := url.Parse(c.baseURL + urlMethod)
	if err != nil {
		return nil, fmt.Errorf("[sqlpad_client] failed to parse urlMethod: %w", err)
	}

	var buf io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("[sqlpad_client] failed to marshal payload: %w", err)
		}
		buf = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("[sqlpad_client] failed to create new http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.username, c.password)

	return req.WithContext(ctx), nil
}

func (c *Client) do(req *http.Request, resPtr interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("[sqlpad_client] failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		var body []byte
		if body, err = ioutil.ReadAll(resp.Body); err != nil {
			return fmt.Errorf("[sqlpad_client] failed to read error response body, code: %d, %w", resp.StatusCode, err)
		}
		return &Error{Response: resp, StatusCode: resp.StatusCode, Body: body}
	}

	if resPtr == nil {
		return nil
	}

	if err = json.NewDecoder(resp.Body).Decode(resPtr); err != nil {
		return fmt.Errorf("[sqlpad_client] failed to unmarshal response: %w", err)
	}

	return nil
}

func (c *Client) doRequest(ctx context.Context, method, urlMethod string, payload interface{}, res interface{}) error {
	req, err := c.request(ctx, method, urlMethod, payload)
	if err != nil {
		return err
	}
	return c.do(req, res)
}
