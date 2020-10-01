package xml

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

func decodeElement(decoder *xml.Decoder, start xml.StartElement, resource, name, path string, template specs.Template, store references.Store) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to decode property '%s': %w", name, err)
		}
	}()

	if start.Name.Local != name {
		return fmt.Errorf("unexpected '%s', expected '%s'", start.Name.Local, name)
	}

	switch {
	case template.Message != nil:
		return NewObject(resource, name, path, template.Message, store).UnmarshalXML(decoder, start)
	case template.Repeated != nil:
		return errors.New("repeated: not implemented")
	case template.Enum != nil:
		return NewEnum(resource, name, path, template.Enum, template.Reference, store).UnmarshalXML(decoder, start)
	case template.Scalar != nil:
		return NewScalar(resource, name, path, template.Scalar, template.Reference, store).UnmarshalXML(decoder, start)
	default:
		return fmt.Errorf("property '%s' has unknown type", name)
	}
}
