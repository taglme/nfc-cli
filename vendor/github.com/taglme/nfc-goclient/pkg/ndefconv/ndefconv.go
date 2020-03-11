package ndefconv

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

// Ndef defines NDEF message with additional ReadOnly flag
type Ndef struct {
	ReadOnly bool
	Message  []NdefRecord
}

// NdefResource defines json structure of Ndef resource
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

//ToResource - conversion Ndef structure to json resource
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

//ToNdefRecord - conversion json resource to Ndef structure
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

//NdefRecord defines NDEF record representation
type NdefRecord struct {
	Type NdefRecordPayloadType
	Data NdefRecordPayload
}

func (ndefRecord *NdefRecord) String() string {
	return ndefRecord.Data.String()
}

//NdefRecordResource json structure of NdefRecord resource
type NdefRecordResource struct {
	Type string                    `json:"type"`
	Data NdefRecordPayloadResource `json:"data"`
}

//UnmarshalJSON - unmarshal NdefRecordResource structure
func (ndefRecordResource *NdefRecordResource) UnmarshalJSON(data []byte) error {

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
	ndefRecordResource.Type = t

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
		ndefRecordResource.Data = r
	case NdefRecordPayloadTypeUrl:
		r := NdefRecordPayloadUrlResource{}
		err := json.Unmarshal(dataBytes, &r)
		if err != nil {
			return err
		}
		if r.Url == "" {
			return errors.New("Url field of Url type record should be not empty")
		}
		ndefRecordResource.Data = r
	case NdefRecordPayloadTypeText:
		r := NdefRecordPayloadTextResource{}
		err := json.Unmarshal(dataBytes, &r)
		if err != nil {
			return err
		}
		if r.Text == "" {
			return errors.New("Text field of Text type record should be not empty")
		}
		ndefRecordResource.Data = r
	case NdefRecordPayloadTypeUri:
		r := NdefRecordPayloadUriResource{}
		err := json.Unmarshal(dataBytes, &r)
		if err != nil {
			return err
		}
		if r.Uri == "" {
			return errors.New("Uri field of Uri type record should be not empty")
		}
		ndefRecordResource.Data = r
	case NdefRecordPayloadTypeVcard:
		r := NdefRecordPayloadVcardResource{}
		err := json.Unmarshal(dataBytes, &r)
		if err != nil {
			return err
		}
		if r.FirstName == "" {
			return errors.New("First name field of Vcard type record should be not empty")
		}
		ndefRecordResource.Data = r
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

		ndefRecordResource.Data = r
	case NdefRecordPayloadTypePhone:
		r := NdefRecordPayloadPhoneResource{}
		err := json.Unmarshal(dataBytes, &r)
		if err != nil {
			return err
		}
		if r.PhoneNumber == "" {
			return errors.New("Phone number field of Phone type record should be not empty")
		}

		ndefRecordResource.Data = r
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

		ndefRecordResource.Data = r
	case NdefRecordPayloadTypeAar:
		r := NdefRecordPayloadAarResource{}
		err := json.Unmarshal(dataBytes, &r)
		if err != nil {
			return err
		}
		if r.PackageName == "" {
			return errors.New("Package name field of Android application type record should be not empty")
		}
		ndefRecordResource.Data = r
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
		ndefRecordResource.Data = r
	}
	return nil
}

//ToResource - NdefRecord conversion to corresponding json resource
func (ndefRecord NdefRecord) ToResource() NdefRecordResource {
	resource := NdefRecordResource{
		Type: ndefRecord.Type.String(),
		Data: ndefRecord.Data.ToResource(),
	}

	return resource
}

//ToNdefRecord - NdefRecordResource conversion to NdefRecord structure
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

//NdefRecordPayload represents NDEF record payload interface
type NdefRecordPayload interface {
	ToResource() NdefRecordPayloadResource
	String() string
}

//NdefRecordPayloadResource defines NDEF record payload resource interface
type NdefRecordPayloadResource interface {
	ToPayload() (NdefRecordPayload, error)
}

