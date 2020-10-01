package xml

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// var errNotAnObject = errors.New("not an object")

type errUndefinedProperty string

func (e errUndefinedProperty) Error() string {
	return fmt.Sprintf("undefined property %q", string(e))
}

type errUnknownEnum string

func (e errUnknownEnum) Error() string {
	return fmt.Sprintf("unrecognized enum value %q", string(e))
}

type errUnexpectedToken struct {
	actual   xml.Token
	expected []xml.Token
}

func (e errUnexpectedToken) printExpected() string {
	var builder strings.Builder

	for index, token := range e.expected {
		if index > 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(fmt.Sprintf(`"%T"`, token))
	}

	return builder.String()
}

func (e errUnexpectedToken) Error() string {
	return fmt.Sprintf(`unexpected element "%T", expected one of [%s]`, e.actual, e.printExpected())
}

type errFailedToDecodeProperty struct {
	property string
	inner    error
}

func (e errFailedToDecodeProperty) Unwrap() error {
	return e.inner
}

func (e errFailedToDecodeProperty) Error() string {
	return fmt.Sprintf("failed to decode property '%s': %s", e.property, e.inner)
}

// type token struct {
// 	kind xml.Token
// 	name string
// }
//
// func (t token) String() string {
// 	switch {
// 	case t.kind != nil && t.name != "":
// 		return fmt.Sprintf("%T<%s>)", t.kind, t.name)
// 	case t.kind != nil:
// 		return fmt.Sprintf("%T", t.kind)
// 	case t.name != "":
// 		return fmt.Sprintf("<%s>", t.name)
// 	default:
// 		return "<empty>"
// 	}
// }
//
// type errUnexpectedTag struct {
// 	actual, expected token
// }
//
// func (e errUnexpectedTag) Error() string {
// 	return fmt.Sprintf("unexpected tag '%s', expected '%s'", e.actual, e.expected)
// }
