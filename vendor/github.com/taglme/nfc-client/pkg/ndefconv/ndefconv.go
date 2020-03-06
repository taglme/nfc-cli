package ndefconv

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/emersion/go-vcard"
	"github.com/hsanjuan/go-ndef"
	"github.com/hsanjuan/go-ndef/types/generic"
	"github.com/hsanjuan/go-ndef/types/wkt/text"
)

type Ndef struct {
	ReadOnly bool
	Message  []NdefRecord
}
type NdefResource struct {
	ReadOnly bool                 `json:"read_only"`
	Message  []NdefRecordResource `json:"message"`
}

func (ndef *Ndef) String() string {
	str := ""
	last := len(ndef.Message) - 1
	for i, r := range ndef.Message {
		str += r.String()
		if i != last {
			str += "\n"
		}
	}
	return str

}

func (ndef Ndef) ToResource() NdefResource {
	var ndefRecordResources []NdefRecordResource
	for _, ndefRecord := range ndef.Message {
		ndefRecordResources = append(ndefRecordResources, ndefRecord.ToResource())
	}
	resource := NdefResource{
		ReadOnly: ndef.ReadOnly,
		Message:  ndefRecordResources,
	}

	return resource
}
func (ndefResource NdefResource) ToNdefRecord() (Ndef, error) {
	var ndefRecords []NdefRecord
	for _, ndefRecordResource := range ndefResource.Message {
		ndefRecord, err := ndefRecordResource.ToNdefRecord()
		if err != nil {
			return Ndef{}, err
		}
		ndefRecords = append(ndefRecords, ndefRecord)
	}
	resource := Ndef{
		ReadOnly: ndefResource.ReadOnly,
		Message:  ndefRecords,
	}
	return resource, nil
}

type NdefRecord struct {
	Type NdefRecordPayloadType
	Data NdefRecordPayload
}

func (ndefRecord *NdefRecord) String() string {
	return ndefRecord.Data.String()
}

type NdefRecordResource struct {
	Type string                    `json:"type"`
	Data NdefRecordPayloadResource `json:"data"`
}

func (ndefRecord *NdefRecordResource) UnmarshalJSON(data []byte) error {

	var obj map[string]interface{}
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}

	t, ok := obj["type"].(string)
	if !ok {
		return errors.New("Ndef record should have 'type' field")
	}

	recordType, isValid := StringToNdefRecordPayloadType(t)
	if !isValid {
		return errors.New("Ndef record have not valid type")
	}
	ndefRecord.Type = t

	_, ok = obj["data"]

	if !ok {
		return errors.New("Ndef record  should have 'data' field")
	}

	var dataBytes []byte
	dataBytes, _ = json.Marshal(obj["data"])
	switch recordType {
	case NdefRecordPayloadTypeRaw:
		r := NdefRecordPayloadRawResource{}
		err := json.Unmarshal(dataBytes, &r)
		if err != nil {
			return err
		}
		if r.Tnf < 0 || r.Tnf > 6 {
			return errors.New("Tnf field of Raw type record should have value from '0' to '6'")
		}
		//Tnf == 0 is empty record, so we don't need payload
		if r.Tnf > 0 && r.Payload == "" {
			return errors.New("Payload field of Raw type record should be not empty")
		}
		ndefRecord.Data = r
	case NdefRecordPayloadTypeUrl:
		r := NdefRecordPayloadUrlResource{}
		err := json.Unmarshal(dataBytes, &r)
		if err != nil {
			return err
		}
		if r.Url == "" {
			return errors.New("Url field of Url type record should be not empty")
		}
		ndefRecord.Data = r
	case NdefRecordPayloadTypeText:
		r := NdefRecordPayloadTextResource{}
		err := json.Unmarshal(dataBytes, &r)
		if err != nil {
			return err
		}
		if r.Text == "" {
			return errors.New("Text field of Text type record should be not empty")
		}
		ndefRecord.Data = r
	case NdefRecordPayloadTypeUri:
		r := NdefRecordPayloadUriResource{}
		err := json.Unmarshal(dataBytes, &r)
		if err != nil {
			return err
		}
		if r.Uri == "" {
			return errors.New("Uri field of Uri type record should be not empty")
		}
		ndefRecord.Data = r
	case NdefRecordPayloadTypeVcard:
		r := NdefRecordPayloadVcardResource{}
		err := json.Unmarshal(dataBytes, &r)
		if err != nil {
			return err
		}
		if r.FirstName == "" {
			return errors.New("First name field of Vcard type record should be not empty")
		}
		ndefRecord.Data = r
	case NdefRecordPayloadTypeMime:
		r := NdefRecordPayloadMimeResource{}
		err := json.Unmarshal(dataBytes, &r)
		if err != nil {
			return err
		}
		if r.Type == "" {
			return errors.New("Type field of Mime type record should be not empty")
		}
		_, isValidType := StringToMimeFormat(r.Format)
		if !isValidType {
			return errors.New("Format field of Mime type record has invalid value")
		}

		ndefRecord.Data = r
	case NdefRecordPayloadTypePhone:
		r := NdefRecordPayloadPhoneResource{}
		err := json.Unmarshal(dataBytes, &r)
		if err != nil {
			return err
		}
		if r.PhoneNumber == "" {
			return errors.New("Phone number field of Phone type record should be not empty")
		}

		ndefRecord.Data = r
	case NdefRecordPayloadTypeGeo:
		r := NdefRecordPayloadGeoResource{}
		err := json.Unmarshal(dataBytes, &r)
		if err != nil {
			return err
		}
		if r.Latitude == "" {
			return errors.New("Latitude field of Geo type record should be not empty")
		}
		if r.Longitude == "" {
			return errors.New("Longitude field of Geo type record should be not empty")
		}

		ndefRecord.Data = r
	case NdefRecordPayloadTypeAar:
		r := NdefRecordPayloadAarResource{}
		err := json.Unmarshal(dataBytes, &r)
		if err != nil {
			return err
		}
		if r.PackageName == "" {
			return errors.New("Package name field of Android application type record should be not empty")
		}
		ndefRecord.Data = r
	case NdefRecordPayloadTypePoster:
		r := NdefRecordPayloadPosterResource{}
		err := json.Unmarshal(dataBytes, &r)
		if err != nil {
			return err
		}
		if r.Title == "" {
			return errors.New("Title name field of Smartposter type record should be not empty")
		}
		if r.Uri == "" {
			return errors.New("Uri name field of Smartposter type record should be not empty")
		}
		ndefRecord.Data = r
	}
	return nil
}

