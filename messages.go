package pushover

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Priority settings of a Message.
type Priority int

const (
	// Lowest will generate no alert/notification on the user's devices.
	Lowest Priority = -2
	// Low will send the message as a quiet notification.
	Low Priority = -1
	// Normal is the normal priority.
	Normal Priority = 0
	// High will bypass the user's quiet hours.
	High Priority = 1
	// Emergency will bypass the user's quiet hours and will require
	// a confirmation by the user.
	Emergency Priority = 2
)

// Message to send.
type Message struct {
	// Message to send.
	Message string
	// HTML enables or disables HTML formatting in the Message (default: false).
	HTML bool
	// Monospace, when enabled, formats the message with a monospace font (default: false).
	Monospace bool
	// Title of the message (optional). If missing, it will use
	// the name of the app.
	Title string
	//// Attachment (optional).
	//Attachment io.Reader
	// Devices is the names of the devices (optional). By default,
	// the message is sent to all devices.
	Devices []string
	// URL to be shown as a supplement in the message (optional).
	URL string
	// URLTitle is a title to be shown for the URL (optional).
	// If missing, only the URL will be shown.
	URLTitle string
	// Priority of the message (optional).
	Priority Priority
	// Sound to be played (optional).
	Sound string
	// Timestamp specifies the date/time of the message rather than the
	// time your message is received by the API servers (optional).
	Timestamp time.Time
}

type messagesAPI struct {
	c *Client
}

// SendResponse is the outcome of the Send method.
type SendResponse struct {
	Receipt string `json:"receipt,omitempty"`
}

// Send a message.
func (api *messagesAPI) Send(ctx context.Context, m Message) (*SendResponse, error) {
	values := url.Values{}
	values.Add("token", api.c.appToken)
	values.Add("user", api.c.userKey)
	values.Add("message", m.Message)
	if v := m.HTML; v {
		values.Add("html", "1")
	}
	if v := m.Monospace; v {
		values.Add("monospace", "1")
	}
	if v := m.Title; v != "" {
		values.Add("title", v)
	}
	if v := m.Devices; len(v) > 0 {
		values.Add("device", strings.Join(v, ","))
	}
	if v := m.URL; v != "" {
		values.Add("url", v)
	}
	if v := m.URLTitle; v != "" {
		values.Add("url_title", v)
	}
	if v := m.Priority; v != Normal {
		values.Add("priority", fmt.Sprint(v))
	}
	if v := m.Sound; v != "" {
		values.Add("sound", v)
	}
	if v := m.Timestamp; !v.IsZero() {
		values.Add("timestamp", fmt.Sprint(v.Unix()))
	}

	body := strings.NewReader(values.Encode())
	req, err := http.NewRequestWithContext(ctx, "POST", "/1/messages.json", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := api.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(resp.Body)
	var ret SendResponse
	if err := parseResponse(resp, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
