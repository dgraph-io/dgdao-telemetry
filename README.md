# dgdao-telemetry

An OpenTelemetry-backed tracer for the
[dgdao](https://github.com/dgraph-io/dgdao) typed client. The typed
client traces every database operation through a pluggable `typed.Tracer` (a no-op by
default); this package provides the OpenTelemetry implementation.

## Install

```
go get github.com/dgraph-io/dgdao-telemetry
```

## Usage

Install the tracer once at startup, after configuring your OpenTelemetry SDK and
exporter:

```go
import (
    "github.com/dgraph-io/dgdao/typed"
    telemetry "github.com/dgraph-io/dgdao-telemetry"
)

func main() {
    // ... configure your OpenTelemetry SDK / exporter ...
    typed.SetTracer(telemetry.New())
}
```

Each typed database operation then emits a `dgdao.<op>` client span carrying the
Dgraph database semantic attributes `db.system=dgraph`, `db.operation.name`, and
`db.collection.name`. With no SDK installed in the process the spans are no-ops.

## License

Apache-2.0. See [LICENSE](LICENSE) and [NOTICE](NOTICE).