func (ndefRecord NdefRecord) ToResource() NdefRecordResource {
	resource := NdefRecordResource{
		Type: ndefRecord.Type.String(),
		Data: ndefRecord.Data.ToResource(),
	}

	return resource
}
func (ndefRecordResource NdefRecordResource) ToNdefRecord() (NdefRecord, error) {
	ndefRecordPayloadType, _ := StringToNdefRecordPayloadType(ndefRecordResource.Type)
	data, err := ndefRecordResource.Data.ToPayload()
	if err != nil {
		return NdefRecord{}, err
	}

	resource := NdefRecord{
		Type: ndefRecordPayloadType,
		Data: data,
	}
	return resource, nil
}

type NdefRecordPayload interface {
	ToResource() NdefRecordPayloadResource
	ToRecord() *ndef.Record
	String() string
}
type NdefRecordPayloadResource interface {
	ToPayload() (NdefRecordPayload, error)
}

type NdefRecordPayloadType int

const (
	NdefRecordPayloadTypeRaw NdefRecordPayloadType = iota + 1
	NdefRecordPayloadTypeUrl
	NdefRecordPayloadTypeText
	NdefRecordPayloadTypeUri
	NdefRecordPayloadTypeVcard
	NdefRecordPayloadTypeMime
	NdefRecordPayloadTypePhone
	NdefRecordPayloadTypeGeo
	NdefRecordPayloadTypeAar
	NdefRecordPayloadTypePoster
)

func StringToNdefRecordPayloadType(s string) (NdefRecordPayloadType, bool) {
	switch s {
	case NdefRecordPayloadTypeRaw.String():
		return NdefRecordPayloadTypeRaw, true
	case NdefRecordPayloadTypeUrl.String():
		return NdefRecordPayloadTypeUrl, true
	case NdefRecordPayloadTypeText.String():
		return NdefRecordPayloadTypeText, true
	case NdefRecordPayloadTypeUri.String():
		return NdefRecordPayloadTypeUri, true
	case NdefRecordPayloadTypeVcard.String():
		return NdefRecordPayloadTypeVcard, true
	case NdefRecordPayloadTypeMime.String():
		return NdefRecordPayloadTypeMime, true
	case NdefRecordPayloadTypePhone.String():
		return NdefRecordPayloadTypePhone, true
	case NdefRecordPayloadTypeGeo.String():
		return NdefRecordPayloadTypeGeo, true
	case NdefRecordPayloadTypeAar.String():
		return NdefRecordPayloadTypeAar, true
	case NdefRecordPayloadTypePoster.String():
		return NdefRecordPayloadTypePoster, true
	}
	return 0, false
}

