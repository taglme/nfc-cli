package repository

import (
	"encoding/base64"
	"fmt"
	apiModels "github.com/taglme/nfc-goclient/pkg/models"
	"github.com/taglme/nfc-goclient/pkg/ndefconv"
	"time"
)

func parseJobRunStruct(data interface{}) (jr apiModels.JobRun) {
	m := data.(map[string]interface{})

	if runId, ok := m["run_id"].(string); ok {
		jr.RunID = runId
	}

	if jobId, ok := m["job_id"].(string); ok {
		jr.JobID = jobId
	}

	if jobName, ok := m["job_name"].(string); ok {
		jr.JobName = jobName
	}

	if status, ok := m["status"].(string); ok {
		jr.Status, ok = apiModels.StringToJobRunStatus(status)
		if !ok {
			fmt.Println("Can't parse Job run status")
		}
	}

	if adapterId, ok := m["adapter_id"].(string); ok {
		jr.AdapterID = adapterId
	}

	if adapterName, ok := m["adapter_name"].(string); ok {
		jr.AdapterName = adapterName
	}

	if createdAt, ok := m["created_at"].(string); ok {
		t, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			fmt.Println("Error on converting created at string to time")
		}

		jr.CreatedAt = t
	}

	jr.Tag = parseTagStruct(m["tag"])

	if results, ok := m["results"].([]interface{}); ok {
		jr.Results = make([]apiModels.StepResult, len(results))
		for i, rd := range results {
			jr.Results[i] = parseStepResultStruct(rd)
		}
	}

	return jr
}

func parseTagStruct(data interface{}) (t apiModels.Tag) {
	tag := data.(map[string]interface{})

	if adapterId, ok := tag["adapter_id"].(string); ok {
		t.AdapterID = adapterId
	}

	if adapterName, ok := tag["adapter_name"].(string); ok {
		t.AdapterName = adapterName
	}

	if tagId, ok := tag["tag_id"].(string); ok {
		t.TagID = tagId
	}

	if p, ok := tag["product"].(string); ok {
		t.Product = p
	}

	if v, ok := tag["vendor"].(string); ok {
		t.Vendor = v
	}

	if tType, ok := tag["type"].(string); ok {
		t.Type, ok = apiModels.StringToTagType(tType)
		if !ok {
			fmt.Println("Can't parse job run tag ty[ee")
		}
	}

	if uid, ok := tag["uid"].(string); ok {
		dUid, err := base64.StdEncoding.DecodeString(uid)
		if err != nil {
			fmt.Println("Could not decode Uid. It should be base64 encoded")
		}
		t.Uid = dUid
	}

	if atr, ok := tag["atr"].(string); ok {
		dAtr, err := base64.StdEncoding.DecodeString(atr)
		if err != nil {
			fmt.Println("Could not decode Atr. It should be base64 encoded")
		}
		t.Atr = dAtr
	}

	return t
}

func parseStepResultStruct(data interface{}) (sr apiModels.StepResult) {
	r := data.(map[string]interface{})
	if msg, ok := r["message"].(string); ok {
		sr.Message = msg
	}

	if status, ok := r["status"].(string); ok {
		sr.Status, ok = apiModels.StringToCommandStatus(status)
		if !ok {
			fmt.Println("Can't parse step result status")
		}
	}

	if cmd, ok := r["command"].(string); ok {
		sr.Command, ok = apiModels.StringToCommand(cmd)
		if !ok {
			fmt.Println("Can't parse step result command")
		}
	}

	if o, ok := r["output"].(map[string]interface{}); ok {
		sr.Output = parseOutputByCmd(sr.Command, o)
	}

	if p, ok := r["params"].(map[string]interface{}); ok {
		sr.Params = parseParamsByCmd(sr.Command, p)
	}

	return sr
}

