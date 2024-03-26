package utils

import (
	"fmt"
	"net/http"
)

func ParseHeader(header http.Header, key string) (string, error) {
	values, ok := header[key]
	if !ok {
		return "", fmt.Errorf("header %s doesn't exist", key)
	}

	return values[0], nil
}
