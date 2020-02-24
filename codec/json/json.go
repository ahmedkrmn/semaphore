package json

import (
	"io"

	"github.com/francoispqt/gojay"
	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
)

// New constructs a new JSON encode/decode manager
func New(resource string, schema schema.Object, specs specs.Object) (codec.Manager, error) {
	return &Manager{
		resource: resource,
		specs:    specs,
	}, nil
}

// Manager manages a specs object and allows to encode/decode messages
type Manager struct {
	resource string
	specs    specs.Object
	keys     int
}

// Marshal marshals the given reference store into a JSON message.
// This method is called during runtime to encode a new message with the values stored inside the given reference store
func (manager *Manager) Marshal(refs *refs.Store) (io.Reader, error) {
	reader, writer := io.Pipe()
	encoder := gojay.BorrowEncoder(writer)

	go func() {
		defer writer.Close()
		defer encoder.Release()

		object := NewObject(manager.resource, manager.specs, refs)
		err := encoder.Encode(object)
		if err != nil {
			return
		}
	}()

	return reader, nil
}

// Unmarshal unmarshals the given JSON io reader into the given reference store.
// This method is called during runtime to decode a new message and store it inside the given reference store
func (manager *Manager) Unmarshal(reader io.Reader, refs *refs.Store) error {
	object := NewObject(manager.resource, manager.specs, refs)
	decoder := gojay.BorrowDecoder(reader)
	defer decoder.Release()

	err := decoder.DecodeObject(object)
	if err != nil {
		return err
	}

	return nil
}
