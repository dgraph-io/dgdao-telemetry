/*
 * SPDX-FileCopyrightText: © 2017-2026 Istari Digital, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

package telemetry_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	noop "go.opentelemetry.io/otel/trace/noop"

	"github.com/dgraph-io/dgdao"
	"github.com/dgraph-io/dgdao/typed"

	telemetry "github.com/dgraph-io/dgdao-telemetry"
)

// widget is a minimal schema struct used to exercise the typed client.
type widget struct {
	UID   string   `json:"uid,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
	Name  string   `json:"name,omitempty" dgraph:"index=exact"`
	Qty   int      `json:"qty,omitempty" dgraph:"index=int"`
}

func newConn(t *testing.T) dgdao.Client {
	t.Helper()
	conn, err := dgdao.NewClient("file://"+t.TempDir(), dgdao.WithAutoSchema(true))
	if err != nil {
		t.Fatalf("dgdao.NewClient: %v", err)
	}
	t.Cleanup(conn.Close)
	return conn
}

// recordSpans installs the telemetry tracer into the typed package and an
// in-memory OTel span recorder, so typed-layer DB operations produce inspectable
// spans. Both are reset on cleanup.
func recordSpans(t *testing.T) *tracetest.SpanRecorder {
	t.Helper()
	typed.SetTracer(telemetry.New())
	t.Cleanup(func() { typed.SetTracer(nil) })

	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	otel.SetTracerProvider(tp)
	t.Cleanup(func() { otel.SetTracerProvider(noop.NewTracerProvider()) })
	return sr
}

func spanNames(sr *tracetest.SpanRecorder) []string {
	var names []string
	for _, s := range sr.Ended() {
		names = append(names, s.Name())
	}
	return names
}

func TestClient_CRUD_EmitsSpans(t *testing.T) {
	sr := recordSpans(t)
	ctx := context.Background()
	c := typed.NewClient[widget](newConn(t))

	w := &widget{Name: "sprocket", Qty: 1}
	if err := c.Insert(ctx, w); err != nil {
		t.Fatalf("Insert: %v", err)
	}
	if _, err := c.Get(ctx, w.UID); err != nil {
		t.Fatalf("Get: %v", err)
	}

	var wantInsert, wantGet bool
	for _, n := range spanNames(sr) {
		if n == "dgdao.insert" {
			wantInsert = true
		}
		if n == "dgdao.get" {
			wantGet = true
		}
	}
	if !wantInsert || !wantGet {
		t.Fatalf("missing CRUD spans; got %v", spanNames(sr))
	}

	for _, s := range sr.Ended() {
		if s.Name() != "dgdao.get" {
			continue
		}
		attrs := map[string]string{}
		for _, kv := range s.Attributes() {
			attrs[string(kv.Key)] = kv.Value.AsString()
		}
		if attrs["db.system"] != "dgraph" {
			t.Errorf("db.system = %q, want dgraph", attrs["db.system"])
		}
		if attrs["db.operation.name"] != "get" {
			t.Errorf("db.operation.name = %q, want get", attrs["db.operation.name"])
		}
		if attrs["db.collection.name"] != "widget" {
			t.Errorf("db.collection.name = %q, want widget", attrs["db.collection.name"])
		}
	}
}

func TestQuery_Terminals_EmitSpans(t *testing.T) {
	sr := recordSpans(t)
	ctx := context.Background()
	c := typed.NewClient[widget](newConn(t))
	if err := c.Insert(ctx, &widget{Name: "a", Qty: 1}); err != nil {
		t.Fatalf("Insert: %v", err)
	}

	if _, err := c.Query(ctx).Nodes(); err != nil {
		t.Fatalf("Nodes: %v", err)
	}
	if _, err := c.Query(ctx).First(); err != nil {
		t.Fatalf("First: %v", err)
	}

	var querySpans int
	for _, s := range sr.Ended() {
		if s.Name() == "dgdao.query" {
			querySpans++
		}
	}
	if querySpans < 2 {
		t.Fatalf("want >=2 dgdao.query spans, got %d (%v)", querySpans, spanNames(sr))
	}
}

func TestIterNodes_EmitsSpan(t *testing.T) {
	sr := recordSpans(t)
	ctx := context.Background()
	c := typed.NewClient[widget](newConn(t))
	for i := range 3 {
		if err := c.Insert(ctx, &widget{Name: "w", Qty: i}); err != nil {
			t.Fatalf("Add %d: %v", i, err)
		}
	}

	// IterNodes opens its span inside the returned closure, so the span only
	// ends after iteration completes. Fully drain the iterator before checking.
	seen := 0
	for w, err := range c.Query(ctx).IterNodes() {
		if err != nil {
			t.Fatalf("IterNodes yielded error: %v", err)
		}
		if w == nil {
			t.Fatal("IterNodes yielded a nil widget")
		}
		seen++
	}
	if seen != 3 {
		t.Fatalf("IterNodes yielded %d records, want 3", seen)
	}

	var querySpans int
	for _, s := range sr.Ended() {
		if s.Name() == "dgdao.query" {
			querySpans++
		}
	}
	if querySpans < 1 {
		t.Fatalf("want >=1 dgdao.query span after IterNodes, got %d (%v)", querySpans, spanNames(sr))
	}
}

func TestSpan_RecordsErrorStatus(t *testing.T) {
	sr := recordSpans(t)
	ctx := context.Background()
	c := typed.NewClient[widget](newConn(t))

	// Get against a well-formed but absent UID returns a "node not found" error
	// from the file:// store, exercising the span error path.
	if _, err := c.Get(ctx, "0xdeadbeef"); err == nil {
		t.Fatal("Get(0xdeadbeef) returned no error; expected not-found")
	}

	var found bool
	for _, s := range sr.Ended() {
		if s.Name() != "dgdao.get" {
			continue
		}
		found = true
		if got := s.Status().Code; got != codes.Error {
			t.Errorf("dgdao.get span status code = %v, want Error", got)
		}
	}
	if !found {
		t.Fatalf("no dgdao.get span recorded; got %v", spanNames(sr))
	}
}
