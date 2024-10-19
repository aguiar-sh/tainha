package util

import (
	"strings"
)

func PathProtocol(path string) (string, string) {
	protocol := "http"
	if strings.HasPrefix(path, "http://") {

		path = strings.TrimPrefix(path, "http://")
	} else if strings.HasPrefix(path, "https://") {

		protocol = "https:"
		path = strings.TrimPrefix(path, "https://")
	}

	return path, protocol
}
