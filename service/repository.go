package service

import (
	"github.com/taglme/nfc-cli/models"
	apiModels "github.com/taglme/nfc-client/pkg/models"
)

type ApiService interface {
	GetVersion() (apiModels.AppInfo, error)
	GetAdapters() ([]apiModels.Adapter, error)
	AddJob(models.Command, string, int, int, []byte) (apiModels.Job, error)
	AddSetPwdJob(models.Command, string, int, int, []byte, []byte) (apiModels.Job, error)
	RunWsConnection(func(models.Event)) error
	StopWsConnection() error
}