func parseParamsByCmd(command apiModels.Command, data map[string]interface{}) apiModels.CommandParams {
	switch command {
	case apiModels.CommandGetTags:
		return apiModels.GetTagsParams{}
	case apiModels.CommandTransmitAdapter:
		if tx, ok := data["tx_bytes"].(string); ok {
			res, err := apiModels.TransmitAdapterParamsResource{
				TxBytes: tx,
			}.ToParams()
			if err != nil {
				fmt.Println("Can't convert TransmitAdapterParamsResource", err)
			}
			return res
		}
	case apiModels.CommandTransmitTag:
		if tx, ok := data["tx_bytes"].(string); ok {
			res, err := apiModels.TransmitTagParamsResource{
				TxBytes: tx,
			}.ToParams()
			if err != nil {
				fmt.Println("Can't convert TransmitTagParams", err)
			}
			return res
		}
	case apiModels.CommandWriteNdef:
		res := apiModels.WriteNdefParams{}
		if m, ok := data["message"].([]interface{}); ok {
			res.Message = make([]ndefconv.NdefRecord, len(m))
			for i, msg := range m {
				md := msg.(map[string]interface{})

				resMsg := ndefconv.NdefRecord{}
				if t, ok := md["type"].(string); ok {
					recT, ok := ndefconv.StringToNdefRecordPayloadType(t)
					if !ok {
						fmt.Println("Can't parse String To Ndef Record Payload Type")
					}
					resMsg.Type = recT
				}

				if ndefData, ok := md["data"].(map[string]interface{}); ok {
					resMsg.Data = parseNdefPayloadResourceByNdefPayloadType(resMsg.Type, ndefData)
				}

				res.Message[i] = resMsg
			}
		}
		return res
	case apiModels.CommandReadNdef:
		return apiModels.ReadNdefParams{}
	case apiModels.CommandFormatDefault:
		return apiModels.FormatDefaultParams{}
	case apiModels.CommandLockPermanent:
		return apiModels.LockPermanentParams{}
	case apiModels.CommandSetPassword:
		if pwd, ok := data["password"].(string); ok {
			res, err := apiModels.SetPasswordParamsResource{
				Password: pwd,
			}.ToParams()
			if err != nil {
				fmt.Println("Can't convert AuthPasswordParamsResource", err)
			}
			return res
		}
	case apiModels.CommandRemovePassword:
		return apiModels.RemovePasswordParams{}
	case apiModels.CommandAuthPassword:
		if pwd, ok := data["password"].(string); ok {
			res, err := apiModels.AuthPasswordParamsResource{
				Password: pwd,
			}.ToParams()
			if err != nil {
				fmt.Println("Can't convert AuthPasswordParamsResource", err)
			}
			return res
		}
	case apiModels.CommandGetDump:
		return apiModels.GetDumpParams{}
	case apiModels.CommandSetLocale:
		if locale, ok := data["locale"].(string); ok {
			res, err := apiModels.SetLocaleParamsResource{
				Locale: locale,
			}.ToParams()
			if err != nil {
				fmt.Println("Can't convert SetLocaleParamsResource", err)
			}
			return res
		}
	}

	return nil
}

