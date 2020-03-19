package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringToMimeFormat(t *testing.T) {
	f, err := StringToMimeFormat("ascii")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, MimeFormatASCII, f)

	f2, err := StringToMimeFormat("hex")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, MimeFormatHex, f2)

	_, err = StringToMimeFormat("any other string")
	assert.Error(t, err)
}

func TestNdefLangValuesSlice_Contains(t *testing.T) {
	val := NdefLangValues.Contains("English")
	assert.Equal(t, true, val)
	val2 := NdefLangValues.Contains("Swahili‎")
	assert.Equal(t, false, val2)
}
