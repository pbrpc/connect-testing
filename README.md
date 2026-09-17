# connect-testing

`connect-testing` provides Connect-specific test doubles for client transports
and streams. Tests can script stream opens, messages, send failures, and
connection failures while inspecting how many streams the client opened.

## Installation

```bash
go get github.com/pbrpc/connect-testing
```

## Transport

The `mocks/transport` package provides:

- `Mock`, a `connect.Transport` whose stream opens are answered by an `Open`
  function and counted by `Opened`.
- `Unopenable`, a transport that answers every stream-open attempt with an
  error.
- `Once`, a transport that answers each open with the same stream or error.
- `Stream`, a `connect.ClientStream` with ordered receive messages and a
  configurable send failure.

Create a transport that always fails to open a stream:

```go
wire := transport.Unopenable(errUnavailable, nil)
rpc := connect.NewClient(wire)
```

Create a scripted stream that receives two messages and begins failing sends
after one successful send:

```go
stream := transport.NewStream(
	errSend,
	1,
	wrapperspb.String("first"),
	wrapperspb.String("second"),
)
wire := transport.Once(stream, nil)
rpc := connect.NewClient(wire)
```
