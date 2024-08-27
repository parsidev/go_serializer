package serilizer

import (
	"bytes"
	"fmt"
	"reflect"
)

func Serialize(data interface{}) string {
	var buffer bytes.Buffer

	err := serializeValue(reflect.ValueOf(data), &buffer)
	if err != nil {
		return ""
	}

	return buffer.String()
}

func serializeValue(value reflect.Value, buffer *bytes.Buffer) error {
	switch value.Kind() {
	case reflect.String:
		buffer.WriteString(fmt.Sprintf("s:%d:\"%s\";", len(value.String()), value.String()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buffer.WriteString(fmt.Sprintf("i:%d;", value.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buffer.WriteString(fmt.Sprintf("i:%d;", value.Uint()))
	case reflect.Float32, reflect.Float64:
		buffer.WriteString(fmt.Sprintf("d:%f;", value.Float()))
	case reflect.Bool:
		val := 0
		if value.Bool() {
			val = 1
		}
		buffer.WriteString(fmt.Sprintf("b:%d;", val))
	case reflect.Slice, reflect.Array:
		buffer.WriteString(fmt.Sprintf("a:%d:{", value.Len()))
		for i := 0; i < value.Len(); i++ {
			_ = serializeValue(reflect.ValueOf(i), buffer)
			_ = serializeValue(value.Index(i), buffer)
		}
		buffer.WriteString("}")
	case reflect.Map:
		keys := value.MapKeys()
		buffer.WriteString(fmt.Sprintf("a:%d:{", len(keys)))
		for _, key := range keys {
			_ = serializeValue(key, buffer)
			_ = serializeValue(value.MapIndex(key), buffer)
		}
		buffer.WriteString("}")
	case reflect.Struct:
		t := value.Type()
		buffer.WriteString(fmt.Sprintf("a:%d:{", value.NumField()))
		for i := 0; i < value.NumField(); i++ {
			fieldName := t.Field(i).Name
			_ = serializeValue(reflect.ValueOf(fieldName), buffer)
			_ = serializeValue(value.Field(i), buffer)
		}
		buffer.WriteString("}")
	default:
		return fmt.Errorf("unsupported type: %s", value.Kind())
	}
	return nil
}
