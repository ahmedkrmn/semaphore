package xml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

func decodeElement(decoder *xml.Decoder, start xml.StartElement, resource, prefix, name string, template specs.Template, store references.Store) (err error) {
	defer func() {
		if err != nil {
			err = errFailedToDecodeProperty{
				property: name,
				inner:    err,
			}
		}
	}()

	log.Println("<", start.Name.Local, ">", name)

	if start.Name.Local != name {
		return fmt.Errorf("unexpected '%s', expected '%s'", start.Name.Local, name)
	}

	var unmarshaler xml.Unmarshaler

	switch {
	case template.Message != nil:
		unmarshaler = NewObject(resource, buildPath(prefix, name), name, template.Message, store)
	case template.Repeated != nil:
		return errors.New("repeated: not implemented")
	case template.Enum != nil:
		unmarshaler = NewEnum(resource, prefix, name, template.Enum, template.Reference, store)
	case template.Scalar != nil:
		unmarshaler = NewScalar(resource, prefix, name, template.Scalar, template.Reference, store)
	default:
		return fmt.Errorf("property '%s' has unknown type", name)
	}

	return unmarshaler.UnmarshalXML(decoder, start)
}

func buildPath(prefix, property string) string {
	if prefix == "" {
		return property
	}

	return prefix + "." + property
}
