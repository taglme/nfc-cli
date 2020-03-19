package repository

import (
	"github.com/stretchr/testify/assert"
	apiModels "github.com/taglme/nfc-goclient/pkg/models"
	"github.com/taglme/nfc-goclient/pkg/ndefconv"
	"testing"
)

func Test_parseTagStruct(t *testing.T) {
	tag := map[string]interface{}{
		"tag_id": "id",
		"kind": "kind",
		"href": "link",
		"type": apiModels.TagTypeBluetooth.String(),
		"adapter_id": "id",
		"adapter_name": "name",
		"uid": "qhIyag==",
		"atr": "qhIyag==",
		"product": "product",
		"vendor": "vendor",
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
	assert.Equal(t, "any text(text)\nAccess: read only", res.String())

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
	assert.Equal(t, "[string] string | string\n", res.String())
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
	assert.Equal(t, "any text(text)\n", res.String())

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
		"type": "type",
		"tnf": 6,
		"id": "id",
		"payload": "qhIyag==",
	}

	res := parseNdefPayloadResourceByNdefPayloadType(ndefconv.NdefRecordPayloadTypeRaw, raw)
	assert.Equal(t, "Unchanged, type, aa 12 32 6a", res.String())
}