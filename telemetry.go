/*
 * SPDX-FileCopyrightText: © 2017-2026 Istari Digital, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

// Package telemetry provides an OpenTelemetry-backed tracer for the modusgraph
// typed client. The typed client traces every database operation through a
// pluggable typed.Tracer; install this implementation once at startup:
//
//	typed.SetTracer(telemetry.New())
//
// With no OpenTelemetry SDK installed in the process the spans are no-ops;
// configuring an SDK and exporter is the application's job.
package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/matthewmcneely/modusgraph/typed"
)

// tracerName is the instrumentation scope for typed-layer DB spans.
const tracerName = "github.com/mlwelles/modusgraph-telemetry"

// New returns a typed.Tracer backed by the global OpenTelemetry tracer provider.
func New() typed.Tracer { return otelTracer{} }

type otelTracer struct{}

// StartSpan opens a span named "modusgraph.<op>" carrying the Dgraph database
// semantic attributes db.system, db.operation.name, and db.collection.name.
func (otelTracer) StartSpan(ctx context.Context, op, collection string) (context.Context, typed.Span) {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "modusgraph."+op,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system", "dgraph"),
			attribute.String("db.operation.name", op),
			attribute.String("db.collection.name", collection),
		),
	)
	return ctx, otelSpan{span}
}

type otelSpan struct{ span trace.Span }

// End records err (if any) on the span and ends it.
func (s otelSpan) End(err error) {
	if err != nil {
		s.span.RecordError(err)
		s.span.SetStatus(codes.Error, err.Error())
	}
	s.span.End()
}
