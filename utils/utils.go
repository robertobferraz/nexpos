package utils

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
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

func CreateRandomUsername(name *string) *string {
	var usernameBuffer bytes.Buffer
	usernameBuffer.WriteString(strings.ToLower(strings.ReplaceAll(*name, " ", "_")))
	usernameBuffer.WriteString(fmt.Sprintf(":%v", rand.Int31()))

	return PString(usernameBuffer.String())
}
