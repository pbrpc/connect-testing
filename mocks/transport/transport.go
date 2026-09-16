// Package transport fakes the connect.Transport under a *connect.Client, for
// the failures a handler cannot produce: a stream that will not open, a send
// that fails, a receive that never answers. A stream that does open is
// scripted here too; a handler served in-process is the better fake when the
// behavior under test is the handler's.
package transport

import (
	"context"
	"io"
	"sync"

	"connectrpc.com/connect/v2"
	"google.golang.org/protobuf/proto"
)

// Open decides what one NewClientStream answers with. n is the call's
// ordinal, from 1.
type Open func(n int, ctx context.Context, spec connect.Spec) (connect.ClientStream, error)

// Mock is a connect.Transport that answers every NewClientStream through
// Open and counts the calls.
type Mock struct {
	Open Open

	mu    sync.Mutex
	calls int
}

// New creates a Mock answering through open.
func New(open Open) *Mock {
	return &Mock{Open: open}
}

// Unopenable creates a Mock on which no stream opens: every NewClientStream
// answers err. onOpen, when given, runs on each attempt with its ordinal, so
// a test can end a loop that keeps trying.
func Unopenable(err error, onOpen func(n int)) *Mock {
	return New(func(n int, _ context.Context, _ connect.Spec) (connect.ClientStream, error) {
		if onOpen != nil {
			onOpen(n)
		}

		return nil, err
	})
}

// Once creates a Mock answering every open with the same stream, or the same
// error.
func Once(stream connect.ClientStream, err error) *Mock {
	return New(func(int, context.Context, connect.Spec) (connect.ClientStream, error) {
		return stream, err
	})
}

// NewClientStream answers through Open.
func (m *Mock) NewClientStream(ctx context.Context, spec connect.Spec) (connect.ClientStream, error) {
	m.mu.Lock()
	m.calls++
	n := m.calls
	m.mu.Unlock()

	return m.Open(n, ctx, spec)
}

// Opened answers with how many streams were asked for.
func (m *Mock) Opened() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.calls
}

// Stream is a connect.ClientStream with scripted answers. Receive hands out
// Messages in order, merging each into the caller's message, and then blocks
// until the stream is closed, answering io.EOF. Send succeeds FailAfter times
// and then answers SendErr.
type Stream struct {
	Messages  []proto.Message
	FailAfter int
	SendErr   error

	mu       sync.Mutex
	sends    int
	receives int
	closed   chan struct{}
}

// NewStream creates a Stream whose Receive offers messages in order and whose
// Send answers sendErr after failAfter successful sends.
func NewStream(sendErr error, failAfter int, messages ...proto.Message) *Stream {
	return &Stream{
		Messages:  messages,
		FailAfter: failAfter,
		SendErr:   sendErr,
		closed:    make(chan struct{}),
	}
}

// SendHeaders does nothing.
func (s *Stream) SendHeaders() error { return nil }

// Send counts the call and fails once FailAfter sends have succeeded.
func (s *Stream) Send(any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sends++
	if s.sends > s.FailAfter {
		return s.SendErr
	}

	return nil
}

// CloseSend does nothing.
func (s *Stream) CloseSend() error { return nil }

// Receive merges the next scripted message into msg, and once they are spent
// blocks until Close, answering io.EOF.
func (s *Stream) Receive(msg any) error {
	s.mu.Lock()
	n := s.receives
	s.receives++
	s.mu.Unlock()

	if n < len(s.Messages) {
		proto.Merge(msg.(proto.Message), s.Messages[n])

		return nil
	}

	<-s.closed

	return io.EOF
}

// Close ends the stream; a blocked Receive returns.
func (s *Stream) Close() error {
	select {
	case <-s.closed:
	default:
		close(s.closed)
	}

	return nil
}
