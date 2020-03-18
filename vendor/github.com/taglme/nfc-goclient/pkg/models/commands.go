package models

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/taglme/nfc-goclient/pkg/ndefconv"
)

type CommandParams interface {
	ToResource() CommandParamsResource
	String() string
}
type CommandParamsResource interface {
	ToParams() (CommandParams, error)
}

type Command int

const (
	CommandGetTags Command = iota + 1
	CommandTransmitAdapter
	CommandTransmitTag
	CommandWriteNdef
	CommandReadNdef
	CommandFormatDefault
	CommandLockPermanent
	CommandSetPassword
	CommandRemovePassword
	CommandAuthPassword
	CommandGetDump
	CommandSetLocale
)

func StringToCommand(s string) (Command, bool) {
	switch s {
	case CommandGetTags.String():
		return CommandGetTags, true
	case CommandTransmitAdapter.String():
		return CommandTransmitAdapter, true
	case CommandTransmitTag.String():
		return CommandTransmitTag, true
	case CommandWriteNdef.String():
		return CommandWriteNdef, true
	case CommandReadNdef.String():
		return CommandReadNdef, true
	case CommandFormatDefault.String():
		return CommandFormatDefault, true
	case CommandLockPermanent.String():
		return CommandLockPermanent, true
	case CommandSetPassword.String():
		return CommandSetPassword, true
	case CommandRemovePassword.String():
		return CommandRemovePassword, true
	case CommandAuthPassword.String():
		return CommandAuthPassword, true
	case CommandGetDump.String():
		return CommandGetDump, true
	case CommandSetLocale.String():
		return CommandSetLocale, true

	}
	return 0, false
}

func (command Command) String() string {
	names := [...]string{
		"unknown",
		"get_tags",
		"transmit_adapter",
		"transmit_tag",
		"write_ndef",
		"read_ndef",
		"format_default",
		"lock_permanent",
		"set_password",
		"remove_password",
		"auth_password",
		"get_dump",
		"set_locale",
	}

	if command < CommandGetTags || command > CommandSetLocale {
		return names[0]
	}
	return names[command]
}

type GetTagsParams struct{}
type GetTagsParamsResource struct{}

func (params GetTagsParams) ToResource() CommandParamsResource { return GetTagsParamsResource{} }
func (params GetTagsParams) String() string                    { return "" }

func (paramsResource GetTagsParamsResource) ToParams() (CommandParams, error) {
	return GetTagsParams{}, nil
}

type TransmitAdapterParams struct {
	TxBytes []byte
}
type TransmitAdapterParamsResource struct {
	TxBytes string `json:"tx_bytes"`
}

func (params TransmitAdapterParams) ToResource() CommandParamsResource {
	encodedString := base64.StdEncoding.EncodeToString(params.TxBytes)
	resource := TransmitAdapterParamsResource{
		TxBytes: encodedString,
	}
	return resource
}
func (params TransmitAdapterParams) String() string {
	return fmt.Sprintf("% x ", params.TxBytes)
}
func (paramsResource TransmitAdapterParamsResource) ToParams() (CommandParams, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(paramsResource.TxBytes)
	if err != nil {
		return TransmitAdapterParams{}, errors.New("Could not decode tx_bytes. It should be base64 encoded")
	}
	params := TransmitAdapterParams{
		TxBytes: decodedBytes,
	}
	return params, nil
}

type TransmitTagParams struct {
	TxBytes []byte
}
type TransmitTagParamsResource struct {
	TxBytes string `json:"tx_bytes"`
}

func (params TransmitTagParams) ToResource() CommandParamsResource {
	encodedString := base64.StdEncoding.EncodeToString(params.TxBytes)
	resource := TransmitTagParamsResource{
		TxBytes: encodedString,
	}
	return resource
}
func (params TransmitTagParams) String() string {
	return fmt.Sprintf("% x ", params.TxBytes)
}
func (paramsResource TransmitTagParamsResource) ToParams() (CommandParams, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(paramsResource.TxBytes)
	if err != nil {
		return TransmitTagParams{}, errors.New("Could not decode tx_bytes. It should be base64 encoded")
	}
	params := TransmitTagParams{
		TxBytes: decodedBytes,
	}
	return params, nil
}

type WriteNdefParams struct {
	Message []ndefconv.NdefRecord
}
type WriteNdefParamsResource struct {
	Message []ndefconv.NdefRecordResource `json:"message"`
}

