package util

import (
	"strings"
)

func GetHeaders(options map[string]string) map[string]string {

	headers := make(map[string]string)
	for key, value := range options {
		if strings.HasPrefix(key, "H-") {
			headers[key[2:]] = value
		}
	}
	return headers
}

func GetCookies(options map[string]string) map[string]string {

	headers := make(map[string]string)
	for key, value := range options {
		if strings.HasPrefix(key, "C-") {
			headers[key[2:]] = value
		}
	}
	return headers
}
