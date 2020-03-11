package models

type Flag = string

const (
	FlagHost    Flag = "host"
	FlagAdapter Flag = "adapter"
	FlagRepeat  Flag = "repeat"
	FlagOutput  Flag = "output"
	FlagAppend  Flag = "append"
	FlagTimeout Flag = "timeout"
	FlagInput   Flag = "input"
	FlagAuth    Flag = "auth"

	FlagPwd Flag = "password"

	FlagTarget  Flag = "target"
	FlagTxBytes Flag = "tx-bytes"

	FlagNdefType Flag = "ndef-type"
	FlagProtect  Flag = "protect"

	FlagNdefTypeRawId      Flag = "id"
	FlagNdefTypeRawTnf     Flag = "tnf"
	FlagNdefTypeType       Flag = "type"
	FlagNdefTypeRawPayload Flag = "payload"

	FlagNdefTypeUrl        Flag = "url"
	FlagNdefTypeText       Flag = "text"
	FlagNdefTypeLang       Flag = "lang"
	FlagNdefUri            Flag = "uri"
	FlagNdefTypeAarPackage Flag = "package-name"
	FlagNdefTypePhone      Flag = "phone-number"

	FlagNdefTypeVcardAddressCity       Flag = "address-city"
	FlagNdefTypeVcardAddressCountry    Flag = "address-country"
	FlagNdefTypeVcardAddressPostalCode Flag = "address-postal-code"
	FlagNdefTypeVcardAddressRegion     Flag = "address-region"
	FlagNdefTypeVcardAddressStreet     Flag = "address-street"
	FlagNdefTypeVcardEmail             Flag = "email"
	FlagNdefTypeVcardFirstName         Flag = "first-name"
	FlagNdefTypeVcardLastName          Flag = "last-name"
	FlagNdefTypeVcardOrganization      Flag = "organization"
	FlagNdefTypeVcardPhoneCell         Flag = "phone-cell"
	FlagNdefTypeVcardPhoneHome         Flag = "phone-home"
	FlagNdefTypeVcardPhoneWork         Flag = "phone-work"
	FlagNdefTypeTitle                  Flag = "title"
	FlagNdefTypeVcardSite              Flag = "site"

	//FlagNdefTypeMimeType Flag = "type"
	FlagNdefTypeMimeFormat  Flag = "format"
	FlagNdefTypeMimeContent Flag = "content"

	FlagNdefTypeGeoLat Flag = "latitude"
	FlagNdefTypeGeoLon Flag = "longitude"

	//FlagNdefTypePosterTitle Flag = "title"
	//FlagNdefTypePosterUri Flag = "uri"
)
