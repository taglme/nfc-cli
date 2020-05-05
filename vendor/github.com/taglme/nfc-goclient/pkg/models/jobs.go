package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	uuid "github.com/nu7hatch/gouuid"
)

type Job struct {
	JobID       string
	JobName     string
	Status      JobStatus
	AdapterID   string
	AdapterName string
	Repeat      int
	TotalRuns   int
	SuccessRuns int
	ErrorRuns   int
	ExpireAfter int
	Steps       []JobStep
	CreatedAt   time.Time
}

type JobResource struct {
	JobID       string            `json:"job_id"`
	Kind        string            `json:"kind"`
	Href        string            `json:"href"`
	JobName     string            `json:"job_name"`
	Status      string            `json:"status"`
	AdapterID   string            `json:"adapter_id"`
	AdapterName string            `json:"adapter_name"`
	Repeat      int               `json:"repeat"`
	TotalRuns   int               `json:"total_runs"`
	SuccessRuns int               `json:"success_runs"`
	ErrorRuns   int               `json:"error_runs"`
	ExpireAfter int               `json:"expire_after"`
	Steps       []JobStepResource `json:"steps"`
	CreatedAt   string            `json:"created_at"`
}

type NewJob struct {
	JobName     string            `json:"job_name" binding:"required"`
	Repeat      int               `json:"repeat"`
	ExpireAfter int               `json:"expire_after" binding:"required"`
	Steps       []JobStepResource `json:"steps" binding:"required"`
}

type JobStatusUpdate struct {
	Status string `json:"status" binding:"required"`
}

func (j Job) ToResource() JobResource {
	var jobStepResources []JobStepResource
	for _, jobStep := range j.Steps {
		jobStepResources = append(jobStepResources, jobStep.ToResource())
	}
	resource := JobResource{
		JobID:       j.JobID,
		Kind:        "Job",
		Href:        fmt.Sprintf(`/adapters/%s/jobs/%s`, j.AdapterID, j.JobID),
		JobName:     j.JobName,
		Status:      j.Status.String(),
		AdapterID:   j.AdapterID,
		AdapterName: j.AdapterName,
		Repeat:      j.Repeat,
		TotalRuns:   j.TotalRuns,
		SuccessRuns: j.SuccessRuns,
		ErrorRuns:   j.ErrorRuns,
		ExpireAfter: j.ExpireAfter,
		Steps:       jobStepResources,
		CreatedAt:   j.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
	}
	return resource
}

func (nj NewJob) ToJob(adapterID string, adapterName string) (Job, error) {
	id, _ := uuid.NewV4()
	createdAt := time.Now().UTC()
	var jobSteps []JobStep
	for _, jobStepResource := range nj.Steps {
		jobStep, err := jobStepResource.ToJobStep()
		if err != nil {
			return Job{}, err
		}
		jobSteps = append(jobSteps, jobStep)
	}

	j := Job{
		JobID:       id.String(),
		JobName:     nj.JobName,
		Status:      JobStatusPending,
		AdapterID:   adapterID,
		AdapterName: adapterName,
		Repeat:      nj.Repeat,
		TotalRuns:   0,
		SuccessRuns: 0,
		ErrorRuns:   0,
		ExpireAfter: nj.ExpireAfter,
		Steps:       jobSteps,
		CreatedAt:   createdAt,
	}

	return j, nil
}

func (j JobResource) ToJob() (job Job, err error) {
	t, err := time.Parse(time.RFC3339, j.CreatedAt)
	if err != nil {
		return job, errors.Wrap(err, "Can't parse job resource created at")
	}

	s, ok := StringToJobStatus(j.Status)
	if !ok {
		return job, errors.Wrap(err, "Can't convert job resource status")
	}

	job = Job{
		JobID:       j.JobID,
		JobName:     j.JobName,
		Status:      s,
		AdapterID:   j.AdapterID,
		AdapterName: j.AdapterName,
		Repeat:      j.Repeat,
		TotalRuns:   j.TotalRuns,
		SuccessRuns: j.SuccessRuns,
		ErrorRuns:   j.ErrorRuns,
		ExpireAfter: j.ExpireAfter,
		CreatedAt:   t,
	}
	var jSteps []JobStep
	for _, s := range j.Steps {
		step, err := s.ToJobStep()

		if err != nil {
			return job, errors.Wrap(err, "Can't convert job step to the step model")
		}

		jSteps = append(jSteps, step)
	}

	job.Steps = jSteps

	return job, nil
}

type JobListResource struct {
	Total  int
	Length int
	Limit  int
	Offset int
	Items  []JobResource
}

type JobStatus int

const (
	JobStatusPending JobStatus = iota + 1
	JobStatusActive
)

func StringToJobStatus(s string) (JobStatus, bool) {
	switch s {
	case JobStatusPending.String():
		return JobStatusPending, true
	case JobStatusActive.String():
		return JobStatusActive, true
	}
	return 0, false
}

