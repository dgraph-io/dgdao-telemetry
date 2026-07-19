# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [0.2.4] - 2026-07-19

### Changed

- Bump the `github.com/dgraph-io/dgdao` dependency to v0.9.0. The typed client's `Add` is
  now `Insert`, so the emitted span name follows (`dgdao.insert`); no change to the tracer
  implementation itself.

## [0.2.3] - 2026-07-17

### Changed

- chore(deps): bump the `github.com/dgraph-io/dgdao` dependency to v0.8.0.

## [0.2.2] - 2026-07-16

### Changed

- Bump the `github.com/dgraph-io/dgdao` dependency to v0.6.1, which adds
  the `Defaulter` hook (before-validation defaulting on writes) and
  caches the per-write schema check.

## [0.2.1] - 2026-07-09

### Changed

- Bump the `github.com/dgraph-io/dgdao` dependency to v0.5.4, which pins
  Dgraph to the released v25.3.8 tag rather than a pre-release
  pseudo-version.

## [0.2.0] - 2026-07-08

### Added

- Initial release: an OpenTelemetry-backed `typed.Tracer` for the dgdao typed
  client, extracted from dgdao's `typed/otel.go`. Emits `dgdao.<op>` client
  spans with Dgraph database semantic attributes. Install with
  `typed.SetTracer(telemetry.New())`.
