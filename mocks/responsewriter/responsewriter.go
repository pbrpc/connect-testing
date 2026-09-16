// Package responsewriter fakes an http.ResponseWriter whose client is gone.
package responsewriter

import (
	"errors"
	"net/http"
)

// ErrBroken is what every write to a Broken writer answers with.
var ErrBroken = errors.New("broken pipe")

// Broken is a ResponseWriter whose connection is gone: headers and the status
// go through and are kept for inspection, and the body fails.
type Broken struct {
	header http.Header

	// Status is the code WriteHeader was called with, 0 when it was not.
	Status int
}

// NewBroken creates a Broken writer.
func NewBroken() *Broken {
	return &Broken{header: http.Header{}}
}

// Header answers with the headers set so far.
func (w *Broken) Header() http.Header { return w.header }

// WriteHeader keeps status.
func (w *Broken) WriteHeader(status int) { w.Status = status }

// Write fails with ErrBroken.
func (*Broken) Write([]byte) (int, error) { return 0, ErrBroken }
