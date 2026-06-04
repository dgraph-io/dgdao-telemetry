# modusgraph-telemetry

An OpenTelemetry-backed tracer for the
[modusgraph](https://github.com/matthewmcneely/modusgraph) typed client. The typed
client traces every database operation through a pluggable `typed.Tracer` (a no-op by
default); this package provides the OpenTelemetry implementation.

## Install

`modusgraph-telemetry` depends on a fork of modusGraph published under a different
import path. Go does not propagate `replace` directives to consumers, so your project
must declare the same one:

```go
// go.mod
require (
    github.com/mlwelles/modusgraph-telemetry v0.1.0
    github.com/matthewmcneely/modusgraph v0.0.0-00010101000000-000000000000
)

replace github.com/matthewmcneely/modusgraph => github.com/mlwelles/modusGraph v0.5.0-dev-mlwelles-20260604c
```

## Usage

Install the tracer once at startup, after configuring your OpenTelemetry SDK and
exporter:

```go
import (
    "github.com/matthewmcneely/modusgraph/typed"
    telemetry "github.com/mlwelles/modusgraph-telemetry"
)

func main() {
    // ... configure your OpenTelemetry SDK / exporter ...
    typed.SetTracer(telemetry.New())
}
```

Each typed database operation then emits a `modusgraph.<op>` client span carrying the
Dgraph database semantic attributes `db.system=dgraph`, `db.operation.name`, and
`db.collection.name`. With no SDK installed in the process the spans are no-ops.

## License

Apache-2.0. See [LICENSE](LICENSE) and [NOTICE](NOTICE).
