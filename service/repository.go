package service

import (
	"github.com/taglme/nfc-cli/models"
	"github.com/taglme/nfc-cli/ndef"
	apiModels "github.com/taglme/nfc-goclient/pkg/models"
)

type ApiService interface {
	GetVersion() (apiModels.AppInfo, error)
	GetAdapters() ([]apiModels.Adapter, error)
	GetJob(adapterId, id string) (apiModels.Job, error)
	DeleteAdapterJobs(adapterId string) error
	AddGenericJob(p models.GenericJobParams) (*apiModels.Job, *apiModels.NewJob, error)
	AddSetPwdJob(p models.GenericJobParams, password []byte) (*apiModels.Job, *apiModels.NewJob, error)
	AddTransmitJob(p models.GenericJobParams, txBytes []byte, target string) (*apiModels.Job, *apiModels.NewJob, error)
	AddWriteJob(p models.GenericJobParams, r ndef.NdefPayload, protect bool) (*apiModels.Job, *apiModels.NewJob, error)
	AddJobFromFile(adapterId string, filename string, p models.GenericJobParams) (int, error)
	RunWsConnection(handler func(models.Event, interface{})) error
	StopWsConnection() error
}
