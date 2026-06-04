# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

- Initial release: an OpenTelemetry-backed `typed.Tracer` for the modusgraph typed
  client, extracted from the fork's `typed/otel.go`. Emits `modusgraph.<op>` client
  spans with Dgraph database semantic attributes. Install with
  `typed.SetTracer(telemetry.New())`.