func parseNdefPayloadResourceByNdefPayloadType(t ndefconv.NdefRecordPayloadType, data map[string]interface{}) ndefconv.NdefRecordPayload {
	switch t {
	case ndefconv.NdefRecordPayloadTypeRaw:
		res := ndefconv.NdefRecordPayloadRawResource{}
		if tnf, ok := data["locale"].(int); ok {
			res.Tnf = tnf
		}
		if t, ok := data["type"].(string); ok {
			res.Type = t
		}
		if id, ok := data["id"].(string); ok {
			res.ID = id
		}
		if payload, ok := data["payload"].(string); ok {
			res.Payload = payload
		}
		p, err := res.ToPayload()
		if err != nil {
			fmt.Println("Can't convert NdefRecordPayloadRawResource", err)
		}
		return p
	case ndefconv.NdefRecordPayloadTypeUrl:
		res := ndefconv.NdefRecordPayloadUrl{}
		if u, ok := data["url"].(string); ok {
			res.Url = u
		}
		return res
	case ndefconv.NdefRecordPayloadTypeText:
		res := ndefconv.NdefRecordPayloadTextResource{}
		if t, ok := data["text"].(string); ok {
			res.Text = t
		}
		if t, ok := data["lang"].(string); ok {
			res.Lang = t
		}
		p, err := res.ToPayload()
		if err != nil {
			fmt.Println("Can't convert NdefRecordPayloadTextResource", err)
		}
		return p
	case ndefconv.NdefRecordPayloadTypeUri:
		res := ndefconv.NdefRecordPayloadUri{}
		if u, ok := data["uri"].(string); ok {
			res.Uri = u
		}
		return res
	case ndefconv.NdefRecordPayloadTypeVcard:
		res := ndefconv.NdefRecordPayloadVcardResource{}
		if u, ok := data["address_city"].(string); ok {
			res.AddressCity = u
		}
		if u, ok := data["address_country"].(string); ok {
			res.AddressCountry = u
		}
		if u, ok := data["address_postal_code"].(string); ok {
			res.AddressPostalCode = u
		}
		if u, ok := data["address_region"].(string); ok {
			res.AddressRegion = u
		}
		if u, ok := data["address_street"].(string); ok {
			res.AddressStreet = u
		}
		if u, ok := data["email"].(string); ok {
			res.Email = u
		}
		if u, ok := data["first_name"].(string); ok {
			res.FirstName = u
		}
		if u, ok := data["last_name"].(string); ok {
			res.LastName = u
		}
		if u, ok := data["organization"].(string); ok {
			res.Organization = u
		}
		if u, ok := data["phone_cell"].(string); ok {
			res.PhoneCell = u
		}
		if u, ok := data["phone_home"].(string); ok {
			res.PhoneHome = u
		}
		if u, ok := data["phone_work"].(string); ok {
			res.PhoneWork = u
		}
		if u, ok := data["title"].(string); ok {
			res.Title = u
		}
		if u, ok := data["site"].(string); ok {
			res.Site = u
		}

		p, err := res.ToPayload()
		if err != nil {
			fmt.Println("Can't convert NdefRecordPayloadVcardResource", err)
		}
		return p
	case ndefconv.NdefRecordPayloadTypeMime:
		res := ndefconv.NdefRecordPayloadMimeResource{}
		if t, ok := data["type"].(string); ok {
			res.Type = t
		}
		if t, ok := data["format"].(string); ok {
			res.Format = t
		}
		if t, ok := data["content"].(string); ok {
			res.Content = t
		}
		p, err := res.ToPayload()
		if err != nil {
			fmt.Println("Can't convert NdefRecordPayloadMimeResource", err)
		}
		return p
	case ndefconv.NdefRecordPayloadTypePhone:
		res := ndefconv.NdefRecordPayloadPhoneResource{}
		if t, ok := data["phone_number"].(string); ok {
			res.PhoneNumber = t
		}
		p, err := res.ToPayload()
		if err != nil {
			fmt.Println("Can't convert NdefRecordPayloadPhoneResource", err)
		}
		return p
	case ndefconv.NdefRecordPayloadTypeGeo:
		res := ndefconv.NdefRecordPayloadGeoResource{}
		if t, ok := data["latitude"].(string); ok {
			res.Latitude = t
		}
		if t, ok := data["longitude"].(string); ok {
			res.Longitude = t
		}
		p, err := res.ToPayload()
		if err != nil {
			fmt.Println("Can't convert NdefRecordPayloadGeoResource", err)
		}
		return p
	case ndefconv.NdefRecordPayloadTypeAar:
		res := ndefconv.NdefRecordPayloadAarResource{}
		if t, ok := data["package_name"].(string); ok {
			res.PackageName = t
		}
		p, err := res.ToPayload()
		if err != nil {
			fmt.Println("Can't convert NdefRecordPayloadAarResource", err)
		}
		return p
	case ndefconv.NdefRecordPayloadTypePoster:
		res := ndefconv.NdefRecordPayloadPosterResource{}
		if t, ok := data["title"].(string); ok {
			res.Title = t
		}
		if t, ok := data["uri"].(string); ok {
			res.Uri = t
		}
		p, err := res.ToPayload()
		if err != nil {
			fmt.Println("Can't convert NdefRecordPayloadPosterResource", err)
		}
		return p
	}

	return nil
}