//NdefRecordPayloadType defines types of NDEF record payloads
type NdefRecordPayloadType int

const (
	//NdefRecordPayloadTypeRaw - raw type of NDEF record payload
	NdefRecordPayloadTypeRaw NdefRecordPayloadType = iota + 1
	//NdefRecordPayloadTypeUrl - URL type of NDEF record payload
	NdefRecordPayloadTypeUrl
	//NdefRecordPayloadTypeText - text type of NDEF record payload
	NdefRecordPayloadTypeText
	//NdefRecordPayloadTypeUri - URI type of NDEF record payload
	NdefRecordPayloadTypeUri
	//NdefRecordPayloadTypeVcard - VCARD type of NDEF record payload
	NdefRecordPayloadTypeVcard
	//NdefRecordPayloadTypeMime - mime type of NDEF record payload
	NdefRecordPayloadTypeMime
	//NdefRecordPayloadTypePhone - phone type of NDEF record payload
	NdefRecordPayloadTypePhone
	//NdefRecordPayloadTypeGeo - geo type of NDEF record payload
	NdefRecordPayloadTypeGeo
	//NdefRecordPayloadTypeAar - android package type of NDEF record payload
	NdefRecordPayloadTypeAar
	//NdefRecordPayloadTypePoster - smart poster type of NDEF record payload
	NdefRecordPayloadTypePoster
)

//StringToNdefRecordPayloadType - string conversion to specific NdefRecordPayloadType type
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

//NdefRecordPayloadRaw defines raw NDEF record payload structure
type NdefRecordPayloadRaw struct {
	Tnf     int
	Type    string
	ID      string
	Payload []byte
}

//NdefRecordPayloadRawResource defines json resource structure of raw NDEF record payload structure
type NdefRecordPayloadRawResource struct {
	Tnf     int    `json:"tnf"`
	Type    string `json:"type"`
	ID      string `json:"id"`
	Payload string `json:"payload"`
}

//ToResource conversion raw NDEF record payload structure to json resource
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

//String pretty print of raw NDEF record payload
func (ndefRecordPayload NdefRecordPayloadRaw) String() string {
	return fmt.Sprintf("%s, %s, % x", TnfToString(ndefRecordPayload.Tnf), ndefRecordPayload.Type, ndefRecordPayload.Payload)
}

//ToPayload conversion raw NDEF record payload resource to corresponding structure
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

//NdefRecordPayloadUrl defines url NDEF record payload structure
type NdefRecordPayloadUrl struct {
	Url string
}

//NdefRecordPayloadUrlResource defines json resource structure of url NDEF record payload structure
type NdefRecordPayloadUrlResource struct {
	Url string `json:"url"`
}

//ToResource conversion url NDEF record payload structure to json resource
func (ndefRecordPayload NdefRecordPayloadUrl) ToResource() NdefRecordPayloadResource {
	return NdefRecordPayloadUrlResource(ndefRecordPayload)
}

//String pretty print of url NDEF record payload
func (ndefRecordPayload NdefRecordPayloadUrl) String() string {
	return ndefRecordPayload.Url
}

//ToPayload conversion url NDEF record payload resource to corresponding structure
func (ndefRecordPayloadResource NdefRecordPayloadUrlResource) ToPayload() (NdefRecordPayload, error) {
	return NdefRecordPayloadUrl(ndefRecordPayloadResource), nil
}

//NdefRecordPayloadText defines text NDEF record payload structure
type NdefRecordPayloadText struct {
	Text string
	Lang string
}

//NdefRecordPayloadTextResource defines json resource structure of text NDEF record payload structure
type NdefRecordPayloadTextResource struct {
	Text string `json:"text"`
	Lang string `json:"lang"`
}

//ToResource conversion text NDEF record payload structure to json resource
func (ndefRecordPayload NdefRecordPayloadText) ToResource() NdefRecordPayloadResource {
	return NdefRecordPayloadTextResource(ndefRecordPayload)
}

