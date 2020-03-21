package repository

import (
	"github.com/stretchr/testify/assert"
	apiModels "github.com/taglme/nfc-goclient/pkg/models"
	"github.com/taglme/nfc-goclient/pkg/ndefconv"
	"testing"
)

func Test_parseTagStruct(t *testing.T) {
	tag := map[string]interface{}{
		"tag_id":       "id",
		"kind":         "kind",
		"href":         "link",
		"type":         apiModels.TagTypeBluetooth.String(),
		"adapter_id":   "id",
		"adapter_name": "name",
		"uid":          "qhIyag==",
		"atr":          "qhIyag==",
		"product":      "product",
		"vendor":       "vendor",
	}

	res := parseTagStruct(tag)

	assert.Equal(t, "id", res.TagID)
	assert.Equal(t, apiModels.TagTypeBluetooth, res.Type)
}

func Test_parseOutputByCmd(t *testing.T) {
	getTags := map[string]interface{}{
		"tags": []interface{}{
			map[string]interface{}{
				"tag_id":       "id",
				"kind":         "kind",
				"href":         "link",
				"type":         apiModels.TagTypeBluetooth.String(),
				"adapter_id":   "id",
				"adapter_name": "name",
				"uid":          "qhIyag==",
				"atr":          "qhIyag==",
				"product":      "product",
				"vendor":       "vendor",
			},
		},
	}
	res := parseOutputByCmd(apiModels.CommandGetTags, getTags)
	assert.Equal(t, "UID: aa 12 32 6a\nATR: aa 12 32 6a\nProduct: product\nVendor: vendor", res.String())

	transmitAdapter := map[string]interface{}{
		"rx_bytes": "qhIyag==",
	}
	res = parseOutputByCmd(apiModels.CommandTransmitAdapter, transmitAdapter)
	assert.Equal(t, "aa 12 32 6a", res.String())

	res = parseOutputByCmd(apiModels.CommandTransmitTag, transmitAdapter)
	assert.Equal(t, "aa 12 32 6a", res.String())

	empty := map[string]interface{}{}
	res = parseOutputByCmd(apiModels.CommandWriteNdef, empty)
	assert.Equal(t, "", res.String())

	readNdef := map[string]interface{}{
		"ndef": map[string]interface{}{
			"read_only": true,
			"message": []interface{}{
				map[string]interface{}{
					"type": "text",
					"data": map[string]interface{}{
						"text": "any text",
						"lang": "English",
					},
				},
			},
		},
	}
	res = parseOutputByCmd(apiModels.CommandReadNdef, readNdef)
	assert.Equal(t, "Record 1: any text (text)\nAccess: read only", res.String())

	res = parseOutputByCmd(apiModels.CommandFormatDefault, empty)
	assert.Equal(t, "", res.String())

	res = parseOutputByCmd(apiModels.CommandLockPermanent, empty)
	assert.Equal(t, "", res.String())

	res = parseOutputByCmd(apiModels.CommandSetPassword, empty)
	assert.Equal(t, "", res.String())

	res = parseOutputByCmd(apiModels.CommandRemovePassword, empty)
	assert.Equal(t, "", res.String())

	auth := map[string]interface{}{
		"ack": "qhIyag==",
	}
	res = parseOutputByCmd(apiModels.CommandAuthPassword, auth)
	assert.Equal(t, "aa 12 32 6a", res.String())

	dump := map[string]interface{}{
		"memory_dump": []interface{}{
			map[string]interface{}{
				"page": "[string]",
				"data": "string",
				"info": "string",
			},
		},
	}
	res = parseOutputByCmd(apiModels.CommandGetDump, dump)
	assert.Equal(t, "[string] string | string", res.String())
}