func (params WriteNdefParams) ToResource() CommandParamsResource {
	var ndefRecordResources []ndefconv.NdefRecordResource
	for _, ndefRecord := range params.Message {
		ndefRecordResources = append(ndefRecordResources, ndefRecord.ToResource())
	}
	resource := WriteNdefParamsResource{
		Message: ndefRecordResources,
	}
	return resource
}
func (params WriteNdefParams) String() string {
	res := ""

	for _, m := range params.Message {
		res += m.String() + "(" + m.Type.String() + ")\n"
	}

	return res
}
func (paramsResource WriteNdefParamsResource) ToParams() (CommandParams, error) {
	var ndefRecords []ndefconv.NdefRecord
	for _, ndefRecordResource := range paramsResource.Message {
		ndefRecord, err := ndefRecordResource.ToNdefRecord()
		if err != nil {
			return nil, err
		}
		ndefRecords = append(ndefRecords, ndefRecord)
	}
	params := WriteNdefParams{
		Message: ndefRecords,
	}
	return params, nil
}

type ReadNdefParams struct{}
type ReadNdefParamsResource struct{}

func (params ReadNdefParams) ToResource() CommandParamsResource { return ReadNdefParamsResource{} }
func (params ReadNdefParams) String() string                    { return "" }
func (paramsResource ReadNdefParamsResource) ToParams() (CommandParams, error) {
	return ReadNdefParams{}, nil
}

type FormatDefaultParams struct{}
type FormatDefaultParamsResource struct{}

func (params FormatDefaultParams) ToResource() CommandParamsResource {
	return FormatDefaultParamsResource{}
}
func (params FormatDefaultParams) String() string { return "" }
func (paramsResource FormatDefaultParamsResource) ToParams() (CommandParams, error) {
	return FormatDefaultParams{}, nil
}

type LockPermanentParams struct{}
type LockPermanentParamsResource struct{}

func (params LockPermanentParams) ToResource() CommandParamsResource {
	return LockPermanentParamsResource{}
}
func (params LockPermanentParams) String() string { return "" }
func (paramsResource LockPermanentParamsResource) ToParams() (CommandParams, error) {
	return LockPermanentParams{}, nil
}

type SetPasswordParams struct {
	Password []byte
}
type SetPasswordParamsResource struct {
	Password string `json:"password"`
}

func (params SetPasswordParams) ToResource() CommandParamsResource {
	encodedString := base64.StdEncoding.EncodeToString(params.Password)
	resource := SetPasswordParamsResource{
		Password: encodedString,
	}
	return resource
}
func (params SetPasswordParams) String() string {
	return fmt.Sprintf("% x ", params.Password)
}
func (paramsResource SetPasswordParamsResource) ToParams() (CommandParams, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(paramsResource.Password)
	if err != nil {
		return SetPasswordParams{}, errors.New("Could not decode password. It should be base64 encoded")
	}
	params := SetPasswordParams{
		Password: decodedBytes,
	}
	return params, nil
}

type RemovePasswordParams struct{}
type RemovePasswordParamsResource struct{}

func (params RemovePasswordParams) ToResource() CommandParamsResource {
	return RemovePasswordParamsResource{}
}
func (params RemovePasswordParams) String() string { return "" }
func (paramsResource RemovePasswordParamsResource) ToParams() (CommandParams, error) {
	return RemovePasswordParams{}, nil
}

type AuthPasswordParams struct {
	Password []byte
}
type AuthPasswordParamsResource struct {
	Password string `json:"password"`
}

func (params AuthPasswordParams) ToResource() CommandParamsResource {
	encodedString := base64.StdEncoding.EncodeToString(params.Password)
	resource := AuthPasswordParamsResource{
		Password: encodedString,
	}
	return resource
}
func (params AuthPasswordParams) String() string {
	return fmt.Sprintf("% x ", params.Password)
}
func (paramsResource AuthPasswordParamsResource) ToParams() (CommandParams, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(paramsResource.Password)
	if err != nil {
		return AuthPasswordParams{}, errors.New("Could not decode password. It should be base64 encoded")
	}
	params := AuthPasswordParams{
		Password: decodedBytes,
	}
	return params, nil
}

type GetDumpParams struct{}
type GetDumpParamsResource struct{}

func (params GetDumpParams) ToResource() CommandParamsResource { return GetDumpParamsResource{} }
func (params GetDumpParams) String() string                    { return "" }
func (paramsResource GetDumpParamsResource) ToParams() (CommandParams, error) {
	return GetDumpParams{}, nil
}