func (ndefRecordPayloadType NdefRecordPayloadType) String() string {
	names := [...]string{
		"unknown",
		"raw",
		"url",
		"text",
		"uri",
		"vcard",
		"mime",
		"phone",
		"geo",
		"aar",
		"poster",
	}
	if ndefRecordPayloadType < NdefRecordPayloadTypeRaw || ndefRecordPayloadType > NdefRecordPayloadTypePoster {
		return names[0]
	}
	return names[ndefRecordPayloadType]
}

type NdefRecordPayloadRaw struct {
	Tnf     int
	Type    string
	ID      string
	Payload []byte
}
type NdefRecordPayloadRawResource struct {
	Tnf     int    `json:"tnf"`
	Type    string `json:"type"`
	ID      string `json:"id"`
	Payload string `json:"payload"`
}

func (ndefRecordPayload NdefRecordPayloadRaw) ToResource() NdefRecordPayloadResource {
	encodedString := base64.StdEncoding.EncodeToString(ndefRecordPayload.Payload)
	resource := NdefRecordPayloadRawResource{
		Tnf:     ndefRecordPayload.Tnf,
		Type:    ndefRecordPayload.Type,
		ID:      ndefRecordPayload.ID,
		Payload: encodedString,
	}
	return resource
}
func (ndefRecordPayload NdefRecordPayloadRaw) ToRecord() *ndef.Record {

	payload := generic.New(ndefRecordPayload.Payload)
	record := ndef.NewRecord(byte(ndefRecordPayload.Tnf), ndefRecordPayload.Type, ndefRecordPayload.ID, payload)

	return record
}

func (ndefRecordPayload NdefRecordPayloadRaw) String() string {
	return fmt.Sprintf("%s, %s, % x", TnfToString(ndefRecordPayload.Tnf), ndefRecordPayload.Type, ndefRecordPayload.Payload)
}

func (ndefRecordPayloadResource NdefRecordPayloadRawResource) ToPayload() (NdefRecordPayload, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(ndefRecordPayloadResource.Payload)
	if err != nil {
		fmt.Println("decode error:", err)
		return NdefRecordPayloadRaw{}, errors.New("Could not decode payload. It should be base64 encoded")
	}

	ndefRecordPayload := NdefRecordPayloadRaw{
		Tnf:     ndefRecordPayloadResource.Tnf,
		Type:    ndefRecordPayloadResource.Type,
		ID:      ndefRecordPayloadResource.ID,
		Payload: decodedBytes,
	}
	return ndefRecordPayload, nil
}

type NdefRecordPayloadUrl struct {
	Url string
}
type NdefRecordPayloadUrlResource struct {
	Url string `json:"url"`
}

func (ndefRecordPayload NdefRecordPayloadUrl) ToResource() NdefRecordPayloadResource {
	return NdefRecordPayloadUrlResource(ndefRecordPayload)
}

func (ndefRecordPayload NdefRecordPayloadUrl) ToRecord() *ndef.Record {
	record := ndef.NewURIRecord(ndefRecordPayload.Url)
	return record
}

func (ndefRecordPayload NdefRecordPayloadUrl) String() string {
	return ndefRecordPayload.Url
}

func (ndefRecordPayloadResource NdefRecordPayloadUrlResource) ToPayload() (NdefRecordPayload, error) {
	return NdefRecordPayloadUrl(ndefRecordPayloadResource), nil
}

type NdefRecordPayloadText struct {
	Text string
	Lang string
}
type NdefRecordPayloadTextResource struct {
	Text string `json:"text"`
	Lang string `json:"lang"`
}

func (ndefRecordPayload NdefRecordPayloadText) ToResource() NdefRecordPayloadResource {
	return NdefRecordPayloadTextResource(ndefRecordPayload)
}

func (ndefRecordPayload NdefRecordPayloadText) ToRecord() *ndef.Record {
	//replace \n with {0x0D, 0x0A} according to spec
	//todo: 1. add escape check for \\n 2. collapse spaces
	crlf := string([]byte{0x0D, 0x0A})
	text := strings.Replace(ndefRecordPayload.Text, "\n", crlf, -1)
	record := ndef.NewTextRecord(text, LangToCode(ndefRecordPayload.Lang))
	return record
}

func (ndefRecordPayload NdefRecordPayloadText) String() string {

	return ndefRecordPayload.Text
}

