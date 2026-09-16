//revive:disable:package-comments
package listener

import (
	"net"
	"sync"
)

// Mock is a net.Listener that doesn't bind to a real port.
type Mock struct {
	acceptCh chan net.Conn
	closed   bool
	addr     addr
}

// New creates a new mock listener.
func New() *Mock {
	return &Mock{
		acceptCh: make(chan net.Conn),
		addr:     addr{},
	}
}

// Accept waits for and returns the next connection.
// Returns net.ErrClosed when the listener is closed.
func (l *Mock) Accept() (net.Conn, error) {
	conn, ok := <-l.acceptCh
	if !ok {
		return nil, net.ErrClosed
	}
	return conn, nil
}

// Close closes the listener.
func (l *Mock) Close() error {
	if !l.closed {
		l.closed = true
		close(l.acceptCh)
	}
	return nil
}

// Addr returns the listener's network address.
func (l *Mock) Addr() net.Addr {
	return l.addr
}

type addr struct{}

func (addr) Network() string { return "tcp" }
func (addr) String() string  { return "127.0.0.1:0" }

// Failing is a net.Listener that returns an error on Accept.
type Failing struct {
	err  error
	addr addr
}

// NewFailing creates a listener that returns the given error on Accept.
func NewFailing(err error) *Failing {
	return &Failing{err: err}
}

// Accept always returns the configured error.
func (l *Failing) Accept() (net.Conn, error) {
	return nil, l.err
}

// Close is a no-op.
func (l *Failing) Close() error {
	return nil
}

// Addr returns the listener's network address.
func (l *Failing) Addr() net.Addr {
	return l.addr
}

// Pipe is a net.Listener with no port whose one connection is an in-memory
// pipe. Accept hands out the server end once, then waits for Close; the
// client end is the test's to speak on, so a server can be driven through a
// real connection without a socket.
type Pipe struct {
	conn   net.Conn
	closed chan struct{}
	once   sync.Once
	addr   addr
}

// NewPipe creates a Pipe and answers with it and the client end of its
// connection.
func NewPipe() (*Pipe, net.Conn) {
	server, client := net.Pipe()

	return &Pipe{conn: server, closed: make(chan struct{})}, client
}

// Accept hands out the pipe's server end on the first call and waits for
// Close on every later one, answering net.ErrClosed then.
func (l *Pipe) Accept() (net.Conn, error) {
	if conn := l.conn; conn != nil {
		l.conn = nil

		return conn, nil
	}

	<-l.closed

	return nil, net.ErrClosed
}

// Close ends the listener; a waiting Accept returns.
func (l *Pipe) Close() error {
	l.once.Do(func() { close(l.closed) })

	return nil
}

// Addr returns the listener's network address.
func (l *Pipe) Addr() net.Addr {
	return l.addr
}
