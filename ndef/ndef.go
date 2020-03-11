package ndef

import (
	"github.com/taglme/nfc-cli/models"
	"github.com/taglme/nfc-goclient/pkg/ndefconv"
)

type NdefPayload interface {
	ToRecord() ndefconv.NdefRecord
}

type NdefRecordPayloadRaw models.NdefRecordPayloadRaw

func (s NdefRecordPayloadRaw) ToRecord() ndefconv.NdefRecord {
	return ndefconv.NdefRecord{
		Type: ndefconv.NdefRecordPayloadTypeRaw,
		Data: ndefconv.NdefRecordPayloadRaw{
			Tnf:     s.Tnf,
			Type:    s.Type,
			ID:      s.ID,
			Payload: s.Payload,
		},
	}
}

type NdefRecordPayloadUrl models.NdefRecordPayloadUrl

func (s NdefRecordPayloadUrl) ToRecord() ndefconv.NdefRecord {
	return ndefconv.NdefRecord{
		Type: ndefconv.NdefRecordPayloadTypeUrl,
		Data: ndefconv.NdefRecordPayloadUrl{
			Url: s.Url,
		},
	}
}

type NdefRecordPayloadTypeText models.NdefRecordPayloadText

func (s NdefRecordPayloadTypeText) ToRecord() ndefconv.NdefRecord {
	return ndefconv.NdefRecord{
		Type: ndefconv.NdefRecordPayloadTypeText,
		Data: ndefconv.NdefRecordPayloadText{
			Text: s.Text,
			Lang: s.Lang,
		},
	}
}

type NdefRecordPayloadUri models.NdefRecordPayloadUri

func (s NdefRecordPayloadUri) ToRecord() ndefconv.NdefRecord {
	return ndefconv.NdefRecord{
		Type: ndefconv.NdefRecordPayloadTypeUri,
		Data: ndefconv.NdefRecordPayloadUri{
			Uri: s.Uri,
		},
	}
}

type NdefRecordPayloadVcard models.NdefRecordPayloadVcard

func (s NdefRecordPayloadVcard) ToRecord() ndefconv.NdefRecord {
	return ndefconv.NdefRecord{
		Type: ndefconv.NdefRecordPayloadTypeVcard,
		Data: ndefconv.NdefRecordPayloadVcard{
			AddressCity:       s.AddressCity,
			AddressCountry:    s.AddressCountry,
			AddressPostalCode: s.AddressPostalCode,
			AddressRegion:     s.AddressRegion,
			AddressStreet:     s.AddressStreet,
			Email:             s.Email,
			FirstName:         s.FirstName,
			LastName:          s.LastName,
			Organization:      s.Organization,
			PhoneCell:         s.PhoneCell,
			PhoneHome:         s.PhoneHome,
			PhoneWork:         s.PhoneWork,
			Title:             s.Title,
			Site:              s.Site,
		},
	}
}

var MapMimeToNdefMime = map[models.MimeFormat]ndefconv.MimeFormat{
	models.MimeFormatASCII: ndefconv.MimeFormatASCII,
	models.MimeFormatHex:   ndefconv.MimeFormatHex,
}

type NdefRecordPayloadMime models.NdefRecordPayloadMime

func (s NdefRecordPayloadMime) ToRecord() ndefconv.NdefRecord {
	return ndefconv.NdefRecord{
		Type: ndefconv.NdefRecordPayloadTypeMime,
		Data: ndefconv.NdefRecordPayloadMime{
			Type:         s.Type,
			Format:       MapMimeToNdefMime[s.Format],
			ContentASCII: s.ContentASCII,
			ContentHEX:   s.ContentHEX,
		},
	}
}

type NdefRecordPayloadPhone models.NdefRecordPayloadPhone

func (s NdefRecordPayloadPhone) ToRecord() ndefconv.NdefRecord {
	return ndefconv.NdefRecord{
		Type: ndefconv.NdefRecordPayloadTypePhone,
		Data: ndefconv.NdefRecordPayloadPhone{
			PhoneNumber: s.PhoneNumber,
		},
	}
}

type NdefRecordPayloadGeo models.NdefRecordPayloadGeo

func (s NdefRecordPayloadGeo) ToRecord() ndefconv.NdefRecord {
	return ndefconv.NdefRecord{
		Type: ndefconv.NdefRecordPayloadTypeGeo,
		Data: ndefconv.NdefRecordPayloadGeo{
			Latitude:  s.Latitude,
			Longitude: s.Longitude,
		},
	}
}

type NdefRecordPayloadAar models.NdefRecordPayloadAar

func (s NdefRecordPayloadAar) ToRecord() ndefconv.NdefRecord {
	return ndefconv.NdefRecord{
		Type: ndefconv.NdefRecordPayloadTypeAar,
		Data: ndefconv.NdefRecordPayloadAar{
			PackageName: s.PackageName,
		},
	}
}

type NdefRecordPayloadPoster models.NdefRecordPayloadPoster

func (s NdefRecordPayloadPoster) ToRecord() ndefconv.NdefRecord {
	return ndefconv.NdefRecord{
		Type: ndefconv.NdefRecordPayloadTypePoster,
		Data: ndefconv.NdefRecordPayloadPoster{
			Title: s.Title,
			Uri:   s.Uri,
		},
	}
}
