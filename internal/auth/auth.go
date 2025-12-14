package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Authorization: ApiKey <api_key>
func GetApiKey(header http.Header) (string, error) {
	result := header.Get("Authorization")
	if result == "" {
		return "", errors.New("missing authorization header")
	}
	results := strings.Split(result, " ")
	if len(results) != 2 || results[0] != "ApiKey" {
		return "", errors.New("invalid api key format")
	}
	return results[1], nil
}
