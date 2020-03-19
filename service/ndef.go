package service

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-cli/models"
	"github.com/taglme/nfc-cli/ndef"
	"github.com/urfave/cli/v2"
	"regexp"
	"strconv"
	"strings"
)

func (s *appService) parseNdefPayloadFlags(ctx *cli.Context) (res ndef.NdefPayload, err error) {
	ndefType := ctx.String(models.FlagNdefType)

	switch ndefType {
	case models.NdefTypeRaw:
		tnf := ctx.Int(models.FlagNdefTypeRawTnf)
		if tnf < 0 || tnf > 6 {
			return nil, errors.New("Wrong tnf flag value. Can be only from 0 to 6")
		}

		t := ctx.String(models.FlagNdefTypeType)
		id := ctx.String(models.FlagNdefTypeRawId)
		payload := ctx.String(models.FlagNdefTypeRawPayload)
		if len(payload) < 1 {
			return nil, errors.New("Payload value can't be empty")
		}
		p, err := s.parseHexString(payload)
		if err != nil {
			return nil, errors.Wrap(err, "Can't parse payload. It should be HEX string i.e. \"03 AD F3 41\"")
		}

		return ndef.NdefRecordPayloadRaw{
			Tnf:     tnf,
			Type:    t,
			ID:      id,
			Payload: p,
		}, nil
	case models.NdefTypeUrl:
		url := ctx.String(models.FlagNdefTypeUrl)
		matched, err := regexp.MatchString(`http(s)?\:\/\/\w+.*`, url)
		if err != nil {
			return nil, errors.Wrap(err, "Error on the url match string")
		}
		if !matched {
			return nil, errors.New("Url has wrong value. It must contain http or https and url origin shouldn't be empty: ")
		}

		return ndef.NdefRecordPayloadUrl{
			Url: url,
		}, nil
	case models.NdefTypeText:
		text := ctx.String(models.FlagNdefTypeText)
		if len(text) < 1 {
			return nil, errors.New("Text value can't be empty")
		}

		lang := ctx.String(models.FlagNdefTypeLang)
		ok := models.NdefLangValues.Contains(lang)
		if !ok {
			return nil, errors.New(fmt.Sprintf("Lang value must be one of the following falues: %s", models.NdefLangValues))
		}

		return ndef.NdefRecordPayloadTypeText{
			Text: text,
			Lang: lang,
		}, nil
	case models.NdefTypeUri:
		uri := ctx.String(models.FlagNdefUri)
		if len(uri) < 1 {
			return nil, errors.New("URI value can't be empty")
		}

		return ndef.NdefRecordPayloadUri{
			Uri: uri,
		}, nil
	case models.NdefTypeVcard:
		city := ctx.String(models.FlagNdefTypeVcardAddressCity)
		country := ctx.String(models.FlagNdefTypeVcardAddressCountry)
		postal := ctx.String(models.FlagNdefTypeVcardAddressPostalCode)
		if len(postal) > 0 && !regexp.MustCompile(`^\d+$`).MatchString(postal) {
			return nil, errors.New("Flag address-postal-code should contain only digits")
		}

		reg := ctx.String(models.FlagNdefTypeVcardAddressRegion)
		street := ctx.String(models.FlagNdefTypeVcardAddressStreet)
		email := ctx.String(models.FlagNdefTypeVcardEmail)
		if len(email) > 0 && !validateEmail(email) {
			return nil, errors.New("Flag email should contain valid email")
		}

		fName := ctx.String(models.FlagNdefTypeVcardFirstName)
		if len(fName) == 0 {
			return nil, errors.New("Flag first-name can't be empty")
		}

		lName := ctx.String(models.FlagNdefTypeVcardLastName)
		org := ctx.String(models.FlagNdefTypeVcardOrganization)
		cell := ctx.String(models.FlagNdefTypeVcardPhoneCell)
		home := ctx.String(models.FlagNdefTypeVcardPhoneHome)
		work := ctx.String(models.FlagNdefTypeVcardPhoneWork)
		title := ctx.String(models.FlagNdefTypeTitle)
		site := ctx.String(models.FlagNdefTypeVcardSite)

		return ndef.NdefRecordPayloadVcard{
			AddressCity:       city,
			AddressCountry:    country,
			AddressPostalCode: postal,
			AddressRegion:     reg,
			AddressStreet:     street,
			Email:             email,
			FirstName:         fName,
			LastName:          lName,
			Organization:      org,
			PhoneCell:         cell,
			PhoneHome:         home,
			PhoneWork:         work,
			Title:             title,
			Site:              site,
		}, nil
	case models.NdefTypeMime:
		res := ndef.NdefRecordPayloadMime{}

		t := ctx.String(models.FlagNdefTypeType)
		if len(t) < 1 {
			return nil, errors.New("Type value can't be empty")
		}
		res.Type = t

		format := ctx.String(models.FlagNdefTypeMimeFormat)
		mimeFormat, err := models.StringToMimeFormat(format)
		if err != nil {
			return nil, errors.Wrap(err, "Format can be either \"hex\" or \"ascii\"")
		}
		res.Format = mimeFormat

		content := ctx.String(models.FlagNdefTypeMimeContent)
		if len(content) < 1 {
			return nil, errors.New("Content value can't be empty")
		}

		if mimeFormat == models.MimeFormatASCII {
			res.ContentASCII = content
			return res, nil
		}

		c, err := s.parseHexString(content)
		if err != nil {
			return nil, errors.Wrap(err, "Can't parse content string. It should be HEX string i.e. \"03 AD F3 41\"")
		}

		res.ContentHEX = c

		return res, nil
	case models.NdefTypePhone:
		p := ctx.String(models.FlagNdefTypePhone)
		if len(p) < 1 {
			return nil, errors.New("Phone number value can't be empty")
		}

		return ndef.NdefRecordPayloadPhone{
			PhoneNumber: p,
		}, nil
	case models.NdefTypeGeo:
		lat := ctx.String(models.FlagNdefTypeGeoLat)
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

		lon := ctx.String(models.FlagNdefTypeGeoLon)
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

		return ndef.NdefRecordPayloadGeo{
			Latitude:  lat,
			Longitude: lon,
		}, nil
	case models.NdefTypeAar:
		p := ctx.String(models.FlagNdefTypeAarPackage)
		if len(p) < 1 {
			return nil, errors.New("Package name value can't be empty")
		}

		return ndef.NdefRecordPayloadAar{
			PackageName: p,
		}, nil
	case models.NdefTypePoster:
		title := ctx.String(models.FlagNdefTypeTitle)
		if len(title) < 1 {
			return nil, errors.New("Title value can't be empty")
		}

		uri := ctx.String(models.FlagNdefUri)
		if len(uri) < 1 {
			return nil, errors.New("Uri value can't be empty")
		}

		return ndef.NdefRecordPayloadPoster{
			Title: title,
			Uri:   uri,
		}, nil
	}

	return nil, errors.New(fmt.Sprintf("There's no Ndef Record Payload struct for such Ndef Type. Choose one from available: %v", models.NdefTypeValues))
}

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}
