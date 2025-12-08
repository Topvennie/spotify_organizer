package spotify

import (
	"strings"
)

func uriToID(uri string) string {
	parts := strings.Split(uri, ":")
	if len(parts) != 3 {
		return ""
	}

	return parts[2]
}
