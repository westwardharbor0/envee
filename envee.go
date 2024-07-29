package envee

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Envee struct {
	prefix string
}

const MaxDefinedTagParts = 2

var ErrMissingRequired = errors.New("missing required variable")

func New() *Envee {
	return new(Envee)
}

func (e *Envee) SetPrefix(s string) {
	e.prefix = s
}

func (e *Envee) Parse(o any) error {
	return e.parse(o, "")
}

func (e *Envee) parse(o any, prefix string) error {
	val := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(val.Interface())

	for i := range val.NumField() {
		if typ.Field(i).Type.Kind() == reflect.Struct {
			p := prefix + typ.Field(i).Tag.Get("prefix")

			err := e.parse(val.Field(i).Addr().Interface(), p)
			if err != nil {
				return fmt.Errorf("failed to parse sub struct %q: %w", typ.Field(i).Name, err)
			}

			continue
		}

		f, err := newField(typ.Field(i).Tag)
		if err != nil {
			return err
		}

		envVal, exists := os.LookupEnv(e.prefix + prefix + f.name)

		if !exists && f.isRequired() {
			return fmt.Errorf("%w: %q", ErrMissingRequired, f.name)
		}

		if !exists {
			envVal = f.defValue
		}

		typed, err := parseTypeValue(envVal, val.Field(i).Type())
		if err != nil && !errors.Is(err, ErrUnsupportedType) {
			return fmt.Errorf("%s: %w", f.name, err)
		}

		if err == nil {
			val.Field(i).Set(reflect.ValueOf(typed))
			continue
		}

		kinded, err := parseKindValue(envVal, val.Field(i).Kind())
		if err != nil {
			return fmt.Errorf("%s: %w", f.name, err)
		}

		val.Field(i).Set(reflect.ValueOf(kinded))
	}

	return nil
}

type field struct {
	name     string
	defValue string
	required bool
}

var (
	ErrNotProcessable  = errors.New("not processable, missing required tags")
	ErrUnknownTagValue = errors.New("unknown tag value")
	ErrUnsupportedType = errors.New("unsupported variable type")
)

func (f *field) isRequired() bool {
	return f.required || f.defValue == ""
}

func newField(tags reflect.StructTag) (*field, error) {
	nameVal, exists := tags.Lookup("env")
	if !exists {
		return nil, ErrNotProcessable
	}

	result := new(field)

	if defVal, exists := tags.Lookup("default"); exists {
		result.defValue = defVal
	}

	vSplices := strings.Split(nameVal, ",")
	vLen := len(vSplices)

	if vLen > MaxDefinedTagParts {
		return nil, fmt.Errorf("%w: env:%q", ErrUnknownTagValue, nameVal)
	}

	if vLen == 1 {
		result.name = vSplices[0]
		return result, nil
	}

	if vSplices[0] == "required" {
		result.required = true
		result.name = vSplices[1]
	} else if vSplices[1] == "required" {
		result.required = true
		result.name = vSplices[0]
	} else {
		return nil, fmt.Errorf("%w: env:%q", ErrUnknownTagValue, nameVal)
	}

	return result, nil
}

func parseTypeValue(v string, t reflect.Type) (any, error) {
	switch t.Name() {
	case "Duration":
		return time.ParseDuration(v)
	default:
		return nil, fmt.Errorf("%w: %s:%q", ErrUnsupportedType, t.Name(), v)
	}
}

func parseKindValue(v string, p reflect.Kind) (any, error) {
	switch p {
	case reflect.Int:
		return strconv.Atoi(v)
	case reflect.Uint:
		i, e := strconv.ParseUint(v, 10, 32)
		return uint(i), e
	case reflect.Int8:
		i, e := strconv.ParseInt(v, 10, 8)
		return int8(i), e
	case reflect.Int16:
		i, e := strconv.ParseInt(v, 10, 16)
		return int16(i), e
	case reflect.Int32:
		i, e := strconv.ParseInt(v, 10, 32)
		return int32(i), e
	case reflect.Int64:
		return strconv.ParseInt(v, 10, 64)
	case reflect.Uint8:
		i, e := strconv.ParseUint(v, 10, 8)
		return uint8(i), e
	case reflect.Uint16:
		i, e := strconv.ParseUint(v, 10, 16)
		return uint16(i), e
	case reflect.Uint32:
		i, e := strconv.ParseUint(v, 10, 32)
		return uint32(i), e
	case reflect.Uint64:
		return strconv.ParseUint(v, 10, 64)
	case reflect.String:
		return v, nil
	case reflect.Bool:
		return strconv.ParseBool(v)
	case reflect.Float32:
		f, e := strconv.ParseFloat(v, 32)
		return float32(f), e
	case reflect.Float64:
		return strconv.ParseFloat(v, 64)
	default:
		return nil, fmt.Errorf("%w: %s:%q", ErrUnsupportedType, p, v)
	}
}
