package repository

import (
	"fmt"
	"github.com/taglme/nfc-cli/models"
	"github.com/taglme/nfc-client/pkg/client"
	apiModels "github.com/taglme/nfc-client/pkg/models"
)

type ApiService struct {
	client  *client.Client
	printer PrinterService
}

func New(c **client.Client) *ApiService {
	return &ApiService{
		client:  *c,
		printer: newPrinter(),
	}
}

func (s *ApiService) GetVersion() (apiModels.AppInfo, error) {
	i, err := s.client.About.Get()
	if err != nil {
		return i, err
	}
	s.printer.PrintAppInfo(i)
	s.printer.Reset()
	return i, err
}

func (s *ApiService) GetAdapters() ([]apiModels.Adapter, error) {
	a, err := s.client.Adapters.GetAll()
	if err != nil {
		return a, err
	}
	s.printer.PrintAdapters(a)
	s.printer.Reset()
	return a, err
}

func (s *ApiService) AddJob(cmd models.Command, adapterId string, repeat, expire int) (apiModels.Job, error) {
	j, err := s.client.Jobs.Add(adapterId, apiModels.NewJob{
		JobName:     MapCliCmdToApiCmd[cmd].String(),
		Repeat:      repeat,
		ExpireAfter: expire,
		Steps:       MapCliCmdToApiJobSteps[cmd],
	})
	if err != nil {
		return j, err
	}

	fmt.Println("Job has been submitted:")
	s.printer.PrintJob(j)
	s.printer.Reset()
	fmt.Println("\nJob steps:")
	s.printer.PrintJobSteps(j.Steps)
	s.printer.Reset()

	return j, err
}

var MapCliCmdToApiCmd = map[models.Command]apiModels.Command{
	models.CommandRead: apiModels.CommandReadNdef,
}

var MapCliCmdToApiJobSteps = map[models.Command][]apiModels.JobStepResource{
	models.CommandRead: {
		{
			Command: apiModels.CommandGetTags.String(),
			Params:  apiModels.GetTagsParamsResource{},
		}, {
			Command: apiModels.CommandReadNdef.String(),
			Params:  apiModels.ReadNdefParamsResource{},
		},
	},
}

var MapApiEventNameToCliEvent = map[apiModels.EventName]models.Event{
	apiModels.EventNameTagDiscovery:     models.EventTagDiscovery,
	apiModels.EventNameTagRelease:       models.EventTagRelease,
	apiModels.EventNameAdapterDiscovery: models.EventAdapterDiscovery,
	apiModels.EventNameAdapterRelease:   models.EventAdapterRelease,
	apiModels.EventNameJobSubmited:      models.EventJobSubmitted,
	apiModels.EventNameJobActivated:     models.EventJobActivated,
	apiModels.EventNameJobPended:        models.EventJobPended,
	apiModels.EventNameJobDeleted:       models.EventJobDeleted,
	apiModels.EventNameJobFinished:      models.EventJobFinished,
	apiModels.EventNameRunStarted:       models.EventRunStarted,
	apiModels.EventNameRunSuccess:       models.EventRunSuccess,
	apiModels.EventNameRunError:         models.EventRunError,
}
