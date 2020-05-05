package repository

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/taglme/nfc-cli/models"
	"github.com/taglme/nfc-cli/ndef"
	"github.com/taglme/nfc-goclient/pkg/client"
	apiModels "github.com/taglme/nfc-goclient/pkg/models"
	"github.com/taglme/nfc-goclient/pkg/ndefconv"
)

func TestNew(t *testing.T) {
	nfc := client.New("url")
	rep := New(&nfc)

	assert.NotNil(t, rep)
}

func TestRepositoryService_getAuthJobStep(t *testing.T) {
	nfc := client.New("url")
	rep := New(&nfc)

	js := rep.getAuthJobStep([]byte{0xa6, 0x12, 0x66, 0xBA})

	assert.Equal(t, apiModels.CommandAuthPassword.String(), js.Command)
	assert.Equal(t, apiModels.AuthPasswordParamsResource{Password: "phJmug=="}, js.Params)
}

func TestRepositoryService_addJob(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/adapters/adapterId/jobs", req.URL.String())
		resp, err := json.Marshal(apiModels.JobResource{
			JobID:       "id",
			JobName:     "Job Name",
			AdapterID:   "adapterId",
			AdapterName: "adname",
			CreatedAt:   "2006-01-02T15:04:05Z",
			Status:      apiModels.JobStatusPending.String(),
			Steps: []apiModels.JobStepResource{{
				Command: apiModels.CommandGetDump.String(),
				Params:  apiModels.GetDumpParamsResource{},
			}},
		})
		if err != nil {
			log.Fatal("Can't marshall test model\n", err)
		}
		rw.WriteHeader(200)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return er\n", err)
		}
	}))

	defer server.Close()

	nfc := client.New(strings.Replace(server.URL, "http://", "", -1))
	rep := New(&nfc)

	nj := apiModels.NewJob{
		JobName:     "Job Name",
		Repeat:      2,
		ExpireAfter: 30,
		Steps: []apiModels.JobStepResource{{
			Command: apiModels.CommandGetDump.String(),
			Params:  apiModels.GetDumpParamsResource{},
		}},
	}

	job, respnj, err := rep.addJob(&nj, "adapterId", nil, false)
	if err != nil {
		t.Error(err)
		log.Fatal(err)
	}

	assert.Equal(t, "id", job.JobID)
	assert.Equal(t, "Job Name", job.JobName)
	assert.Equal(t, "adname", job.AdapterName)
	assert.Equal(t, "2006-01-02T15:04:05Z", job.CreatedAt.Format(time.RFC3339))
	assert.Equal(t, apiModels.JobStatusPending, job.Status)
	assert.Equal(t, apiModels.CommandGetDump, job.Steps[0].Command)
	assert.Equal(t, apiModels.CommandGetDump.String(), respnj.Steps[0].Command)
}

func TestRepositoryService_addJob_withAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/adapters/adapterId/jobs", req.URL.String())
		resp, err := json.Marshal(apiModels.JobResource{
			JobID:       "id",
			JobName:     "Job Name",
			AdapterID:   "adapterId",
			AdapterName: "adname",
			CreatedAt:   "2006-01-02T15:04:05Z",
			Status:      apiModels.JobStatusPending.String(),
			Steps: []apiModels.JobStepResource{
				{
					Command: apiModels.CommandAuthPassword.String(),
					Params:  apiModels.AuthPasswordParamsResource{Password: "phJmug=="},
				},
				{
					Command: apiModels.CommandGetDump.String(),
					Params:  apiModels.GetDumpParamsResource{},
				},
			},
		})
		if err != nil {
			log.Fatal("Can't marshall test model\n", err)
		}
		rw.WriteHeader(200)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return er\n", err)
		}
	}))

	defer server.Close()

	nfc := client.New(strings.Replace(server.URL, "http://", "", -1))
	rep := New(&nfc)

	nj := apiModels.NewJob{
		JobName:     "Job Name",
		Repeat:      2,
		ExpireAfter: 30,
		Steps: []apiModels.JobStepResource{{
			Command: apiModels.CommandGetDump.String(),
			Params:  apiModels.GetDumpParamsResource{},
		}},
	}

	job, respnj, err := rep.addJob(&nj, "adapterId", []byte{0xa6, 0x12, 0x66, 0xBA}, false)
	if err != nil {
		t.Error(err)
		log.Fatal(err)
	}

	assert.Equal(t, "id", job.JobID)
	assert.Equal(t, "Job Name", job.JobName)
	assert.Equal(t, "adname", job.AdapterName)
	assert.Equal(t, "2006-01-02T15:04:05Z", job.CreatedAt.Format(time.RFC3339))
	assert.Equal(t, apiModels.JobStatusPending, job.Status)
	assert.Equal(t, apiModels.CommandAuthPassword, job.Steps[0].Command)
	assert.Equal(t, "a6 12 66 ba", job.Steps[0].Params.String())
	assert.Equal(t, apiModels.CommandGetDump, job.Steps[1].Command)
	assert.Equal(t, apiModels.CommandAuthPassword.String(), respnj.Steps[0].Command)
	assert.Equal(t, apiModels.AuthPasswordParamsResource{Password: "phJmug=="}, respnj.Steps[0].Params)
	assert.Equal(t, apiModels.CommandGetDump.String(), respnj.Steps[1].Command)
}

