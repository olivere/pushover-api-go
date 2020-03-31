package pushover

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
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
	// Message to send. You can use up to 1024 4-byte UTF-8 characters here.
	// Longer messages are automatically truncated.
	Message string
	// HTML enables or disables HTML formatting in the Message (default: false).
	HTML bool
	// Monospace, when enabled, formats the message with a monospace font (default: false).
	Monospace bool
	// Title of the message (optional). If missing, it will use
	// the name of the app. The maximum length is 250 4-byte UTF-8 characters,
	// which will be automatically truncated.
	Title string
	// Attachment (optional) is the name of a file to be attached.
	Attachment string
	// Devices is the names of the devices (optional). By default,
	// the message is sent to all devices.
	Devices []string
	// URL to be shown as a supplement in the message (optional).
	// Use a maximum length of 512 4-byte UTF-8 characters, otherwise
	// the URL will be automatically truncated.
	URL string
	// URLTitle is a title to be shown for the URL (optional).
	// If missing, only the URL will be shown. URLTitle has a
	// maximum length of 100 4-byte UTF-8 characters and will
	// be automatically truncated.
	URLTitle string
	// Priority of the message (optional). By default, we use a Normal priority.
	//
	// If the Priority is Emergency, you must supply Retry and Expire as
	// described https://pushover.net/api#priority.
	Priority Priority
	// Sound to be played (optional).
	Sound string
	// Timestamp specifies the date/time of the message rather than the
	// time your message is received by the API servers (optional).
	Timestamp time.Time

	// Retry is used only when Priority is set to Emergency, in combination
	// with Expire. Retry specifies how often (in seconds) the Pushover
	// servers will send the same notification to the user (a minimum of 30 seconds
	// is enforced).
	//
	// For details, see https://pushover.net/api#priority.
	Retry time.Duration
	// Expire is used when Priority is set to Emergency, in combination with Retry.
	// Expire specifies how many seconds your notification will continue to be retried
	// (every Retry seconds). A maximum of 3 hours or 10800 seconds is enforced.
	//
	// For details, see https://pushover.net/api#priority.
	Expire time.Duration

	// CallbackURL, if set, specifies a publically accessible URL that is invoked
	// when the user acknowledged the message, e.g. when Priority is Emergency.
	CallbackURL string

	// Tags can be used instead of receipts when the server is unable to process
	// receipts. Tags will be stored on Pushover servers and can be used to e.g.
	// cancel receipts by tag.
	Tags []string
}

type messagesAPI struct {
	c *Client
}

// SendResponse is the outcome of the Send method.
type SendResponse struct {
	Receipt string `json:"receipt,omitempty"`
	Status  int    `json:"status,omitempty"`
	Request string `json:"request,omitempty"`
}

