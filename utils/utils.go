package utils

import (
	"errors"
	"strings"
)

func ValidationBearerToken(token *string) (*string, error) {
	if len(*token) <= 0 {
		return nil, errors.New("invalid Token")
	}

	x := strings.Split(*token, "Bearer ")
	if len(x) < 2 {
		return nil, errors.New("invalid Token")
	}

	return &x[1], nil
}
