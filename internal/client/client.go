package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/chiragms-cdcx/coindcx-cs-mcp/internal/logger"
)

const (
	RequestTimeout      = 25 * time.Second
	DialTimeout          = 10 * time.Second
	IdleConnTimeout      = 90 * time.Second
	MaxIdleConnsPerHost  = 32
)

const optionsPathPrefix = "/api/v1/options"
const futuresPathPrefix = "/api/v1/derivatives/futures"
const bffPathPrefix = "/api/v1"
const apiV2Prefix = "/api/v2"

type Client struct {
	baseURL     string
	publicURL   string
	adminURL    string
	adminKey    string
	adminSecret string
	authToken   string
	httpClient  *http.Client
}

var defaultTransport = &http.Transport{
	DialContext:           (&net.Dialer{Timeout: DialTimeout}).DialContext,
	MaxIdleConnsPerHost:   MaxIdleConnsPerHost,
	IdleConnTimeout:       IdleConnTimeout,
	ResponseHeaderTimeout: 15 * time.Second,
}

func New(baseURL, publicURL, adminURL, adminKey, adminSecret, authToken string) *Client {
	return NewWithHTTPClient(baseURL, publicURL, adminURL, adminKey, adminSecret, authToken, nil)
}

func NewWithHTTPClient(baseURL, publicURL, adminURL, adminKey, adminSecret, authToken string, hc *http.Client) *Client {
	if baseURL == "" {
		baseURL = "https://api.coindcx.com"
	}
	if publicURL == "" {
		publicURL = "https://public.coindcx.com"
	}
	if adminURL == "" {
		adminURL = "https://admin-api.coindcx.com"
	}
	if hc == nil {
		hc = &http.Client{
			Transport: defaultTransport,
			Timeout:   RequestTimeout,
		}
	}
	return &Client{
		baseURL:     baseURL,
		publicURL:   publicURL,
		adminURL:    adminURL,
		adminKey:    adminKey,
		adminSecret: adminSecret,
		authToken:   authToken,
		httpClient:  hc,
	}
}

func (c *Client) GetPublic(path string, query url.Values) ([]byte, int, error) {
	return c.get(c.publicURL+path, query)
}

func (c *Client) GetBase(path string, query url.Values) ([]byte, int, error) {
	return c.get(c.baseURL+path, query)
}

func (c *Client) get(fullURL string, query url.Values) ([]byte, int, error) {
	if len(query) > 0 {
		fullURL = fullURL + "?" + query.Encode()
	}
	logger.Debug("HTTP GET %s", fullURL)
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		logger.Error("new request failed: %v", err)
		return nil, 0, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.Error("GET %s failed: %v", fullURL, err)
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("read body GET %s: %v", fullURL, err)
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode >= 400 {
		logger.Error("GET %s status=%d body:\n%s", fullURL, resp.StatusCode, logger.PrettyJSON(string(data)))
	}
	return data, resp.StatusCode, nil
}

func (c *Client) HasCredentials() bool {
	return c.adminKey != "" && c.adminSecret != ""
}

func (c *Client) HasAuthToken() bool {
	return c.authToken != ""
}

func (c *Client) PostSigned(path string, body any) ([]byte, int, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, 0, err
	}
	sig, err := SignPayload(c.adminSecret, payload)
	if err != nil {
		return nil, 0, err
	}
	targetURL := c.baseURL + path
	logger.Debug("HTTP POST (signed) %s", targetURL)
	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(payload))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-AUTH-APIKEY", c.adminKey)
	req.Header.Set("X-AUTH-SIGNATURE", sig)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode >= 400 {
		logger.Error("POST %s status=%d body:\n%s", targetURL, resp.StatusCode, logger.PrettyJSON(string(data)))
	}
	return data, resp.StatusCode, nil
}

