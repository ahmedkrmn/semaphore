package http

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jexia/maestro/pkg/codec/json"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/metadata"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/types"
	"github.com/jexia/maestro/pkg/transport"
)

type DiscardWriter struct {
}

func (d *DiscardWriter) Write(b []byte) (int, error) {
	return ioutil.Discard.Write(b)
}

func (d *DiscardWriter) Close() error {
	return nil
}

func NewMockCaller() *Caller {
	ctx := instance.NewContext()
	caller := &Caller{
		ctx: ctx,
	}
	return caller
}

func TestNewCaller(t *testing.T) {
	ctx := instance.NewContext()
	constructor := NewCaller()
	listener := constructor(ctx)
	if listener == nil {
		t.Fatal("unexpected result, expected a listener to be constructed")
	}

	if listener.Name() != "http" {
		t.Fatalf("unexpected name '%s', expected listener to be called http", listener.Name())
	}
}

func TestCaller(t *testing.T) {
	message := "hello world"
	mock := NewSimpleMockSpecs()

	codec, err := (&json.Constructor{}).New("input", mock)
	if err != nil {
		t.Fatal(err)
	}

	refs := refs.NewReferenceStore(1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message":"` + message + `"}`))
	}))

	defer server.Close()

	service := NewMockService(server.URL, "GET", "/")
	caller, err := NewMockCaller().Dial(service, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	defer caller.Close()

	r, w := io.Pipe()
	rw := &MockResponseWriter{
		header: metadata.MD{},
		writer: w,
	}

	ctx := context.Background()
	req := transport.Request{
		Method: caller.GetMethod("mock"),
	}

	caller.SendMsg(ctx, rw, &req, refs)

	err = codec.Unmarshal(r, refs)
	if err != nil {
		t.Fatal(err)
	}

	ref := refs.Load("input", "message")
	if ref == nil {
		t.Fatal("input:message reference not set")
	}

	result, is := ref.Value.(string)
	if !is {
		t.Fatal("input:message reference is not a string")
	}

	if result != message {
		t.Fatalf("unexpected input:message %s, expected %s", result, message)
	}
}

func TestCallerUnknownMethod(t *testing.T) {
	service := NewMockService("http://localhost", "GET", "/")
	call, err := NewMockCaller().Dial(service, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	method := call.GetMethod("unknown")
	if method != nil {
		t.Fatal("unexpected method returned")
	}

	if len(call.GetMethods()) != 1 {
		t.Fatalf("unexpected methods %+v, expected a single method to be defined", call.GetMethods())
	}
}

func TestCallerReferences(t *testing.T) {
	expected := ":message"
	path := "message"
	resource := ".params"

	service := NewMockService("http://localhost", "GET", "/"+expected)
	call, err := NewMockCaller().Dial(service, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	method := call.GetMethod("mock")
	references := method.References()
	if len(references) != 1 {
		t.Fatalf("unexpected references %+v", references)
	}

	if len(call.GetMethods()) != 1 {
		t.Fatalf("unexpected methods %+v, expected a single method to be defined", call.GetMethods())
	}

	reference := references[0]
	if reference.Path != expected {
		t.Fatalf("unexpected path %s, expected %s", reference.Path, expected)
	}

	if reference.Reference.Resource != resource {
		t.Fatalf("unexpected reference resource %s, expected %s", reference.Reference.Resource, resource)
	}

	if reference.Reference.Path != path {
		t.Fatalf("unexpected reference path %s, expected %s", reference.Reference.Path, path)
	}
}

func TestCallerReferencesLookup(t *testing.T) {
	value := "1"
	expected := "/" + value

	path := "message"
	resource := ".params"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expected {
			t.Log("unexpected url", r.URL, " expected", expected)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))

	defer server.Close()

	service := NewMockService(server.URL, "GET", "/:message")
	caller, err := NewMockCaller().Dial(service, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	method := caller.GetMethod("mock")
	references := method.References()
	if len(references) != 1 {
		t.Fatalf("unexpected references %+v", references)
	}

	references[0].Type = types.String
	references[0].Label = labels.Optional

	store := refs.NewReferenceStore(1)
	ctx := context.Background()
	req := transport.Request{
		Method: method,
	}

	store.StoreValue(resource, path, value)

	rw := &MockResponseWriter{
		header: metadata.MD{},
		writer: &DiscardWriter{},
	}

	err = caller.SendMsg(ctx, rw, &req, store)
	if err != nil {
		t.Fatal(err)
	}
}
