package service

import (
	"github.com/taglme/nfc-cli/models"
	apiModels "github.com/taglme/nfc-client/pkg/models"
)

type ApiService interface {
	GetVersion() (apiModels.AppInfo, error)
	GetAdapters() ([]apiModels.Adapter, error)
	AddJob(models.Command, string, int, int, []byte) (apiModels.Job, error)
	AddSetPwdJob(string, int, int, []byte, []byte) (apiModels.Job, error)
	AddTransmitJob(string, int, int, []byte, []byte, string) (apiModels.Job, error)
	RunWsConnection(func(models.Event)) error
	StopWsConnection() error
}
