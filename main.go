package maestro

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/jexia/maestro/flow"
	"github.com/jexia/maestro/refs"

	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/protocol"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/intermediate"
	"github.com/jexia/maestro/specs/strict"
	"github.com/jexia/maestro/specs/trace"
	"github.com/jexia/maestro/utils"
)

// Client represents a maestro instance
type Client struct {
	Manifest  *specs.Manifest
	Listeners []protocol.Listener
	Options   Options
}

// Serve opens the client listeners
func (client *Client) Serve() {

}

// Option represents a constructor func which sets a given option
type Option func(*Options)

// Options represents all the available options
type Options struct {
	Path      string
	Recursive bool
	Codec     map[string]codec.Constructor
	Callers   []protocol.Caller
	Listeners []protocol.Listener
	Schema    schema.Collection
	Functions specs.CustomDefinedFunctions
}

// NewOptions constructs a options object from the given option constructors
func NewOptions(options ...Option) Options {
	result := Options{
		Codec: make(map[string]codec.Constructor),
	}

	for _, option := range options {
		option(&result)
	}

	return result
}

// WithPath defines the definitions path to be used
func WithPath(path string, recursive bool) Option {
	return func(options *Options) {
		options.Path = path
		options.Recursive = recursive
	}
}

// WithCodec appends the given codec to the collection of available codecs
func WithCodec(constructor codec.Constructor) Option {
	return func(options *Options) {
		options.Codec[constructor.Name()] = constructor
	}
}

// WithCaller appends the given caller to the collection of available callers
func WithCaller(caller protocol.Caller) Option {
	return func(options *Options) {
		options.Callers = append(options.Callers, caller)
	}
}

// WithListener appends the given listener to the collection of available listeners
func WithListener(listener protocol.Listener) Option {
	return func(options *Options) {
		options.Listeners = append(options.Listeners, listener)
	}
}

// WithSchemaCollection defines the schema collection to be used
func WithSchemaCollection(collection schema.Collection) Option {
	return func(options *Options) {
		options.Schema = collection
	}
}

// WithFunctions defines the custom defined functions to be used
func WithFunctions(functions specs.CustomDefinedFunctions) Option {
	return func(options *Options) {
		options.Functions = functions
	}
}

// New constructs a new Maestro instance
func New(opts ...Option) (*Client, error) {
	options := NewOptions(opts...)

	if options.Path == "" {
		return nil, trace.New(trace.WithMessage("undefined path in options"))
	}

	if options.Schema == nil {
		return nil, trace.New(trace.WithMessage("undefined schema in options"))
	}

	manifest, err := ConstructSpecs(options)
	if err != nil {
		return nil, err
	}

	err = ConstructFlowManager(manifest, options)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Manifest:  manifest,
		Listeners: options.Listeners,
		Options:   options,
	}

	return client, nil
}

