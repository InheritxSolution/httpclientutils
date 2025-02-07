package httpclientutils

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/clbanning/mxj/v2"
)

// RequestOptions holds the configuration for the HTTP request.
type RequestOptions struct {
	Method            string
	URL               string
	Body              interface{}
	Headers           map[string]string
	TLSConfig         *tls.Config
	Timeout           time.Duration
	BasicAuth         *BasicAuthOptions
	ResolveResp       interface{}
	XMLToJSON         interface{}
	DisableEscapeHTML bool
}

// BasicAuthOptions holds the username and password for basic authentication.
type BasicAuthOptions struct {
	Username string
	Password string
}

// Option is a functional option for configuring RequestOptions.
type Option func(*RequestOptions)

// Functional option setters

func WithMethod(method string) Option { return func(opts *RequestOptions) { opts.Method = method } }

func WithURL(url string) Option { return func(opts *RequestOptions) { opts.URL = url } }

func WithBody(body interface{}) Option { return func(opts *RequestOptions) { opts.Body = body } }

func WithHeaders(headers map[string]string) Option {
	return func(opts *RequestOptions) { opts.Headers = headers }
}
func WithTLSConfig(config *tls.Config) Option {
	return func(opts *RequestOptions) { opts.TLSConfig = config }
}
func WithTimeout(timeout time.Duration) Option {
	return func(opts *RequestOptions) { opts.Timeout = timeout }
}
func WithBasicAuth(username, password string) Option {
	return func(opts *RequestOptions) { opts.BasicAuth = &BasicAuthOptions{Username: username, Password: password} }
}
func WithResolveResponse(resp interface{}) Option {
	return func(opts *RequestOptions) { opts.ResolveResp = resp }
}
func WithResolveXMLToJSON(resp interface{}) Option {
	return func(opts *RequestOptions) { opts.XMLToJSON = resp }
}
func WithDisableEscapeHTML(disable bool) Option {
	return func(opts *RequestOptions) { opts.DisableEscapeHTML = disable }
}

// MakeHTTPRequest sends an HTTP request with the provided options.
func MakeHTTPRequest(opts ...Option) (int, http.Header, []byte, error) {
	options := &RequestOptions{Method: http.MethodGet, Headers: make(map[string]string)}
	for _, opt := range opts {
		opt(options)
	}

	body, err := prepareBody(options.Body, options.DisableEscapeHTML)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to prepare request body: %w", err)
	}

	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: options.TLSConfig},
		Timeout:   options.Timeout,
	}

	req, err := http.NewRequest(options.Method, options.URL, body)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}
	if options.BasicAuth != nil {
		req.SetBasicAuth(options.BasicAuth.Username, options.BasicAuth.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return http.StatusRequestTimeout, nil, nil, fmt.Errorf("request timed out: %w", err)
		}
		return 0, nil, nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, resp.Header, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if options.ResolveResp != nil {
		if err := resolveResponse(resp.Header.Get("Content-Type"), responseBody, options.ResolveResp, options.XMLToJSON); err != nil {
			return resp.StatusCode, resp.Header, responseBody, fmt.Errorf("failed to resolve response: %w", err)
		}
	}

	return resp.StatusCode, resp.Header, responseBody, nil
}

func prepareBody(body interface{}, disableEscapeHTML bool) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}
	switch v := body.(type) {
	case string:
		return strings.NewReader(v), nil
	case []byte:
		return bytes.NewReader(v), nil
	default:
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(!disableEscapeHTML)
		if err := enc.Encode(v); err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		return &buf, nil
	}
}

func resolveResponse(contentType string, body []byte, resolveResp, xmlToJson interface{}) error {
	contentType = strings.Split(contentType, ";")[0]

	switch {
	case strings.Contains(contentType, "application/json"):
		if err := json.Unmarshal(body, resolveResp); err != nil {
			return fmt.Errorf("failed to unmarshal JSON response: %w", err)
		}
	case strings.Contains(contentType, "application/xml"):
		m, err := mxj.NewMapXml(body)
		if err != nil {
			return fmt.Errorf("failed to parse XML response: %w", err)
		}
		jsonData, err := m.Json()
		if err != nil {
			return fmt.Errorf("failed to convert XML to JSON: %w", err)
		}
		if xmlToJson != nil {
			if err := json.Unmarshal(jsonData, xmlToJson); err != nil {
				return fmt.Errorf("failed to unmarshal XML to JSON: %w", err)
			}
		}
		if resolveResp != nil {
			return json.Unmarshal(jsonData, resolveResp)
		}
	default:
		return fmt.Errorf("unsupported content type: %s", contentType)
	}

	return nil
}
