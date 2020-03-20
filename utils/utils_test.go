package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseHexString(t *testing.T) {
	res, err := ParseHexString("AA A6 12 0B")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, []byte{0xaa, 0xa6, 0x12, 0xb}, res)

	res2, err := ParseHexString("aaA6120B")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, []byte{0xaa, 0xa6, 0x12, 0xb}, res2)

	_, err = ParseHexString("not hex string")
	assert.Error(t, err)

	res3, err := ParseHexString("")
	if err != nil {
		t.Error(err)
	}
	assert.Nil(t, res3)
}

func TestValidateEmail(t *testing.T) {
	assert.True(t, ValidateEmail("email@gmail.com"))
	assert.False(t, ValidateEmail("email"))
}
