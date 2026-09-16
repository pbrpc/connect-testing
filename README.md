# connect-testing

Test doubles for services on the connect stack: the pieces a unit test puts
under a server or a client so nothing listens, dials, or touches a disk, and
what was sent can be inspected afterwards.

## Installation

```bash
go get github.com/pbrpc/connect-testing
```

## Packages

| Package                | Doubles                                                                                                                                                             |
| ---------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `mocks/addr`           | A `net.Addr` for a given string                                                                                                                                     |
| `mocks/certificate`    | `SelfSigned`, a `tls.Certificate` made in memory; `PEM`, the same as the strings configuration carries                                                              |
| `mocks/listener`       | `Mock`, a `net.Listener` that accepts nothing until closed; `Failing`, one whose `Accept` errors; `Pipe`, one whose single connection is a `net.Pipe`               |
| `mocks/meter`          | An OpenTelemetry `metric.Meter` whose counters are inspectable, with error injection                                                                                |
| `mocks/responsewriter` | `Broken`, an `http.ResponseWriter` whose body writes fail, keeping the status and headers written                                                                   |
| `mocks/roundtripper`   | `Func`, a function as `http.RoundTripper`; `Recorder`, one that records every request with its body and answers through `Respond` or `Fail`                         |
| `mocks/slog`           | `CaptureHandler`, an `slog.Handler` that keeps records for inspection                                                                                               |
| `mocks/tracer`         | An in-memory tracer provider whose spans and events are inspectable after `EndSpan`                                                                                 |
| `mocks/transport`      | `Mock`, a `connect.Transport` answering each open through a function (`Unopenable`, `Once`); `Stream`, a `connect.ClientStream` with scripted messages and failures |

## Usage

A handler under test is served in-process, `connect.NewServer()` with the
handler registered and `connectinprocess.New` under a `connect.NewClient`, and
the generated client drives it; that needs nothing from here. These doubles are
for what a handler cannot produce: a stream that will not open
(`transport.Unopenable`), a send that fails partway (`transport.NewStream`), a
replica that does not answer at the HTTP layer (`roundtripper.Fail`), a listener
a server can be told to serve on without a port (`listener.Pipe`), a certificate
for a TLS handshake (`certificate.SelfSigned`), and a client that stopped
reading (`responsewriter.NewBroken`).
