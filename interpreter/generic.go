package interpreter

import (
	"fmt"
	"go/ast"
	"reflect"

	"github.com/pkg/errors"
	"github.com/talon-one/decimal"
	"github.com/talon-one/talang/token"
)

var decimalType = reflect.TypeOf(decimal.Decimal{})

type Unmarshaler interface {
	UnmarshalTaToken(*token.TaToken) error
}

type Marshaler interface {
	MarshalTaToken() (*token.TaToken, error)
}

func genericSetConv(value interface{}) (*token.TaToken, error) {
	v := reflect.ValueOf(value)
	// walk down until we can address something
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}

	if v.Type() == decimalType {
		return token.NewDecimal(v.Interface().(decimal.Decimal)), nil
	}

	for {
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (e.Elem().Kind() == reflect.Ptr) {
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if v.Type().NumMethod() > 0 {
			if u, ok := v.Interface().(Marshaler); ok {
				return u.MarshalTaToken()
			}
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		m := make(map[string]*token.TaToken)
		for i := 0; i < v.NumField(); i++ {
			if fieldStruct := v.Type().Field(i); ast.IsExported(fieldStruct.Name) {
				structValue, err := genericSetConv(v.Field(i).Interface())
				if err != nil {
					return nil, err
				}
				m[fieldStruct.Name] = structValue
			}
		}
		return token.NewMap(m), nil
	case reflect.Map:
		m := make(map[string]*token.TaToken)
		if v.Type().Key().Kind() != reflect.String {
			return nil, errors.New("A different key than `string' is not supported")
		}
		for _, key := range v.MapKeys() {
			var err error
			m[key.String()], err = genericSetConv(v.MapIndex(key).Interface())
			if err != nil {
				return nil, err
			}
		}
		return token.NewMap(m), nil
	case reflect.Slice:
		size := v.Len()
		s := make([]*token.TaToken, size, size)
		for i := 0; i < size; i++ {
			var err error
			s[i], err = genericSetConv(v.Index(i).Interface())
			if err != nil {
				return nil, err
			}
		}
		return token.NewList(s...), nil
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		return token.NewDecimalFromString(fmt.Sprintf("%v", value)), nil
	case reflect.String:
		return token.NewString(value.(string)), nil
	case reflect.Bool:
		return token.NewBool(value.(bool)), nil
	case reflect.Float32:
		return token.NewDecimalFromFloat(float64(value.(float32))), nil
	case reflect.Float64:
		return token.NewDecimalFromFloat(value.(float64)), nil
	}
	return nil, errors.Errorf("Unknown type `%T'", value)
}

func (interp *Interpreter) GenericSet(key string, value interface{}) error {
	block, err := genericSetConv(value)
	if err != nil {
		return err
	}

	if len(key) == 0 {
		interp.Binding = block
	} else {
		interp.Set(key, block)
	}
	return nil
}

func genericGetConv(tkn *token.TaToken, v reflect.Value) (reflect.Value, error) {
	// walk down until we can address something
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}

	for {
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (e.Elem().Kind() == reflect.Ptr) {
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if v.Type().NumMethod() > 0 {
			if u, ok := v.Interface().(Unmarshaler); ok {
				if err := u.UnmarshalTaToken(tkn); err != nil {
					return reflect.Value{}, err
				}
				return v.Elem(), nil
			}
		}
		v = v.Elem()
	}

	if v.Type() == decimalType {
		d, err := decimal.NewFromString(tkn.String)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(d), nil
	}

	var result interface{}
	var err error

	switch v.Kind() {
	case reflect.Struct:
		if !tkn.IsMap() {
			return reflect.Value{}, errors.Errorf("%s is not a map", tkn.String)
		}

		for i := 0; i < v.NumField(); i++ {
			if fieldStruct := v.Type().Field(i); ast.IsExported(fieldStruct.Name) {
				item, err := genericGetConv(tkn.MapItem(fieldStruct.Name), v.Field(i))
				if err != nil {
					return reflect.Value{}, err
				}
				if item.Kind() != reflect.Invalid {
					field := v.Field(i)
					for ; field.Kind() == reflect.Ptr; field = field.Elem() {
						field.Set(reflect.New(field.Type().Elem()))
					}
					field.Set(item)
				}
			}
		}
		return reflect.Value{}, nil
	case reflect.Map:
		if !tkn.IsMap() {
			return reflect.Value{}, errors.Errorf("%s is not a map", tkn.String)
		}
		if v.Type().Key().Kind() != reflect.String {
			return reflect.Value{}, errors.New("A different key than `string' is not supported")
		}

		valueType := v.Type().Elem()

		m := reflect.MakeMap(v.Type())

		for key, value := range tkn.Map() {
			buf := reflect.New(valueType).Elem()
			item, err := genericGetConv(value, buf)
			if err != nil {
				return reflect.Value{}, err
			}
			if item.Kind() != reflect.Invalid {
				m.SetMapIndex(reflect.ValueOf(key), item)
			}
		}
		return m, nil
	case reflect.Slice:
		if !tkn.IsList() {
			return reflect.Value{}, errors.Errorf("%s is not a list", tkn.String)
		}
		valueType := v.Type().Elem()
		size := len(tkn.Children)
		s := reflect.MakeSlice(v.Type(), size, size)
		for i := 0; i < size; i++ {
			buf := reflect.New(valueType).Elem()
			item, err := genericGetConv(tkn.Children[i], buf)
			if err != nil {
				return reflect.Value{}, err
			}
			if item.Kind() != reflect.Invalid {
				s.Index(i).Set(item)
			}
		}
		return s, nil
	case reflect.Int:
		if !tkn.IsDecimal() {
			return reflect.Value{}, errors.Errorf("%s is not a decimal", tkn.String)
		}
		result, err = tkn.Decimal.Int()
		if err != nil {
			return reflect.Value{}, err
		}
	case reflect.Int8:
		if !tkn.IsDecimal() {
			return reflect.Value{}, errors.Errorf("%s is not a decimal", tkn.String)
		}
		result, err = tkn.Decimal.Int8()
		if err != nil {
			return reflect.Value{}, err
		}
	case reflect.Int16:
		if !tkn.IsDecimal() {
			return reflect.Value{}, errors.Errorf("%s is not a decimal", tkn.String)
		}
		result, err = tkn.Decimal.Int16()
		if err != nil {
			return reflect.Value{}, err
		}
	case reflect.Int32:
		if !tkn.IsDecimal() {
			return reflect.Value{}, errors.Errorf("%s is not a decimal", tkn.String)
		}
		result, err = tkn.Decimal.Int32()
		if err != nil {
			return reflect.Value{}, err
		}
	case reflect.Int64:
		if !tkn.IsDecimal() {
			return reflect.Value{}, errors.Errorf("%s is not a decimal", tkn.String)
		}
		result, err = tkn.Decimal.Int64()
		if err != nil {
			return reflect.Value{}, err
		}
	case reflect.Uint:
		if !tkn.IsDecimal() {
			return reflect.Value{}, errors.Errorf("%s is not a decimal", tkn.String)
		}
		result, err = tkn.Decimal.Uint()
		if err != nil {
			return reflect.Value{}, err
		}
	case reflect.Uint8:
		if !tkn.IsDecimal() {
			return reflect.Value{}, errors.Errorf("%s is not a decimal", tkn.String)
		}
		result, err = tkn.Decimal.Uint8()
		if err != nil {
			return reflect.Value{}, err
		}
	case reflect.Uint16:
		if !tkn.IsDecimal() {
			return reflect.Value{}, errors.Errorf("%s is not a decimal", tkn.String)
		}
		result, err = tkn.Decimal.Uint16()
		if err != nil {
			return reflect.Value{}, err
		}
	case reflect.Uint32:
		if !tkn.IsDecimal() {
			return reflect.Value{}, errors.Errorf("%s is not a decimal", tkn.String)
		}
		result, err = tkn.Decimal.Uint32()
		if err != nil {
			return reflect.Value{}, err
		}
	case reflect.Uint64:
		if !tkn.IsDecimal() {
			return reflect.Value{}, errors.Errorf("%s is not a decimal", tkn.String)
		}
		result, err = tkn.Decimal.Uint64()
		if err != nil {
			return reflect.Value{}, err
		}
	case reflect.String:
		if !tkn.IsString() {
			return reflect.Value{}, errors.Errorf("%s is not a string", tkn.String)
		}
		result = tkn.String
	case reflect.Bool:
		if !tkn.IsBool() {
			return reflect.Value{}, errors.Errorf("%s is not a bool", tkn.String)
		}
		result = tkn.Bool
	case reflect.Float32:
		if !tkn.IsDecimal() {
			return reflect.Value{}, errors.Errorf("%s is not a decimal", tkn.String)
		}
		result, err = tkn.Decimal.Float32()
		if err != nil {
			return reflect.Value{}, err
		}
	case reflect.Float64:
		if !tkn.IsDecimal() {
			return reflect.Value{}, errors.Errorf("%s is not a decimal", tkn.String)
		}
		result, err = tkn.Decimal.Float64()
		if err != nil {
			return reflect.Value{}, err
		}
	default:
		return reflect.Value{}, errors.Errorf("Unknown type `%T'", v.Interface())
	}
	return reflect.ValueOf(result), nil
}

func (interp *Interpreter) GenericGet(key string, value interface{}) error {
	var binding *token.TaToken
	if len(key) == 0 {
		binding = interp.Binding
	} else {
		binding = interp.Binding.MapItem(key)
		if binding.IsNull() {
			return errors.Errorf("Found no Item with `%s'", key)
		}
	}

	reflectValue := reflect.ValueOf(value)

	v, err := genericGetConv(binding, reflectValue)
	if err != nil {
		return err
	}
	if v.Kind() == reflect.Invalid {
		return nil
	}
	if reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	reflectValue.Set(v)
	return nil
}
