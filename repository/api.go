package repository

import (
	"encoding/base64"
	"fmt"
	"github.com/taglme/nfc-cli/models"
	"github.com/taglme/nfc-cli/ndef"
	"github.com/taglme/nfc-goclient/pkg/client"
	apiModels "github.com/taglme/nfc-goclient/pkg/models"
	"github.com/taglme/nfc-goclient/pkg/ndefconv"
	"log"
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
	fmt.Printf("Server version: %s\n", i.Version)
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

func (s *ApiService) DeleteAdapterJobs(adapterId string) error {
	return s.client.Jobs.DeleteAll(adapterId)
}

func (s *ApiService) AddJobFromFile(adapterId string, filename string, p models.GenericJobParams) (int, error) {
	newJobs, err := s.readFromFile(filename)
	if err != nil {
		return 0, err
	}

	for _, newJob := range newJobs {
		if p.Expire != 60 {
			newJob.ExpireAfter = p.Expire
		}

		if len(p.JobName) > 0 {
			newJob.JobName = p.JobName
		}

		j, err := s.client.Jobs.Add(adapterId, newJob)
		if err != nil {
			return len(newJobs), err
		}

		fmt.Println("Job has been submitted:")
		s.printer.PrintJob(j)
		s.printer.Reset()
		fmt.Println("\nJob steps:")
		s.printer.PrintJobSteps(j.Steps)
		s.printer.Reset()
	}

	return len(newJobs), err
}

func (s *ApiService) addJob(nj *apiModels.NewJob, adapterId string, auth []byte, export bool) (*apiModels.Job, *apiModels.NewJob, error) {
	if auth != nil {
		nj.Steps = append(nj.Steps, apiModels.JobStepResource{})
		copy(nj.Steps[1:], nj.Steps)
		nj.Steps[0] = *s.getAuthJobStep(auth)
	}

	if export {
		fmt.Println("New job has been exported:")
		s.printer.PrintNewJob(*nj)
		s.printer.Reset()

		var steps []apiModels.JobStep
		for _, sr := range nj.Steps {
			s, err := sr.ToJobStep()
			if err != nil {
				log.Printf("Can't convert new job step resource to new job step: %s\n", err)
				continue
			}
			steps = append(steps, s)
		}
		s.printer.PrintJobSteps(steps)
		s.printer.Reset()
		return nil, nj, nil
	}

	j, err := s.client.Jobs.Add(adapterId, *nj)
	if err != nil {
		return &j, nj, err
	}

	fmt.Println("Job has been submitted:")
	s.printer.PrintJob(j)
	s.printer.Reset()
	fmt.Println("\nJob steps:")
	s.printer.PrintJobSteps(j.Steps)
	s.printer.Reset()

	return &j, nj, err
}

func (s *ApiService) AddGenericJob(p models.GenericJobParams) (*apiModels.Job, *apiModels.NewJob, error) {
	nj := apiModels.NewJob{
		JobName:     MapCliCmdToJobName[p.Cmd],
		Repeat:      p.Repeat,
		ExpireAfter: p.Expire,
		Steps:       MapCliCmdToApiJobSteps[p.Cmd],
	}

	if len(p.JobName) > 0 {
		nj.JobName = p.JobName
	}

	return s.addJob(&nj, p.AdapterId, p.Auth, p.Export)
}

func (s *ApiService) AddSetPwdJob(p models.GenericJobParams, password []byte) (*apiModels.Job, *apiModels.NewJob, error) {
	jobStep := apiModels.JobStep{
		Command: apiModels.CommandSetPassword,
		Params: apiModels.SetPasswordParams{
			Password: password,
		},
	}

	jobStepResource := jobStep.ToResource()

	nj := apiModels.NewJob{
		JobName:     "Set tag password",
		Repeat:      p.Repeat,
		ExpireAfter: p.Expire,
		Steps:       []apiModels.JobStepResource{jobStepResource},
	}

	if len(p.JobName) > 0 {
		nj.JobName = p.JobName
	}

	return s.addJob(&nj, p.AdapterId, p.Auth, p.Export)
}

func (s *ApiService) AddTransmitJob(p models.GenericJobParams, txBytes []byte, target string) (*apiModels.Job, *apiModels.NewJob, error) {
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

	if len(p.JobName) > 0 {
		nj.JobName = p.JobName
	}

	jobStepResource := jobStep.ToResource()

	nj.Repeat = p.Repeat
	nj.ExpireAfter = p.Expire
	nj.Steps = []apiModels.JobStepResource{jobStepResource}

	return s.addJob(&nj, p.AdapterId, p.Auth, p.Export)
}

func (s *ApiService) AddWriteJob(p models.GenericJobParams, r ndef.NdefPayload, protect bool) (*apiModels.Job, *apiModels.NewJob, error) {
	var nj apiModels.NewJob

	jobStep := apiModels.JobStep{
		Command: apiModels.CommandWriteNdef,
		Params: apiModels.WriteNdefParams{
			Message: []ndefconv.NdefRecord{
				r.ToRecord(),
			},
		},
	}

	nj.JobName = "Write tag"
	if len(p.JobName) > 0 {
		nj.JobName = p.JobName
	}
	nj.Repeat = p.Repeat
	nj.ExpireAfter = p.Expire
	nj.Steps = []apiModels.JobStepResource{jobStep.ToResource()}

	if protect {
		lockStep := apiModels.JobStep{
			Command: apiModels.CommandLockPermanent,
			Params:  apiModels.LockPermanentParams{},
		}

		nj.Steps = append(nj.Steps, lockStep.ToResource())
	}

	return s.addJob(&nj, p.AdapterId, p.Auth, p.Export)
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
