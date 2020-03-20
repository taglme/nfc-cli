package service

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-cli/models"
	"github.com/taglme/nfc-cli/ndef"
	"github.com/urfave/cli/v2"
)

func (s *appService) parseNdefPayloadFlags(ctx *cli.Context) (res ndef.NdefPayload, err error) {
	ndefType := ctx.String(models.FlagNdefType)

	switch ndefType {
	case models.NdefTypeRaw:
		tnf := ctx.Int(models.FlagNdefTypeRawTnf)
		t := ctx.String(models.FlagNdefTypeType)
		id := ctx.String(models.FlagNdefTypeRawId)
		payload := ctx.String(models.FlagNdefTypeRawPayload)

		return validateNdefRecordPayloadRaw(tnf, t, id, payload)
	case models.NdefTypeUrl:
		url := ctx.String(models.FlagNdefTypeUrl)
		return validateNdefRecordPayloadUrl(url)
	case models.NdefTypeText:
		text := ctx.String(models.FlagNdefTypeText)
		lang := ctx.String(models.FlagNdefTypeLang)
		return validateNdefRecordPayloadTypeText(text, lang)
	case models.NdefTypeUri:
		uri := ctx.String(models.FlagNdefUri)
		return validateNdefRecordPayloadUri(uri)
	case models.NdefTypeVcard:
		postal := ctx.String(models.FlagNdefTypeVcardAddressPostalCode)
		email := ctx.String(models.FlagNdefTypeVcardEmail)
		fName := ctx.String(models.FlagNdefTypeVcardFirstName)
		res, err := validateNdefTypeVcard(postal, email, fName)
		if err != nil {
			return nil, err
		}

		res.AddressCity = ctx.String(models.FlagNdefTypeVcardAddressCity)
		res.AddressCountry = ctx.String(models.FlagNdefTypeVcardAddressCountry)
		res.AddressRegion = ctx.String(models.FlagNdefTypeVcardAddressRegion)
		res.AddressStreet = ctx.String(models.FlagNdefTypeVcardAddressStreet)
		res.LastName = ctx.String(models.FlagNdefTypeVcardLastName)
		res.Organization = ctx.String(models.FlagNdefTypeVcardOrganization)
		res.PhoneCell = ctx.String(models.FlagNdefTypeVcardPhoneCell)
		res.PhoneHome = ctx.String(models.FlagNdefTypeVcardPhoneHome)
		res.PhoneWork = ctx.String(models.FlagNdefTypeVcardPhoneWork)
		res.Title = ctx.String(models.FlagNdefTypeTitle)
		res.Site = ctx.String(models.FlagNdefTypeVcardSite)

		return res, nil
	case models.NdefTypeMime:
		t := ctx.String(models.FlagNdefTypeType)
		format := ctx.String(models.FlagNdefTypeMimeFormat)
		content := ctx.String(models.FlagNdefTypeMimeContent)
		return validateNdefRecordPayloadMime(t, format, content)
	case models.NdefTypePhone:
		p := ctx.String(models.FlagNdefTypePhone)
		return validateNdefRecordPayloadPhone(p)
	case models.NdefTypeGeo:
		lat := ctx.String(models.FlagNdefTypeGeoLat)
		lon := ctx.String(models.FlagNdefTypeGeoLon)
		return validateNdefRecordPayloadGeo(lat, lon)
	case models.NdefTypeAar:
		p := ctx.String(models.FlagNdefTypeAarPackage)
		return validateNdefRecordPayloadAar(p)
	case models.NdefTypePoster:
		title := ctx.String(models.FlagNdefTypeTitle)
		uri := ctx.String(models.FlagNdefUri)
		return validateNdefRecordPayloadPoster(title, uri)
	}

	return nil, errors.New(fmt.Sprintf("There's no Ndef Record Payload struct for such Ndef Type. Choose one from available: %v", models.NdefTypeValues))
}
