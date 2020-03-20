package utils

import (
	"encoding/hex"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

func ParseHexString(hexStr string) ([]byte, error) {
	if len(hexStr) <= 0 {
		return nil, nil
	}

	decoded, err := hex.DecodeString(strings.Replace(hexStr, " ", "", -1))
	if err != nil {
		return []byte{}, errors.Wrap(err, "Can't decode hex string")
	}

	return decoded, nil
}


func ValidateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}
