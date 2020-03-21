package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/taglme/nfc-cli/models"
	"testing"
)

func Test_validateNdefRecordPayloadRaw(t *testing.T) {
	p, err := validateNdefRecordPayloadRaw(25, "", "", "")
	assert.EqualError(t, err, "Wrong tnf flag value. Can be only from 0 to 6")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadRaw(-1, "", "", "")
	assert.EqualError(t, err, "Wrong tnf flag value. Can be only from 0 to 6")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadRaw(2, "", "", "")
	assert.EqualError(t, err, "Payload value can't be empty")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadRaw(2, "", "", "some")
	assert.Error(t, err, "Can't parse payload. It should be HEX string i.e. \"03 AD F3 41\"")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadRaw(2, "", "", "03 AD F3 41")
	assert.Nil(t, err)
	assert.Equal(t, 2, p.Tnf)
	assert.Equal(t, []byte{0x03, 0xAD, 0xF3, 0x41}, p.Payload)
}

func Test_validateNdefRecordPayloadUrl(t *testing.T) {
	p, err := validateNdefRecordPayloadUrl("")
	assert.EqualError(t, err, "Url flag can't be empty.")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadUrl("any word")
	assert.EqualError(t, err, "Url has wrong value. It must contain http or https and url origin shouldn't be empty.")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadUrl("http://lalka.com")
	assert.Nil(t, err)
	assert.Equal(t, "http://lalka.com", p.Url)
}

func Test_validateNdefRecordPayloadUri(t *testing.T) {
	p, err := validateNdefRecordPayloadUri("")
	assert.EqualError(t, err, "URI value can't be empty.")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadUri("any word")
	assert.Nil(t, err)
	assert.Equal(t, "any word", p.Uri)
}

func Test_validateNdefRecordPayloadTypeText(t *testing.T) {
	p, err := validateNdefRecordPayloadTypeText("", "")
	assert.EqualError(t, err, "Text value can't be empty")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadTypeText("any text", "")
	assert.EqualError(t, err, fmt.Sprintf("Lang value must be one of the following falues: %s", models.NdefLangValues))
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadTypeText("any text", "English")
	assert.Nil(t, err)
	assert.Equal(t, "any text", p.Text)
	assert.Equal(t, "English", p.Lang)
}

func Test_validateNdefTypeVcard(t *testing.T) {
	p, err := validateNdefTypeVcard("1a23", "", "")
	assert.EqualError(t, err, "Flag address-postal-code should contain only digits.")
	assert.Nil(t, p)

	p, err = validateNdefTypeVcard("123", "sadsad", "")
	assert.EqualError(t, err, "Flag email should contain valid email.")
	assert.Nil(t, p)

	p, err = validateNdefTypeVcard("123", "email@email.com", "")
	assert.EqualError(t, err, "Flag first-name can't be empty.")
	assert.Nil(t, p)

	p, err = validateNdefTypeVcard("123", "email@email.com", "FName")
	assert.Nil(t, err)
	assert.Equal(t, "email@email.com", p.Email)
	assert.Equal(t, "123", p.AddressPostalCode)
	assert.Equal(t, "FName", p.FirstName)

	p, err = validateNdefTypeVcard("", "email@email.com", "FName")
	assert.Nil(t, err)
	assert.Equal(t, "email@email.com", p.Email)
	assert.Equal(t, "", p.AddressPostalCode)
	assert.Equal(t, "FName", p.FirstName)

	p, err = validateNdefTypeVcard("", "", "FName")
	assert.Nil(t, err)
	assert.Equal(t, "", p.Email)
	assert.Equal(t, "", p.AddressPostalCode)
	assert.Equal(t, "FName", p.FirstName)
}

func Test_validateNdefRecordPayloadMime(t *testing.T) {
	p, err := validateNdefRecordPayloadMime("", "", "")
	assert.EqualError(t, err, "Type value can't be empty.")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadMime("type", "", "")
	assert.EqualError(t, err, "Format can be either \"hex\" or \"ascii\": Given string is a wrong Mime format")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadMime("type", "hex", "")
	assert.EqualError(t, err, "Content value can't be empty.")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadMime("type", "ascii", "any content")
	assert.Nil(t, err)
	assert.Equal(t, "type", p.Type)
	assert.Equal(t, models.MimeFormatASCII, p.Format)
	assert.Equal(t, "any content", p.ContentASCII)
	assert.Nil(t, p.ContentHEX)

	_, err = validateNdefRecordPayloadMime("type", "hex", "any content")
	assert.Error(t, err)

	p, err = validateNdefRecordPayloadMime("type", "hex", "03 AD F3 41")
	assert.Nil(t, err)
	assert.Equal(t, "type", p.Type)
	assert.Equal(t, models.MimeFormatHex, p.Format)
	assert.Equal(t, "", p.ContentASCII)
	assert.Equal(t, []byte{0x03, 0xAD, 0xF3, 0x41}, p.ContentHEX)
}

func Test_validateNdefRecordPayloadPhone(t *testing.T) {
	p, err := validateNdefRecordPayloadPhone("")
	assert.EqualError(t, err, "Phone number value can't be empty")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadPhone("8 800 555 3535")
	assert.Nil(t, err)
	assert.Equal(t, "8 800 555 3535", p.PhoneNumber)
}

func Test_validateNdefRecordPayloadGeo(t *testing.T) {
	p, err := validateNdefRecordPayloadGeo("", "")
	assert.EqualError(t, err, "Latitude value can't be empty")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadGeo("aaa", "")
	assert.Error(t, err)
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadGeo("-190", "")
	assert.EqualError(t, err, "Wrong latitude value – can be from -90 to 90")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadGeo("190", "")
	assert.EqualError(t, err, "Wrong latitude value – can be from -90 to 90")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadGeo("50", "")
	assert.EqualError(t, err, "Longitude value can't be empty")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadGeo("50", "bbb")
	assert.Error(t, err, "")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadGeo("50", "-190")
	assert.EqualError(t, err, "Wrong latitude value – can be from -180 to 180")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadGeo("50", "190")
	assert.EqualError(t, err, "Wrong latitude value – can be from -180 to 180")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadGeo("50", "170")
	assert.Nil(t, err)
	assert.Equal(t, "50", p.Latitude)
	assert.Equal(t, "170", p.Longitude)

	p, err = validateNdefRecordPayloadGeo("50.2", "170.4")
	assert.Nil(t, err)
	assert.Equal(t, "50.2", p.Latitude)
	assert.Equal(t, "170.4", p.Longitude)

	p, err = validateNdefRecordPayloadGeo("50,2", "170,4")
	assert.Nil(t, err)
	assert.Equal(t, "50.2", p.Latitude)
	assert.Equal(t, "170.4", p.Longitude)
}

func Test_validateNdefRecordPayloadAar(t *testing.T) {
	p, err := validateNdefRecordPayloadAar("")
	assert.EqualError(t, err, "Package name value can't be empty")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadAar("pkg name")
	assert.Nil(t, err)
	assert.Equal(t, "pkg name", p.PackageName)
}

func Test_validateNdefRecordPayloadPoster(t *testing.T) {
	p, err := validateNdefRecordPayloadPoster("", "")
	assert.EqualError(t, err, "Title value can't be empty")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadPoster("title", "")
	assert.EqualError(t, err, "URI value can't be empty")
	assert.Nil(t, p)

	p, err = validateNdefRecordPayloadPoster("title", "any uri")
	assert.Nil(t, err)
	assert.Equal(t, "title", p.Title)
	assert.Equal(t, "any uri", p.Uri)
}
