package actions

import (
	"github.com/taglme/nfc-cli/models"
	"github.com/taglme/nfc-client/pkg/client"
	apiModels "github.com/taglme/nfc-client/pkg/models"
)

type ActionService interface {
	GetVersion() (apiModels.AppInfo, error)
	GetAdapters() ([]apiModels.Adapter, error)
	AddJob(models.Command, string, int, int) (apiModels.Job, error)
}

type actionService struct {
	client *client.Client
}

func New(c *client.Client) ActionService {
	return &actionService{
		client: c,
	}
}

func (s *actionService) GetVersion() (apiModels.AppInfo, error) {
	return s.client.About.Get()
}

func (s *actionService) GetAdapters() ([]apiModels.Adapter, error) {
	return s.client.Adapters.GetAll()
}

func (s *actionService) AddJob(cmd models.Command, adapterId string, repeat, expire int) (apiModels.Job, error) {
	return s.client.Jobs.Add(adapterId, apiModels.NewJob{
		JobName:     MapCliCmdToApiCmd[cmd].String(),
		Repeat:      repeat,
		ExpireAfter: expire,
		Steps:       MapCliCmdToApiJobSteps[cmd],
	})
}

var MapCliCmdToApiCmd = map[models.Command]apiModels.Command{
	models.CommandRead: apiModels.CommandReadNdef,
}

var MapCliCmdToApiJobSteps = map[models.Command][]apiModels.JobStepResource{
	models.CommandRead: {
		{
			Command: apiModels.CommandGetTags.String(),
			Params:  apiModels.GetTagsParamsResource{},
		},{
			Command: apiModels.CommandReadNdef.String(),
			Params:  apiModels.ReadNdefParamsResource{},
		},
	},
}