func Test_parseParamsByCmd(t *testing.T) {
	empty := map[string]interface{}{}

	res := parseParamsByCmd(apiModels.CommandGetTags, empty)
	assert.Equal(t, "", res.String())

	transmitAdapter := map[string]interface{}{
		"tx_bytes": "qhIyag==",
	}
	res = parseParamsByCmd(apiModels.CommandTransmitAdapter, transmitAdapter)
	assert.Equal(t, "aa 12 32 6a", res.String())

	res = parseParamsByCmd(apiModels.CommandTransmitTag, transmitAdapter)
	assert.Equal(t, "aa 12 32 6a", res.String())

	res = parseParamsByCmd(apiModels.CommandReadNdef, empty)
	assert.Equal(t, "", res.String())

	writeNdef := map[string]interface{}{
		"message": []interface{}{
			map[string]interface{}{
				"type": "text",
				"data": map[string]interface{}{
					"text": "any text",
					"lang": "English",
				},
			},
		},
	}
	res = parseParamsByCmd(apiModels.CommandWriteNdef, writeNdef)
	assert.Equal(t, "Record 1: any text (text)", res.String())

	res = parseParamsByCmd(apiModels.CommandFormatDefault, empty)
	assert.Equal(t, "", res.String())

	res = parseParamsByCmd(apiModels.CommandLockPermanent, empty)
	assert.Equal(t, "", res.String())

	auth := map[string]interface{}{
		"password": "qhIyag==",
	}
	res = parseParamsByCmd(apiModels.CommandSetPassword, auth)
	assert.Equal(t, "aa 12 32 6a", res.String())

	auth2 := map[string]interface{}{
		"password": "qhIyag==",
	}
	res = parseParamsByCmd(apiModels.CommandAuthPassword, auth2)
	assert.Equal(t, "aa 12 32 6a", res.String())

	res = parseParamsByCmd(apiModels.CommandRemovePassword, empty)
	assert.Equal(t, "", res.String())

	res = parseParamsByCmd(apiModels.CommandGetDump, empty)
	assert.Equal(t, "", res.String())

	locale := map[string]interface{}{
		"locale": "en",
	}
	res = parseParamsByCmd(apiModels.CommandSetLocale, locale)
	assert.Equal(t, "en", res.String())
}

func Test_parseNdefPayloadResourceByNdefPayloadType(t *testing.T) {
	raw := map[string]interface{}{
		"type":    "type",
		"tnf":     6,
		"id":      "id",
		"payload": "qhIyag==",
	}
	res := parseNdefPayloadResourceByNdefPayloadType(ndefconv.NdefRecordPayloadTypeRaw, raw)
	assert.Equal(t, "Unchanged, type, aa 12 32 6a", res.String())

	url := map[string]interface{}{
		"url": "https://url.com",
	}
	res = parseNdefPayloadResourceByNdefPayloadType(ndefconv.NdefRecordPayloadTypeUrl, url)
	assert.Equal(t, "https://url.com", res.String())

	uri := map[string]interface{}{
		"uri": "https://url.com",
	}
	res = parseNdefPayloadResourceByNdefPayloadType(ndefconv.NdefRecordPayloadTypeUri, uri)
	assert.Equal(t, "https://url.com", res.String())

	text := map[string]interface{}{
		"text": "any text",
		"lang": "English",
	}
	res = parseNdefPayloadResourceByNdefPayloadType(ndefconv.NdefRecordPayloadTypeText, text)
	assert.Equal(t, "any text", res.String())

	vcard := map[string]interface{}{
		"address_city":        "string",
		"address_country":     "string",
		"address_postal_code": "123",
		"address_region":      "string",
		"address_street":      "string",
		"email":               "string@string.cocm",
		"first_name":          "first_name",
		"last_name":           "last_name",
		"organization":        "string",
		"phone_cell":          "string",
		"phone_home":          "string",
		"phone_work":          "string",
		"title":               "string",
		"site":                "string",
	}
	res = parseNdefPayloadResourceByNdefPayloadType(ndefconv.NdefRecordPayloadTypeVcard, vcard)
	assert.Equal(t, "first_name last_name", res.String())

	mime := map[string]interface{}{
		"type":    "hex",
		"format":  "hex",
		"content": "qhIyag==",
	}
	res = parseNdefPayloadResourceByNdefPayloadType(ndefconv.NdefRecordPayloadTypeMime, mime)
	assert.Equal(t, "aa 12 32 6a", res.String())

	mime = map[string]interface{}{
		"type":    "hex",
		"format":  "ascii",
		"content": "ascii text",
	}
	res = parseNdefPayloadResourceByNdefPayloadType(ndefconv.NdefRecordPayloadTypeMime, mime)
	assert.Equal(t, "ascii text", res.String())

	phone := map[string]interface{}{
		"phone_number": "8 800 555 3535",
	}
	res = parseNdefPayloadResourceByNdefPayloadType(ndefconv.NdefRecordPayloadTypePhone, phone)
	assert.Equal(t, "8 800 555 3535", res.String())

	geo := map[string]interface{}{
		"latitude":  "25.3",
		"longitude": "35.3",
	}
	res = parseNdefPayloadResourceByNdefPayloadType(ndefconv.NdefRecordPayloadTypeGeo, geo)
	assert.Equal(t, "25.3, 35.3", res.String())

	aar := map[string]interface{}{
		"package_name": "package name",
	}
	res = parseNdefPayloadResourceByNdefPayloadType(ndefconv.NdefRecordPayloadTypeAar, aar)
	assert.Equal(t, "package name", res.String())

	poster := map[string]interface{}{
		"title": "title",
		"uri":   "https://uri.com",
	}
	res = parseNdefPayloadResourceByNdefPayloadType(ndefconv.NdefRecordPayloadTypePoster, poster)
	assert.Equal(t, "title, https://uri.com", res.String())
}