//String pretty print of text NDEF record payload
func (ndefRecordPayload NdefRecordPayloadText) String() string {

	return ndefRecordPayload.Text
}

//ToPayload conversion text NDEF record payload resource to corresponding structure
func (ndefRecordPayloadResource NdefRecordPayloadTextResource) ToPayload() (NdefRecordPayload, error) {
	return NdefRecordPayloadText(ndefRecordPayloadResource), nil
}

//LangToCode - language conversion to code
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

//CodeToLang - code conversion to language
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

//NdefRecordPayloadUri defines uri NDEF record payload structure
type NdefRecordPayloadUri struct {
	Uri string
}

//NdefRecordPayloadUriResource defines json resource structure of uri NDEF record payload structure
type NdefRecordPayloadUriResource struct {
	Uri string `json:"uri"`
}

//ToResource conversion uri NDEF record payload structure to json resource
func (ndefRecordPayload NdefRecordPayloadUri) ToResource() NdefRecordPayloadResource {
	return NdefRecordPayloadUriResource(ndefRecordPayload)
}

//String pretty print of uri NDEF record payload
func (ndefRecordPayload NdefRecordPayloadUri) String() string {
	return ndefRecordPayload.Uri
}

//ToPayload conversion uri NDEF record payload resource to corresponding structure
func (ndefRecordPayloadResource NdefRecordPayloadUriResource) ToPayload() (NdefRecordPayload, error) {
	return NdefRecordPayloadUri(ndefRecordPayloadResource), nil
}

//NdefRecordPayloadVcard defines vcard NDEF record payload structure
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

//NdefRecordPayloadVcardResource defines json resource structure of vcard NDEF record payload structure
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

//ToResource conversion vcard NDEF record payload structure to json resource
func (ndefRecordPayload NdefRecordPayloadVcard) ToResource() NdefRecordPayloadResource {
	return NdefRecordPayloadVcardResource(ndefRecordPayload)
}

//String pretty print of vcard NDEF record payload
func (ndefRecordPayload NdefRecordPayloadVcard) String() string {
	s := ndefRecordPayload.FirstName
	if ndefRecordPayload.LastName != "" && ndefRecordPayload.FirstName != "" {
		s = s + " "
	}
	s = s + ndefRecordPayload.LastName

	return s
}

//ToPayload conversion vcard NDEF record payload resource to corresponding structure
func (ndefRecordPayloadResource NdefRecordPayloadVcardResource) ToPayload() (NdefRecordPayload, error) {
	return NdefRecordPayloadVcard(ndefRecordPayloadResource), nil
}

//NdefRecordPayloadMime defines mime NDEF record payload structure
type NdefRecordPayloadMime struct {
	Type         string
	Format       MimeFormat
	ContentASCII string
	ContentHEX   []byte
}

//NdefRecordPayloadMimeResource defines json resource structure of mime NDEF record payload structure
type NdefRecordPayloadMimeResource struct {
	Type    string `json:"type"`
	Format  string `json:"format"`
	Content string `json:"content"`
}

//ToResource conversion mime NDEF record payload structure to json resource
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

//String pretty print of mime NDEF record payload
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

//ToPayload conversion mime NDEF record payload resource to corresponding structure
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

//MimeFormat - specific type represented mime data format
type MimeFormat int

const (
	//MimeFormatASCII - ascii format of mime data
	MimeFormatASCII MimeFormat = iota + 1
	//MimeFormatHex - hex format of mime data
	MimeFormatHex
)

//StringToMimeFormat string conversion to specific mime format type
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

//NdefRecordPayloadPhone defines phone NDEF record payload structure
type NdefRecordPayloadPhone struct {
	PhoneNumber string
}

//NdefRecordPayloadPhoneResource defines json resource structure of phone NDEF record payload structure
type NdefRecordPayloadPhoneResource struct {
	PhoneNumber string `json:"phone_number"`
}

