package client

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/baiyz0825/outline-wiki-sync/utils/xlog"
)

// LoggingTransport is a custom Transport that xlog.Logs requests and responses
type LoggingTransport struct {
	Transport http.RoundTripper
}

func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// xlog.Log request
	var reqBodyBytes []byte
	if req.Body != nil {
		reqBodyBytes, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
	}
	xlog.Log.Debugf("Request: %s %s %s\nHeaders: %v\nBody: %s",
		req.Method, req.URL.String(), req.Proto, req.Header, string(reqBodyBytes))

	// Make the request
	resp, err := t.Transport.RoundTrip(req)

	if err != nil {
		xlog.Log.Errorf("Error making request: %v", err)
		return nil, err
	}

	// xlog.Log response
	var respBodyBytes []byte
	if resp.Body != nil {
		respBodyBytes, _ = io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(respBodyBytes))
	}
	xlog.Log.Debugf("Response: %s\nHeaders: %v\nBody: %s\nDuration: %v",
		resp.Status, resp.Header, string(respBodyBytes), time.Since(start))

	return resp, nil
}
