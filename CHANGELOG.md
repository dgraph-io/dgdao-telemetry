# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [0.2.0] - 2026-07-08

### Added

- Initial release: an OpenTelemetry-backed `typed.Tracer` for the dgdao typed
  client, extracted from dgdao's `typed/otel.go`. Emits `dgdao.<op>` client
  spans with Dgraph database semantic attributes. Install with
  `typed.SetTracer(telemetry.New())`.
