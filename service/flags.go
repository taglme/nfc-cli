package service

import (
	"encoding/hex"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-cli/models"
	"github.com/urfave/cli/v2"
	"strings"
)

func (s *appService) parseHexString(hexStr string) (res []byte, err error) {
	if len(hexStr) <= 0 {
		return res, nil
	}

	decoded, err := hex.DecodeString(strings.Replace(hexStr, " ", "", -1))
	if err != nil {
		return res, errors.Wrap(err, "Can't decode hex string")
	}

	return decoded, nil
}

func (s *appService) getFlagsMap() map[string]cli.Flag {
	return map[string]cli.Flag{
		models.FlagHost: &cli.StringFlag{
			Name:        models.FlagHost,
			Value:       "127.0.0.1:3011",
			Usage:       "Target host and port",
			Destination: &s.host,
		},
		models.FlagAdapter: &cli.IntFlag{
			Name:        models.FlagAdapter,
			Value:       1,
			Usage:       "Adapter",
			Destination: &s.adapter,
		},
		models.FlagRepeat: &cli.IntFlag{
			Name:        models.FlagRepeat,
			Value:       1,
			Usage:       "Number of required repetitions of the task. Optional. If missing, the task is run once",
			Destination: &s.repeat,
		},
		models.FlagOutput: &cli.StringFlag{
			Name:        models.FlagOutput,
			Usage:       "File name for recording the results of the task. Optional. If there is no record of the results is not performed.",
			Destination: &s.output,
		},
		models.FlagAppend: &cli.BoolFlag{
			Name:        models.FlagAppend,
			Value:       false,
			Usage:       "Mode of writing the results to a file. Optional. If append = true, the results are added to the file. If absent or append = false after opening the file, its contents are cleared",
			Destination: &s.append,
		},
		models.FlagTimeout: &cli.IntFlag{
			Name:        models.FlagTimeout,
			Value:       60,
			Usage:       "Job timeout time in seconds. Optional. If absent equals 60",
			Destination: &s.timeout,
		},
		models.FlagFile: &cli.StringFlag{
			Name:        models.FlagFile,
			Usage:       "File name for loading data to form a command. Optional. If absent, data is formed from the arguments of the command. If present, then the command arguments are ignored, data is taken from the file.",
			Destination: &s.input,
			Required:    true,
		},
		models.FlagAuth: &cli.StringFlag{
			Name:        models.FlagAuth,
			Usage:       "An indication of the need for authorization before starting operations. The value of the argument is indicated as an array of bytes in hex format. Example \"03 AD F3 41\"",
			Destination: &s.auth,
		},
		models.FlagJobName: &cli.StringFlag{
			Name:        models.FlagJobName,
			Usage:       "Task name of the created task. Optional. If absent, then name created in accordance with the command used.",
			Destination: &s.jobName,
		},
		models.FlagExport: &cli.BoolFlag{
			Name:  models.FlagExport,
			Value: false,
			Usage: "Flag indicating the need to save it instead of sending a job to the server to the job file specified in the output parameter.",
		},
		models.FlagPwd: &cli.StringFlag{
			Name:     models.FlagPwd,
			Usage:    "Password to get an access to the memory of the NFC tag. The value of the argument is indicated as an array of bytes in hex format. Example \"03 AD F3 41\"",
			Required: true,
		},
		models.FlagTarget: &cli.StringFlag{
			Name:  models.FlagTarget,
			Usage: "Indicating to whom the byte array is sent. Optional. Can be tag or adapter. If not specified, the 'tag' value applies.",
			Value: "tag",
		},
		models.FlagTxBytes: &cli.StringFlag{
			Name:     models.FlagTxBytes,
			Usage:    "Array of bytes transmitted in hex format. Mandatory.",
			Required: true,
		},

		models.FlagNdefType: &cli.StringFlag{
			Name:     models.FlagNdefType,
			Usage:    "Indication of the type of record. Mandatory",
			Required: true,
		},
		models.FlagProtect: &cli.BoolFlag{
			Name:  models.FlagProtect,
			Usage: "The need to lock the label after recording. Optional.",
		},

		models.FlagNdefTypeRawId: &cli.StringFlag{
			Name:  models.FlagNdefTypeRawId,
			Usage: "NDEF raw type id field",
		},
		models.FlagNdefTypeRawTnf: &cli.IntFlag{
			Name:  models.FlagNdefTypeRawTnf,
			Value: -1,
			Usage: "NDEF raw type tnf field",
		},
		models.FlagNdefTypeType: &cli.StringFlag{
			Name:  models.FlagNdefTypeType,
			Usage: "NDEF raw type type field",
		},
		models.FlagNdefTypeRawPayload: &cli.StringFlag{
			Name:  models.FlagNdefTypeRawPayload,
			Usage: "NDEF raw type payload field",
		},

		models.FlagNdefTypeUrl: &cli.StringFlag{
			Name:  models.FlagNdefTypeUrl,
			Usage: "NDEF url type url field",
		},
		models.FlagNdefTypeText: &cli.StringFlag{
			Name:  models.FlagNdefTypeText,
			Usage: "NDEF text type text field",
		},
		models.FlagNdefTypeLang: &cli.StringFlag{
			Name:  models.FlagNdefTypeLang,
			Value: "English",
			Usage: "NDEF text type lang field",
		},
		models.FlagNdefUri: &cli.StringFlag{
			Name:  models.FlagNdefUri,
			Usage: "NDEF uri/vcard type uri field",
		},
		models.FlagNdefTypeAarPackage: &cli.StringFlag{
			Name:  models.FlagNdefTypeAarPackage,
			Usage: "NDEF aar type package name field",
		},
		models.FlagNdefTypePhone: &cli.StringFlag{
			Name:  models.FlagNdefTypePhone,
			Usage: "NDEF phone type phone number field",
		},
		models.FlagNdefTypeVcardAddressCity: &cli.StringFlag{
			Name:  models.FlagNdefTypeVcardAddressCity,
			Usage: "NDEF vcard type address city field",
		},
		models.FlagNdefTypeVcardAddressCountry: &cli.StringFlag{
			Name:  models.FlagNdefTypeVcardAddressCountry,
			Usage: "NDEF vcard type address country field",
		},
		models.FlagNdefTypeVcardAddressPostalCode: &cli.StringFlag{
			Name:  models.FlagNdefTypeVcardAddressPostalCode,
			Usage: "NDEF vcard type address postal code field",
		},
		models.FlagNdefTypeVcardAddressRegion: &cli.StringFlag{
			Name:  models.FlagNdefTypeVcardAddressRegion,
			Usage: "NDEF vcard type address region field",
		},
		models.FlagNdefTypeVcardAddressStreet: &cli.StringFlag{
			Name:  models.FlagNdefTypeVcardAddressStreet,
			Usage: "NDEF vcard type address street field",
		},
		models.FlagNdefTypeVcardEmail: &cli.StringFlag{
			Name:  models.FlagNdefTypeVcardEmail,
			Usage: "NDEF vcard type email field",
		},
		models.FlagNdefTypeVcardFirstName: &cli.StringFlag{
			Name:  models.FlagNdefTypeVcardFirstName,
			Usage: "NDEF vcard type first name field",
		},
		models.FlagNdefTypeVcardLastName: &cli.StringFlag{
			Name:  models.FlagNdefTypeVcardLastName,
			Usage: "NDEF vcard type last name field",
		},
		models.FlagNdefTypeVcardOrganization: &cli.StringFlag{
			Name:  models.FlagNdefTypeVcardOrganization,
			Usage: "NDEF vcard type organization field",
		},
		models.FlagNdefTypeVcardPhoneCell: &cli.StringFlag{
			Name:  models.FlagNdefTypeVcardPhoneCell,
			Usage: "NDEF vcard type phone cell field",
		},
		models.FlagNdefTypeVcardPhoneHome: &cli.StringFlag{
			Name:  models.FlagNdefTypeVcardPhoneHome,
			Usage: "NDEF vcard type phone home field",
		},
		models.FlagNdefTypeVcardPhoneWork: &cli.StringFlag{
			Name:  models.FlagNdefTypeVcardPhoneWork,
			Usage: "NDEF vcard type phone work field",
		},
		models.FlagNdefTypeTitle: &cli.StringFlag{
			Name:  models.FlagNdefTypeTitle,
			Usage: "NDEF vcard/poster type title field",
		},
		models.FlagNdefTypeVcardSite: &cli.StringFlag{
			Name:  models.FlagNdefTypeVcardSite,
			Usage: "NDEF vcard type site field",
		},
		models.FlagNdefTypeMimeFormat: &cli.StringFlag{
			Name:  models.FlagNdefTypeMimeFormat,
			Usage: "NDEF mime type format field",
		},
		models.FlagNdefTypeMimeContent: &cli.StringFlag{
			Name:  models.FlagNdefTypeMimeContent,
			Usage: "NDEF mime type content field",
		},
		models.FlagNdefTypeGeoLat: &cli.StringFlag{
			Name:  models.FlagNdefTypeGeoLat,
			Usage: "NDEF geo type latitude field",
		},
		models.FlagNdefTypeGeoLon: &cli.StringFlag{
			Name:  models.FlagNdefTypeGeoLon,
			Usage: "NDEF geo type longitude field",
		},
	}
}