func parseOutputByCmd(command apiModels.Command, data map[string]interface{}) apiModels.CommandOutput {
	switch command {
	case apiModels.CommandGetTags:
		res := apiModels.GetTagsOutput{}
		if tags, ok := data["tags"].([]interface{}); ok {
			sl := make([]apiModels.Tag, len(tags))
			for i, t := range tags {
				sl[i] = parseTagStruct(t)
			}
			res.Tags = sl
		}
		return res
	case apiModels.CommandTransmitAdapter:
		res := apiModels.TransmitAdapterOutputResource{}
		if rx, ok := data["rx_bytes"].(string); ok {
			res.RxBytes = rx
		}
		p, err := res.ToOutput()
		if err != nil {
			fmt.Println("Can't convert TransmitAdapterOutputResource", err)
		}
		return p
	case apiModels.CommandTransmitTag:
		res := apiModels.TransmitTagOutputResource{}
		if rx, ok := data["rx_bytes"].(string); ok {
			res.RxBytes = rx
		}
		p, err := res.ToOutput()
		if err != nil {
			fmt.Println("Can't convert TransmitTagOutputResource", err)
		}
		return p
	case apiModels.CommandWriteNdef:
		return apiModels.WriteNdefOutput{}
	case apiModels.CommandReadNdef:
		res := apiModels.ReadNdefOutput{}

		if ndef, ok := data["ndef"].(map[string]interface{}); ok {
			if r, ok := ndef["read_only"].(bool); ok {
				res.Ndef.ReadOnly = r
			}

			if m, ok := data["message"].([]interface{}); ok {
				res.Ndef.Message = make([]ndefconv.NdefRecord, len(m))
				for i, msg := range m {
					md := msg.(map[string]interface{})

					resMsg := ndefconv.NdefRecord{}
					if t, ok := md["type"].(string); ok {
						recT, ok := ndefconv.StringToNdefRecordPayloadType(t)
						if !ok {
							fmt.Println("Can't parse String To Ndef Record Payload Type")
						}
						resMsg.Type = recT
					}

					if ndefData, ok := md["data"].(map[string]interface{}); ok {
						resMsg.Data = parseNdefPayloadResourceByNdefPayloadType(resMsg.Type, ndefData)
					}

					res.Ndef.Message[i] = resMsg
				}
			}
		}

		return res
	case apiModels.CommandFormatDefault:
		return apiModels.FormatDefaultOutput{}
	case apiModels.CommandLockPermanent:
		return apiModels.LockPermanentOutput{}
	case apiModels.CommandSetPassword:
		return apiModels.SetPasswordOutput{}
	case apiModels.CommandRemovePassword:
		return apiModels.RemovePasswordOutput{}
	case apiModels.CommandAuthPassword:
		res := apiModels.AuthPasswordOutputResource{}
		if ack, ok := data["ack"].(string); ok {
			res.Ack = ack
		}
		p, err := res.ToOutput()
		if err != nil {
			fmt.Println("Can't convert AuthPasswordOutputResource", err)
		}
		return p
	case apiModels.CommandGetDump:
		res := apiModels.GetDumpOutputResource{}
		if dump, ok := data["memory_dump"].([]interface{}); ok {
			sl := make([]apiModels.PageDumpResource, len(dump))
			for i, t := range dump {
				td := t.(map[string]interface{})

				if p, ok := td["page"].(string); ok {
					sl[i].Page = p
				}
				if p, ok := td["data"].(string); ok {
					sl[i].Data = p
				}
				if p, ok := td["info"].(string); ok {
					sl[i].Info = p
				}
			}
			res.MemoryDump = sl
		}
		p, err := res.ToOutput()
		if err != nil {
			fmt.Println("Can't convert GetDumpOutputResource", err)
		}
		return p
	}

	return nil
}
