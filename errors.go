package pushover

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// apiError is used when something fails.
type apiError struct {
	Inner error `json:"-"`

	StatusCode int      `json:"-"`
	Status     int      `json:"status,omitempty"`
	Request    string   `json:"request,omitempty"`
	User       string   `json:"user,omitempty"`
	Errors     []string `json:"errors,omitempty"`
	Receipt    string   `json:"receipt,omitempty"`
}

// Unwrap implements the errors.Wrapper interface, allowing errors.Is and
// errors.As to work with apiErrors.
func (e apiError) Unwrap() error {
	return e.Inner
}

// Error returns a string representation of the API error.
func (e apiError) Error() string {
	var sb strings.Builder
	sb.WriteString("pushover: ")
	sb.WriteString(fmt.Sprint(e.StatusCode))
	sb.WriteString(" ")
	sb.WriteString(http.StatusText(e.StatusCode))
	sb.WriteString("; status=")
	sb.WriteString(fmt.Sprint(e.Status))
	if e.Request != "" {
		sb.WriteString("; request=")
		sb.WriteString(e.Request)
	}
	if e.User != "" {
		sb.WriteString("; status=")
		sb.WriteString(e.User)
	}
	if len(e.Errors) > 0 {
		sb.WriteString("; errors=[")
		sb.WriteString(strings.Join(e.Errors, "; "))
		sb.WriteString("]")
	}
	return sb.String()
}

func errorFromResponse(resp *http.Response) error {
	if resp == nil || (resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return nil
	}
	var ae apiError
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&ae); err != nil {
		return fmt.Errorf("pushover: request failed with status %d", resp.StatusCode)
	}
	ae.StatusCode = resp.StatusCode
	return ae
}

// IsContextErr returns true if the error is from a context that was canceled
// or its deadline exceeded.
func IsContextErr(err error) bool {
	if err == context.Canceled || err == context.DeadlineExceeded {
		return true
	}
	// This happens e.g. on redirect errors, see https://golang.org/src/net/http/client_test.go#L329
	if ue, ok := err.(*url.Error); ok {
		if ue.Temporary() {
			return true
		}
		// Use of an AWS Signing Transport can result in a wrapped url.Error
		return IsContextErr(ue.Err)
	}
	return false
}

// IsNotFound returns true if the given error indicates that a record
// could not be found.
func IsNotFound(err interface{}) bool {
	return IsStatusCode(err, http.StatusNotFound)
}

// IsStatusCode returns true if the given error indicates a specific
// HTTP status code. The err parameter can be of type *http.Response,
// an internal error, or int (indicating the HTTP status code).
func IsStatusCode(err interface{}, code int) bool {
	switch e := err.(type) {
	case *http.Response:
		return e.StatusCode == code
	case *apiError:
		return e.StatusCode == code
	case apiError:
		return e.StatusCode == code
	case int:
		return e == code
	}
	return false
}