func (ndefRecordPayloadResource NdefRecordPayloadTextResource) ToPayload() (NdefRecordPayload, error) {
	return NdefRecordPayloadText(ndefRecordPayloadResource), nil
}

func LangToCode(lang string) string {
	var code string
	switch lang {
	case "Arabic":
		code = "ar"
	case "Bengali":
		code = "bn"
	case "Chinese":
		code = "zh"
	case "Danish":
		code = "da"
	case "Dutch":
		code = "nl"
	case "English":
		code = "en"
	case "Finnish":
		code = "fi"
	case "French":
		code = "fr"
	case "German":
		code = "de"
	case "Greek":
		code = "el"
	case "Hebrew":
		code = "he"
	case "Hindi":
		code = "hi"
	case "Irish":
		code = "ga"
	case "Italian":
		code = "it"
	case "Japanese":
		code = "ja"
	case "Latin":
		code = "la"
	case "Portuguese":
		code = "pt"
	case "Russian":
		code = "ru"
	case "Spanish":
		code = "es"
	default:
		code = "en"
	}
	return code
}

func CodeToLang(code string) string {
	var lang string
	switch code {
	case "ar":
		lang = "Arabic"
	case "bn":
		lang = "Bengali"
	case "zh":
		lang = "Chinese"
	case "da":
		lang = "Danish"
	case "nl":
		lang = "Dutch"
	case "en":
		lang = "English"
	case "fi":
		lang = "Finnish"
	case "fr":
		lang = "French"
	case "de":
		lang = "German"
	case "el":
		lang = "Greek"
	case "he":
		lang = "Hebrew"
	case "hi":
		lang = "Hindi"
	case "ga":
		lang = "Irish"
	case "it":
		lang = "Italian"
	case "ja":
		lang = "Japanese"
	case "la":
		lang = "Latin"
	case "pt":
		lang = "Portuguese"
	case "ru":
		lang = "Russian"
	case "es":
		lang = "Spanish"
	default:
		lang = "English"
	}
	return lang
}

type NdefRecordPayloadUri struct {
	Uri string
}
type NdefRecordPayloadUriResource struct {
	Uri string `json:"uri"`
}

func (ndefRecordPayload NdefRecordPayloadUri) ToResource() NdefRecordPayloadResource {
	return NdefRecordPayloadUriResource(ndefRecordPayload)
}

func (ndefRecordPayload NdefRecordPayloadUri) ToRecord() *ndef.Record {
	record := ndef.NewURIRecord(ndefRecordPayload.Uri)
	return record
}

func (ndefRecordPayload NdefRecordPayloadUri) String() string {
	return ndefRecordPayload.Uri
}

