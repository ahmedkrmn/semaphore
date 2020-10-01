package xml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

func decodeElement(decoder *xml.Decoder, resource, name, path string, template specs.Template, store references.Store) error {
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local != name {
				return fmt.Errorf("unexpected '%s', expected '%s'", t.Name.Local, name)
			}

			switch {
			case template.Message != nil:
				if err := decodeMessage(decoder, name, template.Message, store); err != nil {
					return err
				}

				continue
			case template.Repeated != nil:
				return errors.New("repeated: not implemented")
			case template.Enum != nil:
				return NewEnum(resource, name, path, template.Enum, template.Reference, store).UnmarshalXML(decoder, t)
			case template.Scalar != nil:
				return NewScalar(resource, name, path, template.Scalar, template.Reference, store).UnmarshalXML(decoder, t)
			default:
				return fmt.Errorf("property '%s' has unknown type", name)
			}
		case xml.EndElement:
			// element is closed
			return nil
		default:
			return errUnexpectedToken{
				actual: t,
				expected: []xml.Token{
					xml.StartElement{},
				},
			}
		}
	}
}

func decodeMessage(decoder *xml.Decoder, name string, message specs.Message, store references.Store) error {
	log.Println("decodeMessage:", name)

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			property, ok := message[t.Name.Local]
			if !ok {
				return errUndefinedProperty(t.Name.Local)
			}

			if err := decodeElement(decoder, "TODO", property.Name, property.Path, property.Template, store); err != nil {
				return err
			}

		case xml.EndElement:
			if t.Name.Local == name {
				return nil
			}
		}
	}
}
