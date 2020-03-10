package repository

import (
	"encoding/base64"
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

func (s *ApiService) addJob(nj apiModels.NewJob, adapterId string, pwd []byte) (apiModels.Job, error) {
	if pwd != nil {
		nj.Steps = append(nj.Steps, apiModels.JobStepResource{})
		copy(nj.Steps[1:], nj.Steps)
		nj.Steps[0] = *s.getAuthJobStep(pwd)
	}

	j, err := s.client.Jobs.Add(adapterId, nj)
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

func (s *ApiService) AddGenericJob(cmd models.Command, adapterId string, repeat, expire int, pwd []byte) (apiModels.Job, error) {
	nj := apiModels.NewJob{
		JobName:     MapCliCmdToJobName[cmd],
		Repeat:      repeat,
		ExpireAfter: expire,
		Steps:       MapCliCmdToApiJobSteps[cmd],
	}

	return s.addJob(nj, adapterId, pwd)
}

func (s *ApiService) AddSetPwdJob(adapterId string, repeat, expire int, auth, password []byte) (apiModels.Job, error) {
	jobStep := apiModels.JobStep{
		Command: apiModels.CommandSetPassword,
		Params: apiModels.SetPasswordParams{
			Password: password,
		},
	}

	jobStepResource := jobStep.ToResource()

	nj := apiModels.NewJob{
		JobName:     "Set tag password",
		Repeat:      repeat,
		ExpireAfter: expire,
		Steps:       []apiModels.JobStepResource{jobStepResource},
	}

	return s.addJob(nj, adapterId, auth)
}

func (s *ApiService) AddTransmitJob(adapterId string, repeat, expire int, auth, txBytes []byte, target string) (apiModels.Job, error) {
	var nj apiModels.NewJob
	var jobStep apiModels.JobStep

	if target == "adapter" {
		nj.JobName = "Transmit adapter"
		jobStep = apiModels.JobStep{
			Command: apiModels.CommandTransmitAdapter,
			Params: apiModels.TransmitAdapterParams{
				TxBytes: txBytes,
			},
		}
	} else {
		nj.JobName = "Transmit tag"
		jobStep = apiModels.JobStep{
			Command: apiModels.CommandTransmitTag,
			Params: apiModels.TransmitTagParams{
				TxBytes: txBytes,
			},
		}
	}

	jobStepResource := jobStep.ToResource()

	nj.Repeat = repeat
	nj.ExpireAfter = expire
	nj.Steps = []apiModels.JobStepResource{jobStepResource}

	return s.addJob(nj, adapterId, auth)
}

var MapCliCmdToJobName = map[models.Command]string{
	models.CommandRead:   "Read tag",
	models.CommandDump:   "Dump tag",
	models.CommandLock:   "Lock tag",
	models.CommandFormat: "Format tag",
	models.CommandRmpwd:  "Remove tag password",
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
	models.CommandDump: {
		{
			Command: apiModels.CommandGetDump.String(),
			Params:  apiModels.GetDumpParamsResource{},
		},
	},
	models.CommandLock: {
		{
			Command: apiModels.CommandLockPermanent.String(),
			Params:  apiModels.LockPermanentParamsResource{},
		},
	},
	models.CommandFormat: {
		{
			Command: apiModels.CommandFormatDefault.String(),
			Params:  apiModels.FormatDefaultParamsResource{},
		},
	},
	models.CommandRmpwd: {
		{
			Command: apiModels.CommandRemovePassword.String(),
			Params:  apiModels.RemovePasswordParamsResource{},
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

func (s *ApiService) getAuthJobStep(pwd []byte) *apiModels.JobStepResource {
	encodedString := base64.StdEncoding.EncodeToString(pwd)

	return &apiModels.JobStepResource{
		Command: apiModels.CommandAuthPassword.String(),
		Params:  apiModels.AuthPasswordParamsResource{Password: encodedString},
	}
}
