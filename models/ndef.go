package models

import "github.com/pkg/errors"

type NdefType = string

const (
	NdefTypeRaw    NdefType = "raw"
	NdefTypeUrl    NdefType = "url"
	NdefTypeText   NdefType = "text"
	NdefTypeUri    NdefType = "uri"
	NdefTypeVcard  NdefType = "vcard"
	NdefTypeMime   NdefType = "mime"
	NdefTypePhone  NdefType = "phone"
	NdefTypeGeo    NdefType = "geo"
	NdefTypeAar    NdefType = "aar"
	NdefTypePoster NdefType = "poster"
)

var NdefTypeValues = []NdefType{
	NdefTypeRaw,
	NdefTypeUrl,
	NdefTypeText,
	NdefTypeUri,
	NdefTypeVcard,
	NdefTypeMime,
	NdefTypePhone,
	NdefTypeGeo,
	NdefTypeAar,
	NdefTypePoster,
}

type NdefPayload interface{}

type NdefRecordPayloadRaw struct {
	Tnf     int
	Type    string
	ID      string
	Payload []byte
}
type NdefLangValuesSlice []string

var NdefLangValues = NdefLangValuesSlice{
	"Arabic", "Bengali", "Chinese", "Danish", "Dutch", "English", "Finnish", "French", "German", "Greek", "Hebrew", "Hindi", "Irish", "Italian", "Japanese", "Latin", "Portuguese", "Russian", "Spanish",
}

func (arr NdefLangValuesSlice) Contains(str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

type NdefRecordPayloadText struct {
	Text string
	Lang string
}

type NdefRecordPayloadUrl struct {
	Url string
}

type NdefRecordPayloadUri struct {
	Uri string
}
type MimeFormat string

const (
	MimeFormatASCII MimeFormat = "ascii"
	MimeFormatHex   MimeFormat = "hex"
)

func StringToMimeFormat(s string) (MimeFormat, error) {
	switch s {
	case "ascii":
		return MimeFormatASCII, nil
	case "hex":
		return MimeFormatHex, nil
	}

	return "unknown", errors.New("Given string is a wrong Mime format")
}

type NdefRecordPayloadMime struct {
	Type         string
	Format       MimeFormat
	ContentASCII string
	ContentHEX   []byte
}

type NdefRecordPayloadPhone struct {
	PhoneNumber string
}

type NdefRecordPayloadGeo struct {
	Latitude  string
	Longitude string
}
type NdefRecordPayloadAar struct {
	PackageName string
}
type NdefRecordPayloadPoster struct {
	Title string
	Uri   string
}
type NdefRecordPayloadVcard struct {
	AddressCity       string
	AddressCountry    string
	AddressPostalCode string
	AddressRegion     string
	AddressStreet     string
	Email             string
	FirstName         string
	LastName          string
	Organization      string
	PhoneCell         string
	PhoneHome         string
	PhoneWork         string
	Title             string
	Site              string
}