//ToResource conversion phone NDEF record payload structure to json resource
func (ndefRecordPayload NdefRecordPayloadPhone) ToResource() NdefRecordPayloadResource {
	return NdefRecordPayloadPhoneResource(ndefRecordPayload)
}

//String pretty print of phone NDEF record payload
func (ndefRecordPayload NdefRecordPayloadPhone) String() string {
	return ndefRecordPayload.PhoneNumber
}

//ToPayload conversion phone NDEF record payload resource to corresponding structure
func (ndefRecordPayloadResource NdefRecordPayloadPhoneResource) ToPayload() (NdefRecordPayload, error) {
	return NdefRecordPayloadPhone(ndefRecordPayloadResource), nil
}

//NdefRecordPayloadGeo defines geo NDEF record payload structure
type NdefRecordPayloadGeo struct {
	Latitude  string
	Longitude string
}

//NdefRecordPayloadGeoResource defines json resource structure of geo NDEF record payload structure
type NdefRecordPayloadGeoResource struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

//ToResource conversion geo NDEF record payload structure to json resource
func (ndefRecordPayload NdefRecordPayloadGeo) ToResource() NdefRecordPayloadResource {
	return NdefRecordPayloadGeoResource(ndefRecordPayload)
}

//String pretty print of geo NDEF record payload
func (ndefRecordPayload NdefRecordPayloadGeo) String() string {
	return fmt.Sprintf("%s, %s", ndefRecordPayload.Latitude, ndefRecordPayload.Longitude)
}

//ToPayload conversion geo NDEF record payload resource to corresponding structure
func (ndefRecordPayloadResource NdefRecordPayloadGeoResource) ToPayload() (NdefRecordPayload, error) {
	return NdefRecordPayloadGeo(ndefRecordPayloadResource), nil
}

//NdefRecordPayloadAar defines Android package NDEF record payload structure
type NdefRecordPayloadAar struct {
	PackageName string
}

//NdefRecordPayloadAarResource defines json resource structure of Android package NDEF record payload structure
type NdefRecordPayloadAarResource struct {
	PackageName string `json:"package_name"`
}

//ToResource conversion Android package NDEF record payload structure to json resource
func (ndefRecordPayload NdefRecordPayloadAar) ToResource() NdefRecordPayloadResource {
	return NdefRecordPayloadAarResource(ndefRecordPayload)
}

//String pretty print of Android package NDEF record payload
func (ndefRecordPayload NdefRecordPayloadAar) String() string {
	return ndefRecordPayload.PackageName
}

//ToPayload conversion Android package NDEF record payload resource to corresponding structure
func (ndefRecordPayloadResource NdefRecordPayloadAarResource) ToPayload() (NdefRecordPayload, error) {
	return NdefRecordPayloadAar(ndefRecordPayloadResource), nil
}

//NdefRecordPayloadPoster defines smart poster NDEF record payload structure
type NdefRecordPayloadPoster struct {
	Title string
	Uri   string
}

//NdefRecordPayloadPosterResource defines json resource structure of smart poster NDEF record payload structure
type NdefRecordPayloadPosterResource struct {
	Title string `json:"title"`
	Uri   string `json:"uri"`
}

//ToResource conversion smart poster NDEF record payload structure to json resource
func (ndefRecordPayload NdefRecordPayloadPoster) ToResource() NdefRecordPayloadResource {
	return NdefRecordPayloadPosterResource(ndefRecordPayload)
}

//String pretty print of smart poster NDEF record payload
func (ndefRecordPayload NdefRecordPayloadPoster) String() string {
	return fmt.Sprintf("%s, %s", ndefRecordPayload.Title, ndefRecordPayload.Uri)
}

//ToPayload conversion smart poster NDEF record payload resource to corresponding structure
func (ndefRecordPayloadResource NdefRecordPayloadPosterResource) ToPayload() (NdefRecordPayload, error) {
	return NdefRecordPayloadPoster(ndefRecordPayloadResource), nil
}

//TnfToString - raw NDEF record tnf code conversion to string
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
