package repository

import (
	"encoding/base64"
	"fmt"
	"github.com/taglme/nfc-cli/models"
	"github.com/taglme/nfc-cli/ndef"
	"github.com/taglme/nfc-goclient/pkg/client"
	apiModels "github.com/taglme/nfc-goclient/pkg/models"
	"github.com/taglme/nfc-goclient/pkg/ndefconv"
)

type RepositoryService struct {
	client *client.Client
}

func New(c **client.Client) *RepositoryService {
	return &RepositoryService{
		client: *c,
	}
}

func (s *RepositoryService) GetVersion() (apiModels.AppInfo, error) {
	i, err := s.client.About.Get()
	if err != nil {
		return i, err
	}

	s.printAppInfo(i)
	return i, err
}

func (s *RepositoryService) GetAdapters(withOutput bool) ([]apiModels.Adapter, error) {
	a, err := s.client.Adapters.GetAll()
	if err != nil {
		return a, err
	}

	if withOutput {
		s.printAdapters(a)
	}

	return a, err
}

func (s *RepositoryService) GetJob(adapterId, id string) (apiModels.Job, error) {
	return s.client.Jobs.Get(adapterId, id)
}

func (s *RepositoryService) DeleteAdapterJobs(adapterId string) error {
	return s.client.Jobs.DeleteAll(adapterId)
}

func (s *RepositoryService) AddJobFromFile(adapterId string, filename string, p models.GenericJobParams) (int, error) {
	newJobs, err := s.readFromFile(filename)
	if err != nil {
		return 0, err
	}

	fmt.Printf("Loaded %d jobs.\n", len(newJobs))

	runs := 0
	for _, newJob := range newJobs {
		if p.Expire != 60 {
			newJob.ExpireAfter = p.Expire
		}

		if len(p.JobName) > 0 {
			newJob.JobName = p.JobName
		}

		runs += newJob.Repeat

		_, err := s.client.Jobs.Add(adapterId, newJob)
		if err != nil {
			return len(newJobs), err
		}
	}

	return runs, err
}

func (s *RepositoryService) addJob(nj *apiModels.NewJob, adapterId string, auth []byte, export bool) (*apiModels.Job, *apiModels.NewJob, error) {
	if auth != nil {
		nj.Steps = append(nj.Steps, apiModels.JobStepResource{})
		copy(nj.Steps[1:], nj.Steps)
		nj.Steps[0] = *s.getAuthJobStep(auth)
	}

	if export {
		fmt.Printf("Job %s: successfully exported.\n", nj.JobName)
		return nil, nj, nil
	}

	j, err := s.client.Jobs.Add(adapterId, *nj)
	if err != nil {
		return &j, nj, err
	}
	return &j, nj, err
}

func (s *RepositoryService) AddGenericJob(p models.GenericJobParams) (*apiModels.Job, *apiModels.NewJob, error) {
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

func (s *RepositoryService) AddSetPwdJob(p models.GenericJobParams, password []byte) (*apiModels.Job, *apiModels.NewJob, error) {
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

func (s *RepositoryService) AddTransmitJob(p models.GenericJobParams, txBytes []byte, target string) (*apiModels.Job, *apiModels.NewJob, error) {
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

func (s *RepositoryService) AddWriteJob(p models.GenericJobParams, r ndef.NdefPayload, protect bool) (*apiModels.Job, *apiModels.NewJob, error) {
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

func (s *RepositoryService) getAuthJobStep(pwd []byte) *apiModels.JobStepResource {
	encodedString := base64.StdEncoding.EncodeToString(pwd)

	return &apiModels.JobStepResource{
		Command: apiModels.CommandAuthPassword.String(),
		Params:  apiModels.AuthPasswordParamsResource{Password: encodedString},
	}
}
