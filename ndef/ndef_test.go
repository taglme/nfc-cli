package ndef

import (
	"github.com/stretchr/testify/assert"
	"github.com/taglme/nfc-cli/models"
	"github.com/taglme/nfc-goclient/pkg/ndefconv"
	"testing"
)

func TestNdefRecordPayloadAar_ToRecord(t *testing.T) {
	a := NdefRecordPayloadAar{PackageName: "package name"}

	assert.Equal(t,
		ndefconv.NdefRecord{
			Type: ndefconv.NdefRecordPayloadTypeAar,
			Data: ndefconv.NdefRecordPayloadAar{PackageName: "package name"},
		},
		a.ToRecord())
}

func TestNdefRecordPayloadGeo_ToRecord(t *testing.T) {
	a := NdefRecordPayloadGeo{
		Latitude:  "25.134",
		Longitude: "43",
	}

	assert.Equal(
		t,
		ndefconv.NdefRecord{
			Type: ndefconv.NdefRecordPayloadTypeGeo,
			Data: ndefconv.NdefRecordPayloadGeo{Latitude: "25.134", Longitude: "43"},
		},
		a.ToRecord(),
	)
}

func TestNdefRecordPayloadMime_ToRecord(t *testing.T) {
	a := NdefRecordPayloadMime{
		Type:         "type",
		Format:       models.MimeFormatASCII,
		ContentASCII: "content",
		ContentHEX:   []byte{0xa1},
	}

	assert.Equal(
		t,
		ndefconv.NdefRecord{
			Type: ndefconv.NdefRecordPayloadTypeMime,
			Data: ndefconv.NdefRecordPayloadMime{
				Type:         "type",
				Format:       ndefconv.MimeFormatASCII,
				ContentASCII: "content",
				ContentHEX:   []byte{0xa1},
			},
		},
		a.ToRecord(),
	)
}

func TestNdefRecordPayloadPhone_ToRecord(t *testing.T) {
	a := NdefRecordPayloadPhone{
		PhoneNumber: "8 800 555 3535",
	}

	assert.Equal(
		t,
		ndefconv.NdefRecord{
			Type: ndefconv.NdefRecordPayloadTypePhone,
			Data: ndefconv.NdefRecordPayloadPhone{
				PhoneNumber: "8 800 555 3535",
			},
		},
		a.ToRecord(),
	)
}

func TestNdefRecordPayloadTypeText_ToRecord(t *testing.T) {
	a := NdefRecordPayloadTypeText{
		Text: "Text?",
		Lang: "English",
	}

	assert.Equal(
		t,
		ndefconv.NdefRecord{
			Type: ndefconv.NdefRecordPayloadTypeText,
			Data: ndefconv.NdefRecordPayloadText{
				Text: "Text?",
				Lang: "English",
			},
		},
		a.ToRecord(),
	)
}

func TestNdefRecordPayloadPoster_ToRecord(t *testing.T) {
	a := NdefRecordPayloadPoster{
		Title: "title",
		Uri:   "http://url.com",
	}

	assert.Equal(
		t,
		ndefconv.NdefRecord{
			Type: ndefconv.NdefRecordPayloadTypePoster,
			Data: ndefconv.NdefRecordPayloadPoster{
				Title: "title",
				Uri:   "http://url.com",
			},
		},
		a.ToRecord(),
	)
}

func TestNdefRecordPayloadRaw_ToRecord(t *testing.T) {
	a := NdefRecordPayloadRaw{
		Tnf:     2,
		Type:    "type",
		ID:      "id",
		Payload: []byte{0xa6},
	}

	assert.Equal(
		t,
		ndefconv.NdefRecord{
			Type: ndefconv.NdefRecordPayloadTypeRaw,
			Data: ndefconv.NdefRecordPayloadRaw{
				Tnf:     2,
				Type:    "type",
				ID:      "id",
				Payload: []byte{0xa6},
			},
		},
		a.ToRecord(),
	)
}

func TestNdefRecordPayloadUri_ToRecord(t *testing.T) {
	a := NdefRecordPayloadUri{
		Uri: "http://url.com",
	}

	assert.Equal(
		t,
		ndefconv.NdefRecord{
			Type: ndefconv.NdefRecordPayloadTypeUri,
			Data: ndefconv.NdefRecordPayloadUri{
				Uri: "http://url.com",
			},
		},
		a.ToRecord(),
	)
}

func TestNdefRecordPayloadUrl_ToRecord(t *testing.T) {
	a := NdefRecordPayloadUrl{
		Url: "http://url.com",
	}

	assert.Equal(
		t,
		ndefconv.NdefRecord{
			Type: ndefconv.NdefRecordPayloadTypeUrl,
			Data: ndefconv.NdefRecordPayloadUrl{
				Url: "http://url.com",
			},
		},
		a.ToRecord(),
	)
}

func TestNdefRecordPayloadVcard_ToRecord(t *testing.T) {
	a := NdefRecordPayloadVcard{
		AddressCity:       "any",
		AddressCountry:    "any",
		AddressPostalCode: "123",
		AddressRegion:     "any",
		AddressStreet:     "any",
		Email:             "any@gmail.com",
		FirstName:         "any",
		LastName:          "any",
		Organization:      "any",
		PhoneCell:         "any",
		PhoneHome:         "any",
		PhoneWork:         "any",
		Title:             "any",
		Site:              "any",
	}

	assert.Equal(
		t,
		ndefconv.NdefRecord{
			Type: ndefconv.NdefRecordPayloadTypeVcard,
			Data: ndefconv.NdefRecordPayloadVcard{
				AddressCity:       "any",
				AddressCountry:    "any",
				AddressPostalCode: "123",
				AddressRegion:     "any",
				AddressStreet:     "any",
				Email:             "any@gmail.com",
				FirstName:         "any",
				LastName:          "any",
				Organization:      "any",
				PhoneCell:         "any",
				PhoneHome:         "any",
				PhoneWork:         "any",
				Title:             "any",
				Site:              "any",
			},
		},
		a.ToRecord(),
	)
}
