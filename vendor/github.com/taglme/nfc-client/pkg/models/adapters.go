package models

import (
	"fmt"
)

type Adapter struct {
	AdapterID string
	Name      string
	Type      AdapterType
	Driver    string
}

type AdapterResource struct {
	AdapterID string `json:"adapter_id"`
	Kind      string `json:"kind"`
	Href      string `json:"href"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Driver    string `json:"driver"`
}

type AdapterShortResource struct {
	AdapterID string `json:"adapter_id"`
	Kind      string `json:"kind"`
	Href      string `json:"href"`
	Name      string `json:"name"`
	Type      string `json:"type"`
}

type AdapterListResource []AdapterShortResource

func (a Adapter) ToResource() AdapterResource {
	resource := AdapterResource{
		AdapterID: a.AdapterID,
		Kind:      "Adapter",
		Href:      fmt.Sprintf(`/adapters/%s`, a.AdapterID),
		Name:      a.Name,
		Type:      a.Type.String(),
		Driver:    a.Driver,
	}
	return resource
}
func (a *AdapterResource) ToAdapter() Adapter {
	adapterType, _ := StringToAdapterType(a.Type)
	return Adapter{
		AdapterID: a.AdapterID,
		Name:      a.Name,
		Type:      adapterType,
		Driver:    a.Driver,
	}
}

func (a *AdapterShortResource) ToAdapter() Adapter {
	adapterType, _ := StringToAdapterType(a.Type)
	return Adapter{
		AdapterID: a.AdapterID,
		Name:      a.Name,
		Type:      adapterType,
	}
}

func (a Adapter) ToShortResource() AdapterShortResource {
	resource := AdapterShortResource{
		AdapterID: a.AdapterID,
		Kind:      "AdapterShort",
		Href:      fmt.Sprintf(`/adapters/%s`, a.AdapterID),
		Name:      a.Name,
		Type:      a.Type.String(),
	}
	return resource
}

type AdapterType int

const (
	AdapterTypeNfc AdapterType = iota + 1
	AdapterTypeBarcode
	AdapterTypeBluetooth
)

func StringToAdapterType(s string) (AdapterType, bool) {
	switch s {
	case AdapterTypeNfc.String():
		return AdapterTypeNfc, true
	case AdapterTypeBarcode.String():
		return AdapterTypeBarcode, true
	case AdapterTypeBluetooth.String():
		return AdapterTypeBluetooth, true
	}
	return 0, false
}

func (adapterType AdapterType) String() string {
	names := [...]string{
		"",
		"nfc",
		"barcode",
		"bluetooth"}

	if adapterType < AdapterTypeNfc || adapterType > AdapterTypeBluetooth {
		return names[0]
	}
	return names[adapterType]
}
