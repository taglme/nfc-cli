package models

import (
	"encoding/base64"
	"fmt"

	"github.com/pkg/errors"
)

type Tag struct {
	TagID       string
	Type        TagType
	AdapterID   string
	AdapterName string
	Uid         []byte
	Atr         []byte
	Product     string
	Vendor      string
}

type TagResource struct {
	TagID       string `json:"tag_id"`
	Kind        string `json:"kind"`
	Href        string `json:"href"`
	Type        string `json:"type"`
	AdapterID   string `json:"adapter_id"`
	AdapterName string `json:"adapter_name"`
	Uid         string `json:"uid"`
	Atr         string `json:"atr"`
	Product     string `json:"product"`
	Vendor      string `json:"vendor"`
}

type TagShortResource struct {
	TagID       string `json:"tag_id"`
	Kind        string `json:"kind"`
	Href        string `json:"href"`
	Type        string `json:"type"`
	AdapterID   string `json:"adapter_id"`
	AdapterName string `json:"adapter_name"`
	Uid         string `json:"uid"`
}

type TagListResource []TagShortResource

func (t Tag) ToResource() TagResource {
	encodedUid := base64.StdEncoding.EncodeToString(t.Uid)
	encodedAtr := base64.StdEncoding.EncodeToString(t.Atr)
	resource := TagResource{
		TagID:       t.TagID,
		Kind:        "Tag",
		Href:        fmt.Sprintf(`/adapters/%s/tags/%s`, t.AdapterID, t.TagID),
		Type:        t.Type.String(),
		AdapterID:   t.AdapterID,
		AdapterName: t.AdapterName,
		Uid:         encodedUid,
		Atr:         encodedAtr,
		Product:     t.Product,
		Vendor:      t.Vendor,
	}
	return resource
}

func (t Tag) ToShortResource() TagShortResource {
	encodedUid := base64.StdEncoding.EncodeToString(t.Uid)
	resource := TagShortResource{
		TagID:       t.TagID,
		Kind:        "TagShort",
		Href:        fmt.Sprintf(`/adapters/%s/tags/%s`, t.AdapterID, t.TagID),
		Type:        t.Type.String(),
		AdapterID:   t.AdapterID,
		AdapterName: t.AdapterName,
		Uid:         encodedUid,
	}
	return resource
}

func (t TagShortResource) ToTag() (tag Tag, err error) {
	tType, ok := StringToTagType(t.Type)
	if !ok {
		return tag, errors.New("Can't convert type resource category")
	}

	uuid, err := base64.StdEncoding.DecodeString(t.Uid)
	if err != nil {
		return tag, errors.Wrap(err, "Can't decode tag uuid")
	}

	return Tag{
		TagID:       t.TagID,
		Type:        tType,
		AdapterID:   t.AdapterID,
		AdapterName: t.AdapterName,
		Uid:         uuid,
	}, nil
}

func (t TagResource) ToTag() (Tag, error) {
	tagType, _ := StringToTagType(t.Type)
	decodedUid, err := base64.StdEncoding.DecodeString(t.Uid)
	if err != nil {
		return Tag{}, errors.New("Could not decode Uid. It should be base64 encoded")
	}
	decodedAtr, err := base64.StdEncoding.DecodeString(t.Atr)
	if err != nil {
		return Tag{}, errors.New("Could not decode Atr. It should be base64 encoded")
	}

	resource := Tag{
		TagID:       t.TagID,
		Type:        tagType,
		AdapterID:   t.AdapterID,
		AdapterName: t.AdapterName,
		Uid:         decodedUid,
		Atr:         decodedAtr,
		Product:     t.Product,
		Vendor:      t.Vendor,
	}
	return resource, nil
}

type TagType int

const (
	TagTypeNfc TagType = iota + 1
	TagTypeBarcode
	TagTypeBluetooth
)

func StringToTagType(s string) (TagType, bool) {
	switch s {
	case TagTypeNfc.String():
		return TagTypeNfc, true
	case TagTypeBarcode.String():
		return TagTypeBarcode, true
	case TagTypeBluetooth.String():
		return TagTypeBluetooth, true
	}
	return 0, false
}

func (tagType TagType) String() string {
	names := [...]string{
		"unknown",
		"nfc",
		"barcode",
		"bluetooth"}

	if tagType < TagTypeNfc || tagType > TagTypeBluetooth {
		return names[0]
	}
	return names[tagType]
}
