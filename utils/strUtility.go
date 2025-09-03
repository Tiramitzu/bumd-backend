package utils

import (
	"errors"
	"strings"
)

func ValidateAndReturnFilterMap(filter string, fields []string) (map[string]string, error) {
	splits := strings.Split(filter, ".")
	if len(splits) != 2 {
		return nil, errors.New("malformed sortBy query parameter, should be field.orderdirection")
	}
	field, value := splits[0], splits[1]
	if !StringInSlice(fields, field) {
		return nil, errors.New("unknown field in filter query parameter")
	}
	return map[string]string{field: value}, nil
}

func StringInSlice(strSlice []string, s string) bool {
	for _, v := range strSlice {
		if v == s {
			return true
		}
	}
	return false
}
