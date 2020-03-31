package pushover

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

var byteBufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// parseResponse deserializes the HTTP response body into dst as JSON.
// A maximum size of 8 MB of JSON are permitted.
func parseResponse(resp *http.Response, dst interface{}) error {
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return errorFromResponse(resp)
	}

	buf := byteBufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer byteBufPool.Put(buf)

	// Limit to 8 MB of JSON
	if err := json.NewDecoder(io.TeeReader(io.LimitReader(resp.Body, 8<<20), buf)).Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON data: %v, on input: %s", err, buf.Bytes())
	}
	return nil
}

// closeBody closes rc.
func closeBody(rc io.ReadCloser) {
	if rc != nil {
		rc.Close()
	}
}