func Test_parseStepResultStruct(t *testing.T) {
	getTags := map[string]interface{}{
		"tags": []interface{}{
			map[string]interface{}{
				"tag_id":       "id",
				"kind":         "kind",
				"href":         "link",
				"type":         apiModels.TagTypeBluetooth.String(),
				"adapter_id":   "id",
				"adapter_name": "name",
				"uid":          "qhIyag==",
				"atr":          "qhIyag==",
				"product":      "product",
				"vendor":       "vendor",
			},
		},
	}

	sr := map[string]interface{}{
		"message": "Msg it is",
		"status":  apiModels.CommandStatusSuccess.String(),
		"command": apiModels.CommandGetTags.String(),
		"output":  getTags,
		"params":  map[string]interface{}{},
	}

	res := parseStepResultStruct(sr)
	assert.Equal(t, "Msg it is", res.Message)
	assert.Equal(t, apiModels.CommandStatusSuccess, res.Status)
	assert.Equal(t, apiModels.CommandGetTags, res.Command)
	assert.Equal(t, "", res.Params.String())
	assert.Equal(t, "UID: aa 12 32 6a\nATR: aa 12 32 6a\nProduct: product\nVendor: vendor", res.Output.String())
}

func Test_parseJobRunStruct(t *testing.T) {
	getTags := map[string]interface{}{
		"tags": []interface{}{
			map[string]interface{}{
				"tag_id":       "id",
				"kind":         "kind",
				"href":         "link",
				"type":         apiModels.TagTypeBluetooth.String(),
				"adapter_id":   "id",
				"adapter_name": "name",
				"uid":          "qhIyag==",
				"atr":          "qhIyag==",
				"product":      "product",
				"vendor":       "vendor",
			},
		},
	}

	jr := map[string]interface{}{
		"run_id":       "Run ID",
		"job_id":       "Job ID",
		"job_name":     "Job Name",
		"status":       apiModels.JobRunStatusStarted.String(),
		"adapter_id":   "Adapter ID",
		"adapter_name": "Adapter Name",
		"created_at":   "2020-03-19T16:10:33.580Z",
		"tag": map[string]interface{}{
			"tag_id":       "id",
			"kind":         "kind",
			"href":         "link",
			"type":         apiModels.TagTypeBluetooth.String(),
			"adapter_id":   "id",
			"adapter_name": "name",
			"uid":          "qhIyag==",
			"atr":          "qhIyag==",
			"product":      "product",
			"vendor":       "vendor",
		},
		"results": []interface{}{
			map[string]interface{}{
				"message": "Msg it is",
				"status":  apiModels.CommandStatusSuccess.String(),
				"command": apiModels.CommandGetTags.String(),
				"output":  getTags,
				"params":  map[string]interface{}{},
			},
		},
	}

	res := parseJobRunStruct(jr)
	assert.Equal(t, "Run ID", res.RunID)
	assert.Equal(t, "Job ID", res.JobID)
	assert.Equal(t, "Job Name", res.JobName)
	assert.Equal(t, "id", res.Tag.TagID)
	assert.Equal(t, 1, len(res.Results))
}
