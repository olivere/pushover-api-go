package pushover

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

const (
	// Version of the Pushover API for Go.
	Version = "0.1.0"

	defaultBaseURL = "https://api.pushover.net"
)

// Client represents a Pushover client.
type Client struct {
	tr       http.RoundTripper
	baseURL  string
	url      *url.URL
	appToken string
	userKey  string
	logger   Logger
	ua       string

	appLimit     int64 // # of available API calls (as reported by the last API call)
	appRemaining int64 // # of remaining API calls (as reported by the last API call)
	appReset     int64 // date/time (Unix epoch) when the API call limits are reset

	// Messages allows e.g. sending a notification.
	Messages *messagesAPI
}

// ClientOption for configuring Client settings.
type ClientOption func(*Client)

// NewClient initializes and returns a configured Pushover client.
//
// The client can be configured by options and will pick up the
// defaults from these environment variables (options override
// the environment):
//
//   APP_TOKEN      (default: "")
//   USER_KEY       (default: "")
//   PUSHOVER_URL   (default: "https://api.pushover.net")
func NewClient(options ...ClientOption) (*Client, error) {
	c := &Client{
		tr:       http.DefaultTransport,
		baseURL:  defaultBaseURL,
		appToken: envString("", "APP_TOKEN"),
		userKey:  envString("", "USER_KEY"),
		ua:       fmt.Sprintf("pushover-go-api/%s (%s/%s; Go %s)", Version, runtime.GOOS, runtime.GOARCH, runtime.Version()),
	}
	for _, o := range options {
		o(c)
	}
	if c.baseURL == "" {
		c.baseURL = envString(defaultBaseURL, "PUSHOVER_URL")
	}
	var err error
	c.url, err = url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}
	if c.logger == nil {
		c.logger = &nopLogger{}
	}
	c.Messages = &messagesAPI{c: c}
	return c, nil
}

// Close the client.
func (c *Client) Close() error {
	// Reserved for future use
	return nil
}

// WithTransport allows to use HTTP middleware with a Client.
func WithTransport(tr http.RoundTripper) ClientOption {
	return func(c *Client) {
		c.tr = tr
	}
}

// WithURL specifies the base URL to use.
// It is primarily used in development or testing.
func WithURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithAppToken sets the App Token (or API token) to use for this client.
func WithAppToken(appToken string) ClientOption {
	return func(c *Client) {
		c.appToken = appToken
	}
}

// WithUserKey sets the user key (or user/group key) to use for this client.
func WithUserKey(userKey string) ClientOption {
	return func(c *Client) {
		c.userKey = userKey
	}
}

// WithLogger specifies a new logger.
func WithLogger(logger Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// Do executes the HTTP request.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Set URL
	req.URL.Scheme = c.url.Scheme
	req.URL.Host = c.url.Host

	// Set HTTP header data
	if c.ua != "" {
		req.Header.Set("User-Agent", c.ua)
	}
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	start := time.Now()
	resp, err := c.tr.RoundTrip(req)
	duration := time.Since(start)

	c.logger.Log(req, resp, err, start, duration)

	return resp, err
}

// Limits returns the API limits as reported by the last API call.
// If you want to know the current limits without relying on the last
// API call, use Messages.Limits instead.
//
// See https://pushover.net/api#limits for details.
func (c *Client) Limits() Limits {
	var t time.Time
	if c.appReset > 0 {
		t = time.Unix(c.appLimit, 0)
	}
	return Limits{
		Limit:     c.appLimit,
		Remaining: c.appLimit,
		Reset:     c.appReset,
		ResetTime: t,
	}
}
