package utils

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func CpfValidator(cpf string) bool {
	cpf = strings.TrimSpace(cpf)
	cpf = strings.ReplaceAll(cpf, ".", "")
	cpf = strings.ReplaceAll(cpf, "-", "")

	if len(cpf) != 11 {
		return false
	}

	if cpf == "00000000000" || cpf == "11111111111" || cpf == "22222222222" ||
		cpf == "33333333333" || cpf == "44444444444" || cpf == "55555555555" ||
		cpf == "66666666666" || cpf == "77777777777" || cpf == "88888888888" ||
		cpf == "99999999999" {
		return false
	}

	sum := 0
	weight := 10
	digit1 := int(cpf[9] - '0')
	for i := 0; i < 9; i++ {
		sum += int(cpf[i]-'0') * weight
		weight--
	}
	rest := sum % 11
	if rest < 2 {
		rest = 0
	} else {
		rest = 11 - rest
	}
	if rest != digit1 {
		return false
	}

	sum = 0
	weight = 11
	digit2 := int(cpf[10] - '0')
	for i := 0; i < 10; i++ {
		sum += int(cpf[i]-'0') * weight
		weight--
	}
	rest = sum % 11
	if rest < 2 {
		rest = 0
	} else {
		rest = 11 - rest
	}
	if rest != digit2 {
		return false
	}

	return true
}

func CryptPassword(password *string) (*string, error) {
	if password == nil {
		return nil, errors.New("password is nil")
	}
	hashedInput := sha256.Sum256([]byte(*password))
	trimmedHash := bytes.Trim(hashedInput[:], "\x00")
	preparedPassword := string(trimmedHash)
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(preparedPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	bcryptPasswordString := string(bcryptPassword)

	return PString(bcryptPasswordString), nil
}

func CompareHash(password, hash *string) bool {
	if password == nil || hash == nil {
		return false
	}

	hashedInput := sha256.Sum256([]byte(*password))
	trimmedHash := bytes.Trim(hashedInput[:], "\x00")
	preparedPassword := string(trimmedHash)

	plainTextInBytes := []byte(preparedPassword)
	hasTextInBytes := []byte(*hash)

	err := bcrypt.CompareHashAndPassword(hasTextInBytes, plainTextInBytes)
	if err != nil {
		return false
	} else {
		return true
	}
}