func (c *Client) GetSigned(path string, query url.Values, body any) ([]byte, int, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, 0, err
	}
	sig, err := SignPayload(c.adminSecret, payload)
	if err != nil {
		return nil, 0, err
	}
	fullURL := c.baseURL + path
	if len(query) > 0 {
		fullURL = fullURL + "?" + query.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("X-AUTH-APIKEY", c.adminKey)
	req.Header.Set("X-AUTH-SIGNATURE", sig)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode >= 400 {
		logger.Error("GET %s status=%d body:\n%s", fullURL, resp.StatusCode, logger.PrettyJSON(string(data)))
	}
	return data, resp.StatusCode, nil
}

func (c *Client) OptionsGetPublic(path string, query url.Values) ([]byte, int, error) {
	targetURL := c.publicURL + optionsPathPrefix + "/" + path
	return c.get(targetURL, query)
}

func (c *Client) OptionsGetPrivate(path string, query url.Values) ([]byte, int, error) {
	targetURL := c.baseURL + optionsPathPrefix + "/" + path
	if len(query) > 0 {
		targetURL = targetURL + "?" + query.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode >= 400 {
		logger.Error("options GET %s status=%d", targetURL, resp.StatusCode)
	}
	return data, resp.StatusCode, nil
}

func (c *Client) OptionsPostPrivate(path string, body any) ([]byte, int, error) {
	targetURL := c.baseURL + optionsPathPrefix + "/" + path
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(payload))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.authToken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode >= 400 {
		logger.Error("options POST %s status=%d", targetURL, resp.StatusCode)
	}
	return data, resp.StatusCode, nil
}

func (c *Client) FuturesGetPrivate(path string, query url.Values) ([]byte, int, error) {
	urlStr := c.baseURL + futuresPathPrefix + "/" + path
	if len(query) > 0 {
		urlStr = urlStr + "?" + query.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode >= 400 {
		logger.Error("futures GET %s status=%d", urlStr, resp.StatusCode)
	}
	return data, resp.StatusCode, nil
}

func (c *Client) FuturesPostPrivate(path string, body any) ([]byte, int, error) {
	urlStr := c.baseURL + futuresPathPrefix + "/" + path
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewReader(payload))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.authToken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode >= 400 {
		logger.Error("futures POST %s status=%d", urlStr, resp.StatusCode)
	}
	return data, resp.StatusCode, nil
}

func (c *Client) SpotGetV2Private(path string, query url.Values) ([]byte, int, error) {
	urlStr := c.baseURL + apiV2Prefix + "/" + path
	if len(query) > 0 {
		urlStr = urlStr + "?" + query.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode >= 400 {
		logger.Error("spot v2 GET %s status=%d", urlStr, resp.StatusCode)
	}
	return data, resp.StatusCode, nil
}

func (c *Client) AdminGet(path string, query url.Values) ([]byte, int, error) {
	urlStr := c.adminURL + path
	if len(query) > 0 {
		urlStr = urlStr + "?" + query.Encode()
	}
	logger.Debug("HTTP GET (admin) %s", urlStr)
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("X-ADMIN-SERVICE-KEY", c.adminKey)
	req.Header.Set("X-ADMIN-SERVICE-SECRET", c.adminSecret)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.Error("admin GET %s failed: %v", urlStr, err)
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode >= 400 {
		logger.Error("admin GET %s status=%d", urlStr, resp.StatusCode)
	}
	return data, resp.StatusCode, nil
}

func (c *Client) AdminPost(path string, body any) ([]byte, int, error) {
	urlStr := c.adminURL + path
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, 0, err
	}
	logger.Debug("HTTP POST (admin) %s", urlStr)
	req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewReader(payload))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ADMIN-SERVICE-KEY", c.adminKey)
	req.Header.Set("X-ADMIN-SERVICE-SECRET", c.adminSecret)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.Error("admin POST %s failed: %v", urlStr, err)
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode >= 400 {
		logger.Error("admin POST %s status=%d", urlStr, resp.StatusCode)
	}
	return data, resp.StatusCode, nil
}

var ErrMissingCredentials = fmt.Errorf("admin service credentials (X-ADMIN-SERVICE-KEY and X-ADMIN-SERVICE-SECRET) must be set")
