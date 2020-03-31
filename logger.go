package pushover

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strconv"
	"sync"
	"time"
)

// Logger specifies an interface for log output.
type Logger interface {
	// Log request, response, error, start and duration of request.
	Log(*http.Request, *http.Response, error, time.Time, time.Duration) error
}

// -- nopLogger --

// nopLogger doesn't log anything.
type nopLogger struct{}

// Log a roundtrip.
func (nopLogger) Log(_ *http.Request, _ *http.Response, _ error, _ time.Time, _ time.Duration) error {
	return nil
}

// bufPool is a pool for *bytes.Buffer.
var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// -- JSONLogger --

// JSONLogger logs output formatted as JSON.
type JSONLogger struct {
	w io.Writer
}

// NewJSONLogger creates a new JSONLogger that writes to w.
func NewJSONLogger(w io.Writer) *JSONLogger {
	return &JSONLogger{w: w}
}

// Log a roundtrip.
func (l *JSONLogger) Log(req *http.Request, resp *http.Response, err error, start time.Time, duration time.Duration) error {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	quote := func(s string) []byte {
		b := make([]byte, 0, 200)
		return strconv.AppendQuote(b, s)
	}

	buf.WriteRune('{')

	buf.WriteString(`"@timestamp":"`)
	buf.WriteString(start.UTC().Format(time.RFC3339))
	buf.WriteRune('"')

	buf.WriteString(`,"event":{`)

	buf.WriteString(`"duration":`)
	buf.WriteString(fmt.Sprint(duration.Nanoseconds()))

	buf.WriteRune(',')
	buf.WriteString(`"request":{`)
	buf.WriteString(`"url":"`)
	buf.WriteString(req.URL.String())
	buf.WriteRune('"')
	buf.WriteRune('}') // end of request

	buf.WriteRune(',')
	buf.WriteString(`"response":{`)
	if resp != nil {
		buf.WriteString(`"status_code":`)
		buf.WriteString(fmt.Sprint(resp.StatusCode))
	}
	if err != nil {
		if resp != nil {
			buf.WriteRune(',')
		}
		buf.WriteString(`"error":{"message":`)
		buf.Write(quote(err.Error()))
		buf.WriteRune('}')
	}
	buf.WriteRune('}') // end of response

	buf.WriteRune('}') // end of event

	buf.WriteRune('}')
	buf.WriteRune('\n')

	buf.WriteTo(l.w)

	return nil
}

// -- RawLogger --

// RawLogger logs output raw HTTP requests and responses.
type RawLogger struct {
	w io.Writer
}

// NewRawLogger creates a new RawLogger that writes to w.
func NewRawLogger(w io.Writer) *RawLogger {
	return &RawLogger{w: w}
}

// Log a roundtrip.
func (l *RawLogger) Log(req *http.Request, resp *http.Response, err error, start time.Time, duration time.Duration) error {
	out, _ := httputil.DumpRequest(req, true)
	fmt.Fprintln(l.w, string(out))

	out, _ = httputil.DumpResponse(resp, true)
	fmt.Fprintln(l.w, string(out))

	return nil
}
