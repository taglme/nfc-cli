package service

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-cli/models"
	"github.com/taglme/nfc-cli/ndef"
	"github.com/taglme/nfc-cli/utils"
	"regexp"
	"strconv"
	"strings"
)

func validateNdefRecordPayloadRaw(tnf int, t, id, payload string) (*ndef.NdefRecordPayloadRaw, error) {
	if tnf < 0 || tnf > 6 {
		return nil, errors.New("Wrong tnf flag value. Can be only from 0 to 6")
	}
	if len(payload) < 1 {
		return nil, errors.New("Payload value can't be empty")
	}
	p, err := utils.ParseHexString(payload)
	if err != nil {
		return nil, errors.Wrap(err, "Can't parse payload. It should be HEX string i.e. \"03 AD F3 41\"")
	}

	return &ndef.NdefRecordPayloadRaw{
		Tnf:     tnf,
		Type:    t,
		ID:      id,
		Payload: p,
	}, nil
}

func validateNdefRecordPayloadUrl(url string) (*ndef.NdefRecordPayloadUrl, error) {
	if len(url) == 0 {
		return nil, errors.New("Url flag can't be empty.")
	}

	matched, err := regexp.MatchString(`http(s)?\:\/\/\w+.*`, url)
	if err != nil {
		return nil, errors.Wrap(err, "Error on the url match string")
	}
	if !matched {
		return nil, errors.New("Url has wrong value. It must contain http or https and url origin shouldn't be empty.")
	}

	return &ndef.NdefRecordPayloadUrl{
		Url: url,
	}, nil
}

func validateNdefRecordPayloadUri(uri string) (*ndef.NdefRecordPayloadUri, error) {
	if len(uri) < 1 {
		return nil, errors.New("URI value can't be empty.")
	}

	return &ndef.NdefRecordPayloadUri{
		Uri: uri,
	}, nil
}

func validateNdefRecordPayloadTypeText(text, lang string) (*ndef.NdefRecordPayloadTypeText, error) {
	if len(text) < 1 {
		return nil, errors.New("Text value can't be empty")
	}

	ok := models.NdefLangValues.Contains(lang)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Lang value must be one of the following falues: %s", models.NdefLangValues))
	}

	return &ndef.NdefRecordPayloadTypeText{
		Text: text,
		Lang: lang,
	}, nil
}

func validateNdefTypeVcard(postal, email, fName string) (*ndef.NdefRecordPayloadVcard, error) {
	if len(postal) > 0 && !regexp.MustCompile(`^\d+$`).MatchString(postal) {
		return nil, errors.New("Flag address-postal-code should contain only digits.")
	}

	if len(email) > 0 && !utils.ValidateEmail(email) {
		return nil, errors.New("Flag email should contain valid email.")
	}

	if len(fName) == 0 {
		return nil, errors.New("Flag first-name can't be empty.")
	}

	return &ndef.NdefRecordPayloadVcard{

		Email:             email,
		AddressPostalCode: postal,
		FirstName:         fName,
	}, nil
}

func validateNdefRecordPayloadMime(t, format, content string) (*ndef.NdefRecordPayloadMime, error) {
	if len(t) < 1 {
		return nil, errors.New("Type value can't be empty.")
	}
	res := ndef.NdefRecordPayloadMime{}
	res.Type = t

	mimeFormat, err := models.StringToMimeFormat(format)
	if err != nil {
		return nil, errors.Wrap(err, "Format can be either \"hex\" or \"ascii\"")
	}
	res.Format = mimeFormat

	if len(content) < 1 {
		return nil, errors.New("Content value can't be empty.")
	}

	if mimeFormat == models.MimeFormatASCII {
		res.ContentASCII = content
		return &res, nil
	}

	c, err := utils.ParseHexString(content)
	if err != nil {
		return nil, errors.Wrap(err, "Can't parse content string. It should be HEX string i.e. \"03 AD F3 41\"")
	}

	res.ContentHEX = c

	return &res, nil
}

func validateNdefRecordPayloadPhone(p string) (*ndef.NdefRecordPayloadPhone, error) {
	if len(p) < 1 {
		return nil, errors.New("Phone number value can't be empty")
	}

	return &ndef.NdefRecordPayloadPhone{
		PhoneNumber: p,
	}, nil
}

func validateNdefRecordPayloadGeo(lat, lon string) (*ndef.NdefRecordPayloadGeo, error) {
	if len(lat) < 1 {
		return nil, errors.New("Latitude value can't be empty")
	}
	lat = strings.Replace(lat, ",", ".", -1)
	latF, err := strconv.ParseFloat(lat, 16)
	if err != nil {
		return nil, errors.Wrap(err, "Error on parsing float from latitude string")
	}

	if latF < -90 || latF > 90 {
		return nil, errors.New("Wrong latitude value – can be from -90 to 90")
	}

	if len(lon) < 1 {
		return nil, errors.New("Longitude value can't be empty")
	}
	lon = strings.Replace(lon, ",", ".", -1)
	lonF, err := strconv.ParseFloat(lon, 16)
	if err != nil {
		return nil, errors.Wrap(err, "Error on parsing float from longitude string")
	}

	if lonF < -180 || lonF > 180 {
		return nil, errors.New("Wrong latitude value – can be from -180 to 180")
	}

	return &ndef.NdefRecordPayloadGeo{
		Latitude:  lat,
		Longitude: lon,
	}, nil
}

func validateNdefRecordPayloadAar(p string) (*ndef.NdefRecordPayloadAar, error) {
	if len(p) < 1 {
		return nil, errors.New("Package name value can't be empty")
	}

	return &ndef.NdefRecordPayloadAar{
		PackageName: p,
	}, nil
}

func validateNdefRecordPayloadPoster(title, uri string) (*ndef.NdefRecordPayloadPoster, error) {
	if len(title) < 1 {
		return nil, errors.New("Title value can't be empty")
	}

	if len(uri) < 1 {
		return nil, errors.New("URI value can't be empty")
	}

	return &ndef.NdefRecordPayloadPoster{
		Title: title,
		Uri:   uri,
	}, nil
}
