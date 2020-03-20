package utils

import (
	"encoding/hex"
	"github.com/pkg/errors"
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