// Send a message.
func (api *messagesAPI) Send(ctx context.Context, m Message) (*SendResponse, error) {
	var (
		body        io.Reader
		contentType string
	)
	if m.Attachment != "" {
		buf := &bytes.Buffer{}
		w := multipart.NewWriter(buf)
		contentType = w.FormDataContentType()
		if err := w.WriteField("token", api.c.appToken); err != nil {
			return nil, fmt.Errorf("unable to write token field: %w", err)
		}
		if err := w.WriteField("user", api.c.userKey); err != nil {
			return nil, fmt.Errorf("unable to write user field: %w", err)
		}
		if err := w.WriteField("message", cut(m.Message, 1024, ellipsis)); err != nil {
			return nil, fmt.Errorf("unable to write message field: %w", err)
		}
		if v := m.HTML; v {
			if err := w.WriteField("html", "1"); err != nil {
				return nil, fmt.Errorf("unable to write html field: %w", err)
			}
		}
		if v := m.Monospace; v {
			if err := w.WriteField("monospace", "1"); err != nil {
				return nil, fmt.Errorf("unable to write monospace field: %w", err)
			}
		}
		if v := m.Title; v != "" {
			if err := w.WriteField("title", cut(v, 250, ellipsis)); err != nil {
				return nil, fmt.Errorf("unable to write title field: %w", err)
			}
		}
		if v := m.Devices; len(v) > 0 {
			if err := w.WriteField("device", strings.Join(v, ",")); err != nil {
				return nil, fmt.Errorf("unable to write device field: %w", err)
			}
		}
		if v := m.URL; v != "" {
			if err := w.WriteField("url", cut(v, 512, "")); err != nil {
				return nil, fmt.Errorf("unable to write url field: %w", err)
			}
		}
		if v := m.URLTitle; v != "" {
			if err := w.WriteField("url_title", cut(v, 100, ellipsis)); err != nil {
				return nil, fmt.Errorf("unable to write url_title field: %w", err)
			}
		}
		if v := m.Priority; v != Normal {
			if err := w.WriteField("priority", fmt.Sprint(v)); err != nil {
				return nil, fmt.Errorf("unable to write priority field: %w", err)
			}
			if v == Emergency {
				retry := int64(m.Retry.Seconds())
				if retry < 30 {
					retry = 30
				}
				if err := w.WriteField("retry", fmt.Sprint(retry)); err != nil {
					return nil, fmt.Errorf("unable to write retry field: %w", err)
				}

				expire := int64(m.Expire.Seconds())
				if expire > 10800 {
					expire = 10800
				}
				if err := w.WriteField("expire", fmt.Sprint(expire)); err != nil {
					return nil, fmt.Errorf("unable to write expire field: %w", err)
				}
			}
		}
		if v := m.CallbackURL; v != "" {
			if err := w.WriteField("callback", v); err != nil {
				return nil, fmt.Errorf("unable to write callback field: %w", err)
			}
		}
		if v := m.Sound; v != "" {
			if err := w.WriteField("sound", v); err != nil {
				return nil, fmt.Errorf("unable to write sound field: %w", err)
			}
		}
		if v := m.Timestamp; !v.IsZero() {
			if err := w.WriteField("timestamp", fmt.Sprint(v.Unix())); err != nil {
				return nil, fmt.Errorf("unable to write monospace field: %w", err)
			}
		}
		if v := m.Tags; len(v) > 0 {
			if err := w.WriteField("tags", strings.Join(v, ",")); err != nil {
				return nil, fmt.Errorf("unable to write monospace field: %w", err)
			}
		}
		// Attachment
		{
			contents, err := ioutil.ReadFile(m.Attachment)
			if err != nil {
				return nil, fmt.Errorf("unable to read attachment: %w", err)
			}
			part, err := w.CreateFormFile("attachment", filepath.Base(m.Attachment))
			if err != nil {
				return nil, fmt.Errorf("unable to create form-data part for attachment: %w", err)
			}
			part.Write(contents)
		}
		if err := w.Close(); err != nil {
			return nil, fmt.Errorf("unable to close and write form-data: %w", err)
		}
		body = buf
	} else {
		values := url.Values{}
		values.Add("token", api.c.appToken)
		values.Add("user", api.c.userKey)
		values.Add("message", cut(m.Message, 1024, ellipsis))
		if v := m.HTML; v {
			values.Add("html", "1")
		}
		if v := m.Monospace; v {
			values.Add("monospace", "1")
		}
		if v := m.Title; v != "" {
			values.Add("title", cut(v, 250, ellipsis))
		}
		if v := m.Devices; len(v) > 0 {
			values.Add("device", strings.Join(v, ","))
		}
		if v := m.URL; v != "" {
			values.Add("url", cut(v, 512, ""))
		}
		if v := m.URLTitle; v != "" {
			values.Add("url_title", cut(v, 100, ellipsis))
		}
		if v := m.Priority; v != Normal {
			values.Add("priority", fmt.Sprint(v))
			if v == Emergency {
				retry := int64(m.Retry.Seconds())
				if retry < 30 {
					retry = 30
				}
				values.Add("retry", fmt.Sprint(retry))

				expire := int64(m.Expire.Seconds())
				if expire > 10800 {
					expire = 10800
				}
				values.Add("expire", fmt.Sprint(expire))
			}
		}
		if v := m.CallbackURL; v != "" {
			values.Add("callback", v)
		}
		if v := m.Sound; v != "" {
			values.Add("sound", v)
		}
		if v := m.Timestamp; !v.IsZero() {
			values.Add("timestamp", fmt.Sprint(v.Unix()))
		}
		if v := m.Tags; len(v) > 0 {
			values.Add("tags", strings.Join(v, ","))
		}
		body = strings.NewReader(values.Encode())
		contentType = "application/x-www-form-urlencoded"
	}
	req, err := http.NewRequestWithContext(ctx, "POST", "/1/messages.json", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
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

// Limits represents the API limits of the current application,
// i.e the current call limit, the number of remaining calls,
// and the time when the limits are reset.
//
// See https://pushover.net/api#limits for details.
type Limits struct {
	Limit     int64     `json:"limit"`
	Remaining int64     `json:"remaining"`
	Reset     int64     `json:"reset"`
	ResetTime time.Time `json:"-"`
}

// Limits returns the API limits of the current application.
func (api *messagesAPI) Limits(ctx context.Context) (*Limits, error) {
	u, err := url.Parse("/1/apps/limits.json")
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Add("token", api.c.appToken)
	u.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), http.NoBody)
	if err != nil {
		return nil, err
	}
	resp, err := api.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(resp.Body)
	var ret Limits
	if err := parseResponse(resp, &ret); err != nil {
		return nil, err
	}
	if ret.Reset > 0 {
		ret.ResetTime = time.Unix(ret.Reset, 0)
	}
	return &ret, nil
}
