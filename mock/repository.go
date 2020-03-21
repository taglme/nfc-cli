package mock

import (
	"github.com/taglme/nfc-cli/models"
	"github.com/taglme/nfc-cli/ndef"
	"github.com/taglme/nfc-goclient/pkg/client"
	apiModels "github.com/taglme/nfc-goclient/pkg/models"
)

type MockedRepositoryService struct{}

func NewRepositoryService(c **client.Client) *MockedRepositoryService {
	return &MockedRepositoryService{}
}

func (s *MockedRepositoryService) GetVersion() (apiModels.AppInfo, error) {
	return apiModels.AppInfo{}, nil
}

func (s *MockedRepositoryService) GetAdapters(withOutput bool) ([]apiModels.Adapter, error) {
	return []apiModels.Adapter{
		{
			AdapterID: "mocked adapter id",
			Name:      "Mocker adapter name",
		},
	}, nil
}

func (s *MockedRepositoryService) GetJob(adapterId, id string) (apiModels.Job, error) {
	return apiModels.Job{}, nil
}

func (s *MockedRepositoryService) DeleteAdapterJobs(adapterId string) error {
	return nil
}

func (s *MockedRepositoryService) AddJobFromFile(adapterId string, filename string, p models.GenericJobParams) (int, error) {
	return 25, nil
}

func (s *MockedRepositoryService) addJob(nj *apiModels.NewJob, adapterId string, auth []byte, export bool) (*apiModels.Job, *apiModels.NewJob, error) {
	return &apiModels.Job{}, nj, nil
}

func (s *MockedRepositoryService) AddGenericJob(p models.GenericJobParams) (*apiModels.Job, *apiModels.NewJob, error) {
	nj := apiModels.NewJob{
		JobName:     "Job Name",
		Repeat:      p.Repeat,
		ExpireAfter: p.Expire,
		Steps:       []apiModels.JobStepResource{},
	}

	return s.addJob(&nj, p.AdapterId, p.Auth, p.Export)
}

func (s *MockedRepositoryService) AddSetPwdJob(p models.GenericJobParams, password []byte) (*apiModels.Job, *apiModels.NewJob, error) {
	nj := apiModels.NewJob{
		JobName:     "Job Name",
		Repeat:      p.Repeat,
		ExpireAfter: p.Expire,
		Steps:       []apiModels.JobStepResource{},
	}

	return s.addJob(&nj, p.AdapterId, p.Auth, p.Export)
}

func (s *MockedRepositoryService) AddTransmitJob(p models.GenericJobParams, txBytes []byte, target string) (*apiModels.Job, *apiModels.NewJob, error) {
	nj := apiModels.NewJob{
		JobName:     "Job Name",
		Repeat:      p.Repeat,
		ExpireAfter: p.Expire,
		Steps:       []apiModels.JobStepResource{},
	}

	return s.addJob(&nj, p.AdapterId, p.Auth, p.Export)
}

func (s *MockedRepositoryService) AddWriteJob(p models.GenericJobParams, r ndef.NdefPayload, protect bool) (*apiModels.Job, *apiModels.NewJob, error) {
	nj := apiModels.NewJob{
		JobName:     "Job Name",
		Repeat:      p.Repeat,
		ExpireAfter: p.Expire,
		Steps:       []apiModels.JobStepResource{},
	}

	return s.addJob(&nj, p.AdapterId, p.Auth, p.Export)
}

func (s *MockedRepositoryService) RunWsConnection(handler func(models.Event, interface{})) error {
	return nil
}

func (s *MockedRepositoryService) StopWsConnection() error {
	return nil
}