type SetLocaleParams struct {
	Locale Locale
}
type SetLocaleParamsResource struct {
	Locale string `json:"locale"`
}

func (params SetLocaleParams) ToResource() CommandParamsResource {
	resource := SetLocaleParamsResource{
		Locale: params.Locale.String(),
	}
	return resource
}
func (params SetLocaleParams) String() string {
	return params.Locale.String()
}
func (paramsResource SetLocaleParamsResource) ToParams() (CommandParams, error) {
	locale, isValid := StringToLocale(paramsResource.Locale)
	if !isValid {
		return nil, errors.New("Not valid locale. Should be 'en' or 'ru'")
	}
	params := SetLocaleParams{
		Locale: locale,
	}
	return params, nil
}

//CommandOutput

type CommandOutput interface {
	ToResource() CommandOutputResource
	String() string
}
type CommandOutputResource interface {
	ToOutput() (CommandOutput, error)
}

type GetTagsOutput struct {
	Tags []Tag
}
type GetTagsOutputResource struct {
	Tags []TagResource `json:"tags"`
}

func (output GetTagsOutput) ToResource() CommandOutputResource {
	var tagResources []TagResource
	for _, tag := range output.Tags {
		tagResources = append(tagResources, tag.ToResource())
	}
	resource := GetTagsOutputResource{
		Tags: tagResources,
	}
	return resource
}
func (output GetTagsOutput) String() string {
	res := ""

	for _, t := range output.Tags {
		res += fmt.Sprintf("UID: % x\nATR: % x\nProduct: %s\nVendor: %s", t.Uid, t.Atr, t.Product, t.Vendor)
	}

	return res
}
func (outputResource GetTagsOutputResource) ToOutput() (CommandOutput, error) {
	var tags []Tag
	for _, tagResource := range outputResource.Tags {
		tag, err := tagResource.ToTag()
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	output := GetTagsOutput{
		Tags: tags,
	}
	return output, nil
}

type TransmitAdapterOutput struct {
	RxBytes []byte
}
type TransmitAdapterOutputResource struct {
	RxBytes string `json:"rx_bytes"`
}

func (output TransmitAdapterOutput) ToResource() CommandOutputResource {
	encodedString := base64.StdEncoding.EncodeToString(output.RxBytes)
	resource := TransmitAdapterOutputResource{
		RxBytes: encodedString,
	}
	return resource
}
func (output TransmitAdapterOutput) String() string {
	return fmt.Sprintf("% x ", output.RxBytes)
}
func (outputResource TransmitAdapterOutputResource) ToOutput() (CommandOutput, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(outputResource.RxBytes)
	if err != nil {
		return TransmitAdapterOutput{}, errors.New("Could not decode rx_bytes. It should be base64 encoded")
	}
	output := TransmitAdapterOutput{
		RxBytes: decodedBytes,
	}
	return output, nil
}

type TransmitTagOutput struct {
	RxBytes []byte
}
type TransmitTagOutputResource struct {
	RxBytes string `json:"rx_bytes"`
}

func (output TransmitTagOutput) ToResource() CommandOutputResource {
	encodedString := base64.StdEncoding.EncodeToString(output.RxBytes)
	resource := TransmitTagOutputResource{
		RxBytes: encodedString,
	}
	return resource
}
func (output TransmitTagOutput) String() string {
	return fmt.Sprintf("% x ", output.RxBytes)
}
func (outputResource TransmitTagOutputResource) ToOutput() (CommandOutput, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(outputResource.RxBytes)
	if err != nil {
		return TransmitTagOutput{}, errors.New("Could not decode rx_bytes. It should be base64 encoded")
	}
	output := TransmitTagOutput{
		RxBytes: decodedBytes,
	}
	return output, nil
}

type ReadNdefOutput struct {
	Ndef ndefconv.Ndef
}
type ReadNdefOutputResource struct {
	Ndef ndefconv.NdefResource `json:"ndef"`
}

func (output ReadNdefOutput) ToResource() CommandOutputResource {
	return ReadNdefOutputResource{Ndef: output.Ndef.ToResource()}
}
func (output ReadNdefOutput) String() string {
	res := ""

	for _, m := range output.Ndef.Message {
		res += m.String() + "(" + m.Type.String() + ")\n"
	}

	if output.Ndef.ReadOnly {
		res += "Access: read only"
	} else {
		res += "Access: read and write"
	}

	return res
}
func (outputResource ReadNdefOutputResource) ToOutput() (CommandOutput, error) {
	ndef, err := outputResource.Ndef.ToNdefRecord()
	if err != nil {
		return nil, err
	}
	return ReadNdefOutput{
		Ndef: ndef}, nil
}

type WriteNdefOutput struct{}
type WriteNdefOutputResource struct{}

func (output WriteNdefOutput) ToResource() CommandOutputResource { return WriteNdefOutputResource{} }
func (output WriteNdefOutput) String() string { return "" }
func (outputResource WriteNdefOutputResource) ToOutput() (CommandOutput, error) {
	return WriteNdefOutput{}, nil
}

type LockPermanentOutput struct{}
type LockPermanentOutputResource struct{}

func (output LockPermanentOutput) ToResource() CommandOutputResource {
	return LockPermanentOutputResource{}
}
func (output LockPermanentOutput) String() string { return "" }
func (outputResource LockPermanentOutputResource) ToOutput() (CommandOutput, error) {
	return LockPermanentOutput{}, nil
}

type SetPasswordOutput struct{}
type SetPasswordOutputResource struct{}

func (output SetPasswordOutput) ToResource() CommandOutputResource {
	return SetPasswordOutputResource{}
}
func (output SetPasswordOutput) String() string { return "" }
func (outputResource SetPasswordOutputResource) ToOutput() (CommandOutput, error) {
	return SetPasswordOutput{}, nil
}

type RemovePasswordOutput struct{}
type RemovePasswordOutputResource struct{}

func (output RemovePasswordOutput) ToResource() CommandOutputResource {
	return RemovePasswordOutputResource{}
}
func (output RemovePasswordOutput) String() string { return "" }
func (outputResource RemovePasswordOutputResource) ToOutput() (CommandOutput, error) {
	return RemovePasswordOutput{}, nil
}

type AuthPasswordOutput struct {
	Ack []byte
}
type AuthPasswordOutputResource struct {
	Ack string `json:"ack"`
}

func (output AuthPasswordOutput) ToResource() CommandOutputResource {
	encodedString := base64.StdEncoding.EncodeToString(output.Ack)
	resource := AuthPasswordOutputResource{
		Ack: encodedString,
	}
	return resource
}
func (output AuthPasswordOutput) String() string {
	return fmt.Sprintf("% x ", output.Ack)
}
func (outputResource AuthPasswordOutputResource) ToOutput() (CommandOutput, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(outputResource.Ack)
	if err != nil {
		return AuthPasswordOutput{}, errors.New("Could not decode ack. It should be base64 encoded")
	}
	output := AuthPasswordOutput{
		Ack: decodedBytes,
	}
	return output, nil
}

type FormatDefaultOutput struct{}
type FormatDefaultOutputResource struct{}

func (output FormatDefaultOutput) ToResource() CommandOutputResource {
	return FormatDefaultOutputResource{}
}
func (output FormatDefaultOutput) String() string { return "" }
func (outputResource FormatDefaultOutputResource) ToOutput() (CommandOutput, error) {
	return FormatDefaultOutput{}, nil
}

type GetDumpOutput struct {
	MemoryDump []PageDump
}
type GetDumpOutputResource struct {
	MemoryDump []PageDumpResource `json:"memory_dump"`
}

func (output GetDumpOutput) ToResource() CommandOutputResource {
	var pageDumpResources []PageDumpResource
	for _, memoryDump := range output.MemoryDump {
		pageDumpResources = append(pageDumpResources, memoryDump.ToResource())
	}
	resource := GetDumpOutputResource{
		MemoryDump: pageDumpResources,
	}
	return resource
}
func (output GetDumpOutput) String() string {
	res := ""

	for _, m := range output.MemoryDump {
		res += fmt.Sprintf("[%s] %s | %s", m.Page, m.Data, m.Info)
	}

	return res
}
func (outputResource GetDumpOutputResource) ToOutput() (CommandOutput, error) {
	var pageDumps []PageDump
	for _, pageDumpResource := range outputResource.MemoryDump {
		pageDumps = append(pageDumps, pageDumpResource.ToPageDump())
	}
	output := GetDumpOutput{
		MemoryDump: pageDumps,
	}
	return output, nil
}

type PageDump struct {
	Page string
	Data string
	Info string
}

type PageDumpResource struct {
	Page string `json:"page"`
	Data string `json:"data"`
	Info string `json:"info"`
}

func (pageDump PageDump) ToResource() PageDumpResource {
	return PageDumpResource(pageDump)
}
func (pageDumpResource PageDumpResource) ToPageDump() PageDump {
	return PageDump(pageDumpResource)
}
