// Package roundtripper fakes the http.RoundTripper under an http.Client or a
// proxy, so a request is answered without a socket and what was sent can be
// inspected afterwards.
package roundtripper

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"sync"
)

// Func adapts a function to http.RoundTripper.
type Func func(*http.Request) (*http.Response, error)

// RoundTrip calls f.
func (f Func) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

// Sent is one request the Recorder carried: the request as it arrived and
// its body, read in full, since a transport reads the body and the caller
// cannot afterwards.
type Sent struct {
	Request *http.Request
	Body    string
}

// Recorder is an http.RoundTripper that records every request, body
// included, and answers each through Answer.
type Recorder struct {
	Answer Func

	mu   sync.Mutex
	sent []Sent
}

// Record creates a Recorder answering through answer.
func Record(answer Func) *Recorder {
	return &Recorder{Answer: answer}
}

// RoundTrip reads request's body, records the request, and answers it
// through Answer with the body restored, so Answer can read it too.
func (r *Recorder) RoundTrip(request *http.Request) (*http.Response, error) {
	var body []byte
	if request.Body != nil {
		body, _ = io.ReadAll(request.Body)
		request.Body = io.NopCloser(bytes.NewReader(body))
	}

	r.mu.Lock()
	r.sent = append(r.sent, Sent{Request: request, Body: string(body)})
	r.mu.Unlock()

	return r.Answer(request)
}

// Sent answers with every request carried so far, in order.
func (r *Recorder) Sent() []Sent {
	r.mu.Lock()
	defer r.mu.Unlock()

	return append([]Sent(nil), r.sent...)
}

// Respond answers every request with status, header, and body.
func Respond(status int, header http.Header, body string) Func {
	return func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: status,
			Status:     http.StatusText(status),
			Header:     header.Clone(),
			Body:       io.NopCloser(strings.NewReader(body)),
			Request:    request,
		}, nil
	}
}

// Fail answers every request with err at the transport.
func Fail(err error) Func {
	return func(*http.Request) (*http.Response, error) {
		return nil, err
	}
}