// ConstructSpecs construct a specs manifest from the given options
func ConstructSpecs(options Options) (*specs.Manifest, error) {
	files, err := utils.ReadDir(options.Path, options.Recursive, intermediate.Ext)
	if err != nil {
		return nil, err
	}

	manifest := &specs.Manifest{}

	for _, file := range files {
		reader, err := os.Open(filepath.Join(file.Path, file.Name()))
		if err != nil {
			return nil, err
		}

		definition, err := intermediate.UnmarshalHCL(file.Name(), reader)
		if err != nil {
			return nil, err
		}

		result, err := intermediate.ParseManifest(definition, options.Functions)
		if err != nil {
			return nil, err
		}

		manifest.MergeLeft(result)

		err = specs.CheckManifestDuplicates(file.Name(), manifest)
		if err != nil {
			return nil, err
		}
	}

	err = specs.ResolveManifestDependencies(manifest)
	if err != nil {
		return nil, err
	}

	err = strict.Define(options.Schema, manifest)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

// ConstructFlowManager constructs the flow managers from the given specs manifest
func ConstructFlowManager(manifest *specs.Manifest, options Options) error {
	endpoints := make([]*protocol.Endpoint, len(manifest.Endpoints))

	for index, endpoint := range manifest.Endpoints {
		f := GetFlow(manifest, endpoint.Flow)
		nodes := make([]*flow.Node, len(f.Nodes))

		for index, node := range f.Nodes {
			caller, err := ConstructCall(manifest, node, node.Call, options)
			if err != nil {
				return err
			}

			// rollback, err := ConstructCall(manifest, call.Rollback, options)
			// if err != nil {
			// 	return err
			// }

			nodes[index] = flow.NewNode(node, caller, nil)
		}

		collection, has := options.Codec[endpoint.Codec]
		if !has {
			return trace.New(trace.WithMessage("unkown endpoint codec %s", endpoint.Codec))
		}

		req, err := collection.New(specs.InputResource, f.GetInput())
		if err != nil {
			return err
		}

		res, err := collection.New(specs.InputResource, f.GetOutput())
		if err != nil {
			return err
		}

		manager := flow.NewManager(f.GetName(), nodes)

		endpoints[index] = &protocol.Endpoint{
			Flow:     manager,
			Listener: endpoint.Listener,
			Options:  endpoint.Options,
			Request:  req,
			Response: res,
		}
	}

	err := ConstructListeners(endpoints, options)
	if err != nil {
		return err
	}

	return nil
}

type rw struct {
	writer io.Writer
}

func (rw *rw) Header() protocol.Header {
	return nil
}
func (rw *rw) Write(bb []byte) (int, error) {
	return rw.writer.Write(bb)
}
func (rw *rw) WriteHeader(int) {}

func ConstructCall(manifest *specs.Manifest, node *specs.Node, call *specs.Call, options Options) (flow.Call, error) {
	service := GetService(manifest, call.GetService())
	if service == nil {
		return nil, trace.New(trace.WithMessage("the service for %s was not found", call.GetMethod()))
	}

	constructor := GetCaller(options.Callers, service.Caller)
	codec := options.Codec[service.Codec]

	req, err := codec.New(node.GetName(), call.GetRequest())
	if err != nil {
		return nil, err
	}

	res, err := codec.New(node.GetName(), call.GetResponse())
	if err != nil {
		return nil, err
	}

	caller, err := constructor.New(service.Host, options.Schema.GetService(service.Schema), service.Options)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context, refs *refs.Store) error {
		body, err := req.Marshal(refs)
		if err != nil {
			return err
		}

		reader, writer := io.Pipe()
		req := &protocol.Request{
			Endpoint: call.GetMethod(),
			Context:  ctx,
			Body:     body,
			// Header:  protocol.Header{},
		}

		w := &rw{
			writer: writer,
		}

		go func() {
			defer writer.Close()
			caller.Call(w, req, refs)
		}()

		err = res.Unmarshal(reader, refs)
		if err != nil {
			return nil
		}

		return nil
	}, nil
}

// ConstructListeners constructs the listeners from the given collection of endpoints
func ConstructListeners(endpoints []*protocol.Endpoint, options Options) error {
	collections := make(map[string][]*protocol.Endpoint, len(options.Listeners))

	for _, endpoint := range endpoints {
		listener := GetListener(options.Listeners, endpoint.Listener)
		if listener == nil {
			return trace.New(trace.WithMessage("unkown listener %s", endpoint.Listener))
		}

		collections[endpoint.Listener] = append(collections[endpoint.Listener], endpoint)
	}

	for key, collection := range collections {
		listener := GetListener(options.Listeners, key)
		err := listener.Handle(collection)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetListener attempts to retrieve the requested listener
func GetListener(listeners []protocol.Listener, name string) protocol.Listener {
	for _, listener := range listeners {
		if listener.Name() == name {
			return listener
		}
	}

	return nil
}

// GetCaller attempts to retrieve a caller from the given options matching the given name
func GetCaller(callers []protocol.Caller, name string) protocol.Caller {
	for _, caller := range callers {
		if caller.Name() == name {
			return caller
		}
	}

	return nil
}

// GetService attempts to retrieve a service from the given manifest matching the given name
func GetService(manifest *specs.Manifest, name string) *specs.Service {
	for _, service := range manifest.Services {
		if service.Name == name {
			return service
		}
	}

	return nil
}

// GetFlow attempts to retrieve a flow from the given manifest matching the given name
func GetFlow(manifest *specs.Manifest, name string) *specs.Flow {
	for _, flow := range manifest.Flows {
		if flow.GetName() == name {
			return flow
		}
	}

	return nil
}