func (jobStatus JobStatus) String() string {
	names := [...]string{
		"unknown",
		"pending",
		"active",
	}
	if jobStatus < JobStatusPending || jobStatus > JobStatusActive {
		return names[0]
	}
	return names[jobStatus]
}

type JobStep struct {
	Command Command
	Params  CommandParams
}

type JobStepResource struct {
	Command string                `json:"command" binding:"required"`
	Params  CommandParamsResource `json:"params" binding:"required"`
}

func (list JobListResource) GetPaginationInfo() PageInfo {
	return PageInfo{
		Total:  list.Total,
		Length: list.Length,
		Limit:  list.Limit,
		Offset: list.Offset,
	}
}

func (jobStep *JobStepResource) UnmarshalJSON(data []byte) error {

	var obj map[string]interface{}
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}

	t, ok := obj["command"].(string)
	if !ok {
		return errors.New("Job step should have 'command' field")
	}

	command, isValid := StringToCommand(t)
	if !isValid {
		return errors.New("Job step have not valid command name")
	}
	jobStep.Command = t

	_, ok = obj["params"]

	if !ok {
		return errors.New("Job step should have 'params' field")
	}

	var paramsBytes []byte
	paramsBytes, _ = json.Marshal(obj["params"])
	switch command {
	case CommandGetTags:
		r := GetTagsParamsResource{}
		err := json.Unmarshal(paramsBytes, &r)
		if err != nil {
			return err
		}
		jobStep.Params = r
	case CommandTransmitAdapter:
		r := TransmitAdapterParamsResource{}
		err := json.Unmarshal(paramsBytes, &r)
		if err != nil {
			return err
		}
		if r.TxBytes == "" {
			return errors.New("Command 'transmit_adapter'should have not empty 'tx_bytes' param")
		}
		jobStep.Params = r
	case CommandTransmitTag:
		r := TransmitTagParamsResource{}
		err := json.Unmarshal(paramsBytes, &r)
		if err != nil {
			return err
		}
		if r.TxBytes == "" {
			return errors.New("Command 'transmit_tag'should have not empty 'tx_bytes' param")
		}
		jobStep.Params = r
	case CommandWriteNdef:
		r := WriteNdefParamsResource{}
		err := json.Unmarshal(paramsBytes, &r)
		if err != nil {
			return err
		}
		jobStep.Params = r
	case CommandReadNdef:
		r := ReadNdefParamsResource{}
		err := json.Unmarshal(paramsBytes, &r)
		if err != nil {
			return err
		}
		jobStep.Params = r
	case CommandFormatDefault:
		r := FormatDefaultParamsResource{}
		err := json.Unmarshal(paramsBytes, &r)
		if err != nil {
			return err
		}
		jobStep.Params = r
	case CommandLockPermanent:
		r := LockPermanentParamsResource{}
		err := json.Unmarshal(paramsBytes, &r)
		if err != nil {
			return err
		}
		jobStep.Params = r
	case CommandSetPassword:
		r := SetPasswordParamsResource{}
		err := json.Unmarshal(paramsBytes, &r)
		if err != nil {
			return err
		}
		if r.Password == "" {
			return errors.New("Command 'set_pasword' should have not empty 'password' param")
		}

		jobStep.Params = r
	case CommandRemovePassword:
		r := RemovePasswordParamsResource{}
		err := json.Unmarshal(paramsBytes, &r)
		if err != nil {
			return err
		}
		jobStep.Params = r
	case CommandAuthPassword:
		r := AuthPasswordParamsResource{}
		err := json.Unmarshal(paramsBytes, &r)
		if err != nil {
			return err
		}
		if r.Password == "" {
			return errors.New("Command 'auth_password' should have not empty 'password' param")
		}
		jobStep.Params = r
	case CommandGetDump:
		r := GetDumpParamsResource{}
		err := json.Unmarshal(paramsBytes, &r)
		if err != nil {
			return err
		}
		jobStep.Params = r

	case CommandSetLocale:
		r := SetLocaleParamsResource{}
		err := json.Unmarshal(paramsBytes, &r)
		if err != nil {
			return err
		}
		if r.Locale == "" {
			return errors.New("Command 'set_locale' should have not empty 'locale' param")
		}
		jobStep.Params = r
	}

	return nil

}

func (jobStep JobStep) ToResource() JobStepResource {
	var params CommandParamsResource
	if jobStep.Params != nil {
		params = jobStep.Params.ToResource()
	}

	resource := JobStepResource{
		Command: jobStep.Command.String(),
		Params:  params,
	}

	return resource
}
func (jobStepResource JobStepResource) ToJobStep() (JobStep, error) {
	var params CommandParams
	var err error
	if jobStepResource.Params != nil {
		params, err = jobStepResource.Params.ToParams()
		if err != nil {
			return JobStep{}, err
		}
	}

	command, _ := StringToCommand(jobStepResource.Command)
	resource := JobStep{
		Command: command,
		Params:  params,
	}
	return resource, nil
}

func (jobStep JobStep) ToStepResult() StepResult {
	stepResult := StepResult{
		Command: jobStep.Command,
		Params:  jobStep.Params,
	}
	return stepResult
}