func TestRepositoryService_addJob_exported(t *testing.T) {
	nfc := client.New("")
	rep := New(&nfc)

	nj := apiModels.NewJob{
		JobName:     "Job Name",
		Repeat:      2,
		ExpireAfter: 30,
		Steps: []apiModels.JobStepResource{{
			Command: apiModels.CommandGetDump.String(),
			Params:  apiModels.GetDumpParamsResource{},
		}},
	}

	job, respnj, err := rep.addJob(&nj, "adapterId", []byte{0xa6, 0x12, 0x66, 0xBA}, true)
	if err != nil {
		t.Error(err)
		log.Fatal(err)
	}

	assert.Nil(t, job)
	assert.Equal(t, apiModels.CommandAuthPassword.String(), respnj.Steps[0].Command)
	assert.Equal(t, apiModels.AuthPasswordParamsResource{Password: "phJmug=="}, respnj.Steps[0].Params)
	assert.Equal(t, apiModels.CommandGetDump.String(), respnj.Steps[1].Command)
}

func TestRepositoryService_AddGenericJob(t *testing.T) {
	p := models.GenericJobParams{
		Cmd:       models.CommandDump,
		AdapterId: "adapterId",
		Repeat:    3,
		Expire:    60,
		Auth:      []byte{0xa6, 0x12, 0x66, 0xBA},
		Export:    true,
		JobName:   "Dump tag",
	}

	nfc := client.New("url")
	rep := New(&nfc)

	_, nj, err := rep.AddGenericJob(p)
	if err != nil {
		t.Error(err)
		log.Fatal(err)
	}

	assert.Equal(t, "Dump tag", nj.JobName)
	assert.Equal(t, apiModels.CommandAuthPassword.String(), nj.Steps[0].Command)
	assert.Equal(t, apiModels.AuthPasswordParamsResource{Password: "phJmug=="}, nj.Steps[0].Params)
	assert.Equal(t, apiModels.CommandGetDump.String(), nj.Steps[1].Command)
}

func TestRepositoryService_AddSetPwdJob(t *testing.T) {
	p := models.GenericJobParams{
		Cmd:       models.CommandSetpwd,
		AdapterId: "adapterId",
		Repeat:    3,
		Expire:    60,
		Auth:      []byte{0xa6, 0x12, 0x66, 0xBA},
		Export:    true,
		JobName:   "job name",
	}

	nfc := client.New("url")
	rep := New(&nfc)

	_, nj, err := rep.AddSetPwdJob(p, []byte{0xa6, 0x12, 0x66, 0xBA})
	if err != nil {
		t.Error(err)
		log.Fatal(err)
	}

	assert.Equal(t, "job name", nj.JobName)
	assert.Equal(t, apiModels.CommandAuthPassword.String(), nj.Steps[0].Command)
	assert.Equal(t, apiModels.AuthPasswordParamsResource{Password: "phJmug=="}, nj.Steps[0].Params)
	assert.Equal(t, apiModels.CommandSetPassword.String(), nj.Steps[1].Command)
	assert.Equal(t, apiModels.SetPasswordParamsResource{Password: "phJmug=="}, nj.Steps[1].Params)
}

func TestRepositoryService_AddTransmitJob_Adapter(t *testing.T) {
	p := models.GenericJobParams{
		Cmd:       models.CommandTransmit,
		AdapterId: "adapterId",
		Repeat:    3,
		Expire:    60,
		Auth:      []byte{0xa6, 0x12, 0x66, 0xBA},
		Export:    true,
		JobName:   "job name",
	}

	nfc := client.New("url")
	rep := New(&nfc)

	_, nj, err := rep.AddTransmitJob(p, []byte{0xa6, 0x12, 0x66, 0xBA}, "adapter")
	if err != nil {
		t.Error(err)
		log.Fatal(err)
	}

	assert.Equal(t, "job name", nj.JobName)
	assert.Equal(t, apiModels.CommandAuthPassword.String(), nj.Steps[0].Command)
	assert.Equal(t, apiModels.AuthPasswordParamsResource{Password: "phJmug=="}, nj.Steps[0].Params)
	assert.Equal(t, apiModels.CommandTransmitAdapter.String(), nj.Steps[1].Command)
	assert.Equal(t, apiModels.TransmitAdapterParamsResource{TxBytes: "phJmug=="}, nj.Steps[1].Params)
}

