package service

import (
	"github.com/taglme/nfc-cli/models"
	"github.com/taglme/nfc-cli/ndef"
	apiModels "github.com/taglme/nfc-client/pkg/models"
)

type ApiService interface {
	GetVersion() (apiModels.AppInfo, error)
	GetAdapters() ([]apiModels.Adapter, error)
	AddGenericJob(p models.GenericJobParams) (*apiModels.Job, *apiModels.NewJob, error)
	AddSetPwdJob(p models.GenericJobParams, password []byte) (*apiModels.Job, *apiModels.NewJob, error)
	AddTransmitJob(p models.GenericJobParams, txBytes []byte, target string) (*apiModels.Job, *apiModels.NewJob, error)
	AddWriteJob(p models.GenericJobParams, r ndef.NdefPayload, protect bool) (*apiModels.Job, *apiModels.NewJob, error)
	RunWsConnection(handler func(models.Event, interface{})) error
	StopWsConnection() error
}