func (ndefRecordPayloadResource NdefRecordPayloadUriResource) ToPayload() (NdefRecordPayload, error) {
	return NdefRecordPayloadUri(ndefRecordPayloadResource), nil
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
type NdefRecordPayloadVcardResource struct {
	AddressCity       string `json:"address_city"`
	AddressCountry    string `json:"address_country"`
	AddressPostalCode string `json:"address_postal_code"`
	AddressRegion     string `json:"address_region"`
	AddressStreet     string `json:"address_street"`
	Email             string `json:"email"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Organization      string `json:"organization"`
	PhoneCell         string `json:"phone_cell"`
	PhoneHome         string `json:"phone_home"`
	PhoneWork         string `json:"phone_work"`
	Title             string `json:"title"`
	Site              string `json:"site"`
}

func (ndefRecordPayload NdefRecordPayloadVcard) ToResource() NdefRecordPayloadResource {
	return NdefRecordPayloadVcardResource(ndefRecordPayload)
}

func (ndefRecordPayload NdefRecordPayloadVcard) ToRecord() *ndef.Record {
	var addresses []string

	if ndefRecordPayload.AddressStreet != "" {
		addresses = append(addresses, ndefRecordPayload.AddressStreet)
	}
	if ndefRecordPayload.AddressCity != "" {
		addresses = append(addresses, ndefRecordPayload.AddressCity)
	}
	if ndefRecordPayload.AddressRegion != "" {
		addresses = append(addresses, ndefRecordPayload.AddressRegion)
	}
	if ndefRecordPayload.AddressPostalCode != "" {
		addresses = append(addresses, ndefRecordPayload.AddressPostalCode)
	}
	if ndefRecordPayload.AddressCountry != "" {
		addresses = append(addresses, ndefRecordPayload.AddressCountry)
	}

	var vcardString = "BEGIN:VCARD\n" + "VERSION:3.0\n"

	if ndefRecordPayload.FirstName != "" && ndefRecordPayload.LastName != "" {
		vcardString += "N:" + ndefRecordPayload.LastName + ";" + ndefRecordPayload.FirstName + "\n"
		vcardString += "FN:" + ndefRecordPayload.FirstName + " " + ndefRecordPayload.LastName + "\n"
	} else if ndefRecordPayload.LastName == "" {
		vcardString += "N:" + ndefRecordPayload.FirstName + "\n"
		vcardString += "FN:" + ndefRecordPayload.FirstName + "\n"
	} else if ndefRecordPayload.FirstName == "" {
		vcardString += "N:" + ndefRecordPayload.LastName + "\n"
		vcardString += "FN:" + ndefRecordPayload.LastName + "\n"
	}

	if ndefRecordPayload.Title != "" {
		vcardString += "TITLE:" + ndefRecordPayload.Title + "\n"
	}

	if ndefRecordPayload.Organization != "" {
		vcardString += "ORG:" + ndefRecordPayload.Organization + "\n"
	}

	if len(addresses) > 0 {
		vcardString += "ADR:;;"

		for i := range addresses {
			if i == (len(addresses) - 1) {
				vcardString += addresses[i] + ";\n"
			} else {
				vcardString += addresses[i] + ";"
			}

		}

	}

	if ndefRecordPayload.PhoneHome != "" {
		vcardString += "TEL;TYPE=HOME,VOICE:" + ndefRecordPayload.PhoneHome + "\n"
	}

	if ndefRecordPayload.PhoneWork != "" {
		vcardString += "TEL;TYPE=WORK,VOICE:" + ndefRecordPayload.PhoneWork + "\n"
	}

	if ndefRecordPayload.PhoneCell != "" {
		vcardString += "TEL;TYPE=CELL:" + ndefRecordPayload.PhoneCell + "\n"
	}

	if ndefRecordPayload.Email != "" {
		vcardString += "EMAIL:" + ndefRecordPayload.Email + "\n"
	}
	if ndefRecordPayload.Site != "" {
		vcardString += "URL:" + ndefRecordPayload.Site + "\n"
	}

	vcardString += "END:VCARD"

	record := ndef.NewMediaRecord("text/vcard", []byte(vcardString))
	return record
}
func (ndefRecordPayload NdefRecordPayloadVcard) String() string {
	s := ndefRecordPayload.FirstName
	if ndefRecordPayload.LastName != "" && ndefRecordPayload.FirstName != "" {
		s = s + " "
	}
	s = s + ndefRecordPayload.LastName

	return s
}

func (ndefRecordPayloadResource NdefRecordPayloadVcardResource) ToPayload() (NdefRecordPayload, error) {
	return NdefRecordPayloadVcard(ndefRecordPayloadResource), nil
}

type NdefRecordPayloadMime struct {
	Type         string
	Format       MimeFormat
	ContentASCII string
	ContentHEX   []byte
}
type NdefRecordPayloadMimeResource struct {
	Type    string `json:"type"`
	Format  string `json:"format"`
	Content string `json:"content"`
}

func (ndefRecordPayload NdefRecordPayloadMime) ToResource() NdefRecordPayloadResource {
	var content string
	if ndefRecordPayload.Format == MimeFormatASCII {
		content = ndefRecordPayload.ContentASCII
	} else if ndefRecordPayload.Format == MimeFormatHex {
		content = base64.StdEncoding.EncodeToString(ndefRecordPayload.ContentHEX)
	}
	resource := NdefRecordPayloadMimeResource{
		Type:    ndefRecordPayload.Type,
		Format:  ndefRecordPayload.Format.String(),
		Content: content,
	}
	return resource
}

func (ndefRecordPayload NdefRecordPayloadMime) ToRecord() *ndef.Record {
	var contentByte []byte

	if ndefRecordPayload.Format == MimeFormatASCII {
		//replace \n with {0x0D, 0x0A} according to spec
		//todo: 1. add escape check for \\n 2. collapse spaces
		crlf := string([]byte{0x0D, 0x0A})
		contentFormated := strings.Replace(ndefRecordPayload.ContentASCII, "\n", crlf, -1)
		contentByte = []byte(contentFormated)
	} else if ndefRecordPayload.Format == MimeFormatHex {
		contentByte = ndefRecordPayload.ContentHEX
	}
	record := ndef.NewMediaRecord(ndefRecordPayload.Type, contentByte)
	return record
}

func (ndefRecordPayload NdefRecordPayloadMime) String() string {
	var s string
	if ndefRecordPayload.Format == MimeFormatHex {
		s = fmt.Sprintf("% x", ndefRecordPayload.ContentHEX)
	}
	if ndefRecordPayload.Format == MimeFormatASCII {
		s = ndefRecordPayload.ContentASCII
	}
	return s
}
func (ndefRecordPayloadResource NdefRecordPayloadMimeResource) ToPayload() (NdefRecordPayload, error) {
	var contentASCII string
	var contentHEX []byte
	mimeFormat, _ := StringToMimeFormat(ndefRecordPayloadResource.Format)

	if mimeFormat == MimeFormatASCII {
		contentASCII = ndefRecordPayloadResource.Content
	}
	if mimeFormat == MimeFormatHex {
		decodedBytes, err := base64.StdEncoding.DecodeString(ndefRecordPayloadResource.Content)
		if err != nil {
			fmt.Println("decode error:", err)
			return NdefRecordPayloadMime{}, errors.New("Could not decode content. It should be base64 encoded for hex type mime")
		}

		contentHEX = decodedBytes
	}

	ndefRecordPayload := NdefRecordPayloadMime{
		Type:         ndefRecordPayloadResource.Type,
		Format:       mimeFormat,
		ContentASCII: contentASCII,
		ContentHEX:   contentHEX,
	}
	return ndefRecordPayload, nil
}

type MimeFormat int

const (
	MimeFormatASCII MimeFormat = iota + 1
	MimeFormatHex
)

func StringToMimeFormat(s string) (MimeFormat, bool) {
	switch s {
	case MimeFormatASCII.String():
		return MimeFormatASCII, true
	case MimeFormatHex.String():
		return MimeFormatHex, true
	}
	return 0, false
}

func (mimeFormat MimeFormat) String() string {
	names := [...]string{
		"unknown",
		"ascii",
		"hex",
	}

	if mimeFormat < MimeFormatASCII || mimeFormat > MimeFormatHex {
		return names[0]
	}
	return names[mimeFormat]
}

type NdefRecordPayloadPhone struct {
	PhoneNumber string
}
type NdefRecordPayloadPhoneResource struct {
	PhoneNumber string `json:"phone_number"`
}

func (ndefRecordPayload NdefRecordPayloadPhone) ToResource() NdefRecordPayloadResource {
	return NdefRecordPayloadPhoneResource(ndefRecordPayload)
}

func (ndefRecordPayload NdefRecordPayloadPhone) ToRecord() *ndef.Record {
	record := ndef.NewURIRecord("tel:" + ndefRecordPayload.PhoneNumber)
	return record
}
func (ndefRecordPayload NdefRecordPayloadPhone) String() string {
	return ndefRecordPayload.PhoneNumber
}

func (ndefRecordPayloadResource NdefRecordPayloadPhoneResource) ToPayload() (NdefRecordPayload, error) {
	return NdefRecordPayloadPhone(ndefRecordPayloadResource), nil
}

type NdefRecordPayloadGeo struct {
	Latitude  string
	Longitude string
}
type NdefRecordPayloadGeoResource struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

func (ndefRecordPayload NdefRecordPayloadGeo) ToResource() NdefRecordPayloadResource {
	return NdefRecordPayloadGeoResource(ndefRecordPayload)
}

func (ndefRecordPayload NdefRecordPayloadGeo) ToRecord() *ndef.Record {
	record := ndef.NewURIRecord("geo:" + ndefRecordPayload.Latitude + "," + ndefRecordPayload.Longitude)
	return record
}

func (ndefRecordPayload NdefRecordPayloadGeo) String() string {
	return fmt.Sprintf("%s, %s", ndefRecordPayload.Latitude, ndefRecordPayload.Longitude)
}

func (ndefRecordPayloadResource NdefRecordPayloadGeoResource) ToPayload() (NdefRecordPayload, error) {
	return NdefRecordPayloadGeo(ndefRecordPayloadResource), nil
}

type NdefRecordPayloadAar struct {
	PackageName string
}
type NdefRecordPayloadAarResource struct {
	PackageName string `json:"package_name"`
}

func (ndefRecordPayload NdefRecordPayloadAar) ToResource() NdefRecordPayloadResource {
	return NdefRecordPayloadAarResource(ndefRecordPayload)
}
func (ndefRecordPayload NdefRecordPayloadAar) ToRecord() *ndef.Record {
	record := ndef.NewExternalRecord("android.com:pkg", []byte(ndefRecordPayload.PackageName))
	return record
}
func (ndefRecordPayload NdefRecordPayloadAar) String() string {
	return ndefRecordPayload.PackageName
}

func (ndefRecordPayloadResource NdefRecordPayloadAarResource) ToPayload() (NdefRecordPayload, error) {
	return NdefRecordPayloadAar(ndefRecordPayloadResource), nil
}

type NdefRecordPayloadPoster struct {
	Title string
	Uri   string
}
type NdefRecordPayloadPosterResource struct {
	Title string `json:"title"`
	Uri   string `json:"uri"`
}

func (ndefRecordPayload NdefRecordPayloadPoster) ToResource() NdefRecordPayloadResource {
	return NdefRecordPayloadPosterResource(ndefRecordPayload)
}

func (ndefRecordPayload NdefRecordPayloadPoster) ToRecord() *ndef.Record {

	recordURI := ndef.NewURIRecord(ndefRecordPayload.Uri)
	recordTitle := ndef.NewTextRecord(ndefRecordPayload.Title, "en")
	posterPayload := ndef.NewMessageFromRecords(recordTitle, recordURI)
	poster := ndef.NewSmartPosterMessage(posterPayload)
	record := poster.Records[0]
	return record
}

func (ndefRecordPayload NdefRecordPayloadPoster) String() string {
	return fmt.Sprintf("%s, %s", ndefRecordPayload.Title, ndefRecordPayload.Uri)
}

func (ndefRecordPayloadResource NdefRecordPayloadPosterResource) ToPayload() (NdefRecordPayload, error) {
	return NdefRecordPayloadPoster(ndefRecordPayloadResource), nil
}

func RecordToNdefRecord(record *ndef.Record) NdefRecord {
	var ndefRecord NdefRecord
	useRaw := false
	recordTNF := record.TNF()
	recordType := record.Type()
	recordPayload, err := record.Payload()
	if err != nil {
		recordPayload = generic.New([]byte{})
	}
	recordPayloadStr := recordPayload.String()
	recordPayloadBytes := recordPayload.Marshal()
	recordID := record.ID()

	switch recordTNF {
	case ndef.NFCForumWellKnownType:
		switch recordType {
		case "U":
			if strings.HasPrefix(recordPayloadStr, "http://") || strings.HasPrefix(recordPayloadStr, "https://") {
				payloadUrl := NdefRecordPayloadUrl{
					Url: recordPayloadStr,
				}
				ndefRecord = NdefRecord{
					Type: NdefRecordPayloadTypeUrl,
					Data: payloadUrl,
				}
			} else if strings.HasPrefix(recordPayloadStr, "tel:") {
				payloadPhone := NdefRecordPayloadPhone{
					PhoneNumber: recordPayloadStr[4:],
				}
				ndefRecord = NdefRecord{
					Type: NdefRecordPayloadTypePhone,
					Data: payloadPhone,
				}
			} else if strings.HasPrefix(recordPayloadStr, "geo:") {
				latitude := ""
				longitude := ""
				recordPayloadStr = recordPayloadStr[4:]
				i := strings.Index(recordPayloadStr, ",")

				if i > -1 {
					latitude = recordPayloadStr[:i]
					longitude = recordPayloadStr[i+1:]
				} else {
					latitude = recordPayloadStr
				}
				payloadGeo := NdefRecordPayloadGeo{
					Latitude:  latitude,
					Longitude: longitude,
				}
				ndefRecord = NdefRecord{
					Type: NdefRecordPayloadTypeGeo,
					Data: payloadGeo,
				}
			} else {
				payloadUri := NdefRecordPayloadUri{
					Uri: recordPayloadStr,
				}
				ndefRecord = NdefRecord{
					Type: NdefRecordPayloadTypeUri,
					Data: payloadUri,
				}
			}

		case "T":
			langStr := ""
			textStr := ""
			textPayload, ok := recordPayload.(*text.Payload)
			if ok {
				langStr = CodeToLang(textPayload.Language)
				//replace new line {0x0D, 0x0A} with \n
				crlf := string([]byte{0x0D, 0x0A})
				textStr = strings.Replace(textPayload.Text, crlf, "\n", -1)
			}
			payloadText := NdefRecordPayloadText{
				Text: textStr,
				Lang: langStr,
			}
			ndefRecord = NdefRecord{
				Type: NdefRecordPayloadTypeText,
				Data: payloadText,
			}
		case "Sp":
			title := ""
			uri := ""
			smartPosterPayload, ok := recordPayload.(*ndef.SmartPosterPayload)
			if ok {
				posterMessage := smartPosterPayload.Message
				for _, spRecord := range posterMessage.Records {
					if spRecord.TNF() == ndef.NFCForumWellKnownType {
						spRecordPayload, err := spRecord.Payload()
						if err == nil {
							if spRecord.Type() == "U" {
								uri = spRecordPayload.String()
							} else if spRecord.Type() == "T" {
								textPayload, ok := spRecordPayload.(*text.Payload)
								if ok {
									title = textPayload.Text
								}

							}

						}

					}

				}

			}
			payloadPoster := NdefRecordPayloadPoster{
				Title: title,
				Uri:   uri,
			}
			ndefRecord = NdefRecord{
				Type: NdefRecordPayloadTypePoster,
				Data: payloadPoster,
			}
		default:
			useRaw = true

		}

	case ndef.MediaType:
		if recordType == "text/vcard" {
			addressCity := ""
			addressCountry := ""
			addressPostalCode := ""
			addressRegion := ""
			addressStreet := ""
			email := ""
			firstName := ""
			lastName := ""
			organization := ""
			phoneCell := ""
			phoneHome := ""
			phoneWork := ""
			title := ""
			site := ""

			newReader := strings.NewReader(string(recordPayloadBytes))
			dec := vcard.NewDecoder(newReader)
			card, err := dec.Decode()
			if err == nil {
				name := card.Name()
				if name != nil {
					firstName = name.GivenName
					lastName = name.FamilyName
				}
				address := card.Address()
				if address != nil {
					addressCity = address.Locality
					addressCountry = address.Country
					addressPostalCode = address.PostalCode
					addressRegion = address.Region
					addressStreet = address.StreetAddress

				}
				organization = card.Value(vcard.FieldOrganization)
				title = card.Value(vcard.FieldTitle)
				phoneFields, ok := card[vcard.FieldTelephone]
				if ok {
					for _, phoneField := range phoneFields {
						params := phoneField.Params
						if params.HasType(vcard.TypeHome) && params.HasType(vcard.TypeVoice) {
							phoneHome = phoneField.Value
						} else if params.HasType(vcard.TypeWork) && params.HasType(vcard.TypeVoice) {
							phoneWork = phoneField.Value
						} else if params.HasType(vcard.TypeCell) {
							phoneCell = phoneField.Value
						}

					}
				}
				email = card.Value(vcard.FieldEmail)
				site = card.Value(vcard.FieldURL)
			}
			payloadVcard := NdefRecordPayloadVcard{
				AddressCity:       addressCity,
				AddressCountry:    addressCountry,
				AddressPostalCode: addressPostalCode,
				AddressRegion:     addressRegion,
				AddressStreet:     addressStreet,
				Email:             email,
				FirstName:         firstName,
				LastName:          lastName,
				Organization:      organization,
				PhoneCell:         phoneCell,
				PhoneHome:         phoneHome,
				PhoneWork:         phoneWork,
				Title:             title,
				Site:              site,
			}
			ndefRecord = NdefRecord{
				Type: NdefRecordPayloadTypeVcard,
				Data: payloadVcard,
			}

		} else {
			payloadMime := NdefRecordPayloadMime{
				Type:       recordType,
				Format:     MimeFormatHex,
				ContentHEX: recordPayloadBytes,
			}

			ndefRecord = NdefRecord{
				Type: NdefRecordPayloadTypeMime,
				Data: payloadMime,
			}

		}

	case ndef.NFCForumExternalType:
		if recordType == "android.com:pkg" {
			packageName := string(recordPayloadBytes)
			payloadAar := NdefRecordPayloadAar{
				PackageName: packageName,
			}
			ndefRecord = NdefRecord{
				Type: NdefRecordPayloadTypeAar,
				Data: payloadAar,
			}

		} else {
			useRaw = true
		}

	default:
		useRaw = true

	}

	if useRaw {
		payloadRaw := NdefRecordPayloadRaw{
			Tnf:     int(recordTNF),
			Type:    recordType,
			ID:      recordID,
			Payload: recordPayload.Marshal(),
		}

		ndefRecord = NdefRecord{
			Type: NdefRecordPayloadTypeRaw,
			Data: payloadRaw,
		}
	}
	return ndefRecord

}

func TnfToString(tnf int) string {
	var s string
	switch tnf {
	case 0:
		s = "Empty"
	case 1:
		s = "Well-Known"
	case 2:
		s = "MIME media-type"
	case 3:
		s = "Absolute URI"
	case 4:
		s = "External"
	case 5:
		s = "Unknown"
	case 6:
		s = "Unchanged"
	case 7:
		s = "Reserved"

	}
	return s
}