func TestRepositoryService_AddTransmitJob_Tag(t *testing.T) {
	p := models.GenericJobParams{
		Cmd:       models.CommandTransmit,
		AdapterId: "adapterId",
		Repeat:    3,
		Expire:    60,
		Auth:      []byte{0xa6, 0x12, 0x66, 0xBA},
		Export:    true,
		JobName:   "job name",
	}

	nfc := client.New("url")
	rep := New(&nfc)

	_, nj, err := rep.AddTransmitJob(p, []byte{0xa6, 0x12, 0x66, 0xBA}, "tag")
	if err != nil {
		t.Error(err)
		log.Fatal(err)
	}

	assert.Equal(t, "job name", nj.JobName)
	assert.Equal(t, apiModels.CommandAuthPassword.String(), nj.Steps[0].Command)
	assert.Equal(t, apiModels.AuthPasswordParamsResource{Password: "phJmug=="}, nj.Steps[0].Params)
	assert.Equal(t, apiModels.CommandTransmitTag.String(), nj.Steps[1].Command)
	assert.Equal(t, apiModels.TransmitTagParamsResource{TxBytes: "phJmug=="}, nj.Steps[1].Params)
}

func TestRepositoryService_AddWriteJob(t *testing.T) {
	p := models.GenericJobParams{
		Cmd:       models.CommandWrite,
		AdapterId: "adapterId",
		Repeat:    3,
		Expire:    60,
		Auth:      nil,
		Export:    true,
		JobName:   "job name",
	}

	nfc := client.New("url")
	rep := New(&nfc)

	_, nj, err := rep.AddWriteJob(p, ndef.NdefRecordPayloadUrl{Url: "http://url"}, false)
	if err != nil {
		t.Error(err)
		log.Fatal(err)
	}

	expectedNdef := apiModels.WriteNdefParamsResource{
		Message: []ndefconv.NdefRecordResource{{
			Type: ndefconv.NdefRecordPayloadTypeUrl.String(),
			Data: ndefconv.NdefRecordPayloadUrlResource{Url: "http://url"},
		}},
	}

	assert.Equal(t, "job name", nj.JobName)
	assert.Equal(t, apiModels.CommandWriteNdef.String(), nj.Steps[0].Command)
	assert.Equal(t, expectedNdef, nj.Steps[0].Params)
}

func TestRepositoryService_AddWriteJob_Protect(t *testing.T) {
	p := models.GenericJobParams{
		Cmd:       models.CommandWrite,
		AdapterId: "adapterId",
		Repeat:    3,
		Expire:    60,
		Auth:      nil,
		Export:    true,
		JobName:   "job name",
	}

	nfc := client.New("url")
	rep := New(&nfc)

	_, nj, err := rep.AddWriteJob(p, ndef.NdefRecordPayloadUrl{Url: "http://url"}, true)
	if err != nil {
		t.Error(err)
		log.Fatal(err)
	}

	expectedNdef := apiModels.WriteNdefParamsResource{
		Message: []ndefconv.NdefRecordResource{{
			Type: ndefconv.NdefRecordPayloadTypeUrl.String(),
			Data: ndefconv.NdefRecordPayloadUrlResource{Url: "http://url"},
		}},
	}

	assert.Equal(t, "job name", nj.JobName)
	assert.Equal(t, apiModels.CommandWriteNdef.String(), nj.Steps[0].Command)
	assert.Equal(t, expectedNdef, nj.Steps[0].Params)
	assert.Equal(t, apiModels.CommandLockPermanent.String(), nj.Steps[1].Command)
	assert.Equal(t, apiModels.LockPermanentParamsResource{}, nj.Steps[1].Params)
}

func TestRepositoryService_AddJobFromFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/adapters/adapterId/jobs", req.URL.String())
		resp, err := json.Marshal(apiModels.JobResource{
			JobID:       "id",
			JobName:     "Job Name",
			AdapterID:   "adapterId",
			AdapterName: "adname",
			CreatedAt:   "2006-01-02T15:04:05Z",
			Status:      apiModels.JobStatusPending.String(),
			Steps: []apiModels.JobStepResource{
				{
					Command: apiModels.CommandAuthPassword.String(),
					Params:  apiModels.AuthPasswordParamsResource{Password: "phJmug=="},
				},
				{
					Command: apiModels.CommandGetDump.String(),
					Params:  apiModels.GetDumpParamsResource{},
				},
			},
		})
		if err != nil {
			log.Fatal("Can't marshall test model\n", err)
		}
		rw.WriteHeader(200)
		_, err = rw.Write(resp)
		if err != nil {
			log.Fatal("Can't return er\n", err)
		}
	}))

	defer server.Close()

	nfc := client.New(strings.Replace(server.URL, "http://", "", -1))
	rep := New(&nfc)

	p := models.GenericJobParams{
		Cmd:       models.CommandWrite,
		AdapterId: "adapterId",
		Repeat:    3,
		Expire:    60,
		Auth:      nil,
		Export:    false,
		JobName:   "job name",
	}
	amountOfRuns, err := rep.AddJobFromFile("adapterId", "reader_test_file.json", p)
	if err != nil {
		t.Error(err)
		log.Fatal(err)
	}

	assert.Equal(t, 3, amountOfRuns)
}
