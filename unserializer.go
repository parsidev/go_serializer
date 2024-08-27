package serilizer

import (
	"errors"
	"strconv"
	"strings"
)

func findNextValue(data string) (interface{}, string, error) {
	switch {
	case strings.HasPrefix(data, "s:"):
		endIndex := strings.Index(data, ";")
		length, _ := strconv.Atoi(data[2:endIndex])
		value := data[endIndex+2 : endIndex+2+length]
		rest := data[endIndex+2+length+1:]
		return value, rest, nil
	case strings.HasPrefix(data, "i:"):
		endIndex := strings.Index(data, ";")
		value, _ := strconv.Atoi(data[2:endIndex])
		rest := data[endIndex+1:]
		return value, rest, nil
	case strings.HasPrefix(data, "d:"):
		endIndex := strings.Index(data, ";")
		value, _ := strconv.ParseFloat(data[2:endIndex], 64)
		rest := data[endIndex+1:]
		return value, rest, nil
	case strings.HasPrefix(data, "b:"):
		endIndex := strings.Index(data, ";")
		value, _ := strconv.Atoi(data[2:endIndex])
		rest := data[endIndex+1:]
		return value == 1, rest, nil
	case strings.HasPrefix(data, "a:"):
		startIndex := strings.Index(data, "{") + 1
		endIndex := strings.LastIndex(data, "}")
		value, _ := unSerializeArray(data[startIndex:endIndex])
		rest := data[endIndex+1:]
		return value, rest, nil
	default:
		return nil, "", errors.New("unknown type")
	}
}

func UnSerialize(data string) (value interface{}) {
	value, _ = unSerializeValue(data)

	return value
}

func unSerializeValue(data string) (interface{}, error) {
	switch {
	case strings.HasPrefix(data, "s:"):
		endIndex := strings.Index(data, "\";")
		startIndex := strings.Index(data, ":\"")
		startIndex = startIndex + len(":\"")
		return data[startIndex:endIndex], nil
	case strings.HasPrefix(data, "i:"):
		endIndex := strings.Index(data, ";")
		value, _ := strconv.Atoi(data[2:endIndex])
		return value, nil
	case strings.HasPrefix(data, "d:"):
		endIndex := strings.Index(data, ";")
		value, _ := strconv.ParseFloat(data[2:endIndex], 64)
		return value, nil
	case strings.HasPrefix(data, "b:"):
		endIndex := strings.Index(data, ";")
		value, _ := strconv.Atoi(data[2:endIndex])
		return value == 1, nil
	case strings.HasPrefix(data, "a:"):
		startIndex := strings.Index(data, "{") + 1
		endIndex := strings.LastIndex(data, "}")
		arrayData := data[startIndex:endIndex]
		return unSerializeArray(arrayData)
	default:
		return nil, errors.New("unknown type")
	}
}

func unSerializeArray(data string) (map[interface{}]interface{}, error) {
	result := make(map[interface{}]interface{})
	for len(data) > 0 {
		key, rest, err := findNextValue(data)
		if err != nil {
			return nil, err
		}
		value, rest, err := findNextValue(rest)
		if err != nil {
			return nil, err
		}
		result[key] = value
		data = rest
	}
	return result, nil
}
