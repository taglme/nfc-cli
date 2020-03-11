package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/pkg/errors"
	"github.com/taglme/nfc-goclient/pkg/models"
)

//JobFilter defines filter for job list
type JobFilter struct {
	Status  *models.JobStatus
	SortBy  *string
	SortDir *string
	Limit   *int
	Offset  *int
}

//JobService job service interface
type JobService interface {
	GetAll(adapterID string) ([]models.Job, models.PageInfo, error)
	GetFiltered(adapterID string, filter JobFilter) ([]models.Job, models.PageInfo, error)
	Delete(adapterID string, jobID string) error
	DeleteAll(adapterID string) error
	Get(adapterID string, jobID string) (models.Job, error)
	Add(adapterID string, job models.NewJob) (models.Job, error)
	UpdateStatus(adapterID string, jobID string, status models.JobStatus) (models.Job, error)
}

type jobService struct {
	url      string
	basePath string
	path     string
	client   *http.Client
}

func newJobService(c *http.Client, url string) JobService {
	return &jobService{
		url:      url,
		client:   c,
		path:     "/jobs",
		basePath: "/adapters",
	}
}

// Get Job list for adapter with all details
func (s *jobService) GetAll(adapterID string) ([]models.Job, models.PageInfo, error) {
	return s.GetFiltered(adapterID, JobFilter{})
}

// Get Job list for adapter with all details
// adapterId – Unique identifier in form of UUID representing a specific adapter.
// filter.status – Jobs' status filter
// filter.limit – Limit number of jobs in response.
// filter.offset – Offset from start of list.
// filter.sortBy – Sort field for list.
// filter.sortDir – Sort direction for list
func (s *jobService) GetFiltered(adapterID string, filter JobFilter) (jobs []models.Job, pagInfo models.PageInfo, err error) {
	targetUrl := s.url + s.basePath + "/" + adapterID + s.path + buildJobsQueryParams(filter)
	resp, err := s.client.Get(targetUrl)
	if err != nil {
		return jobs, pagInfo, errors.Wrap(err, "Can't get jobs\n")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return jobs, pagInfo, errors.Wrap(err, "Can't convert jobs to byte slice\n")
	}

	err = handleHttpResponseCode(resp.StatusCode, body)
	if err != nil {
		return jobs, pagInfo, errors.Wrap(err, "Error in fetching jobs\n")
	}

	var jListResource models.JobListResource
	err = json.Unmarshal(body, &jListResource)
	if err != nil {
		return jobs, pagInfo, errors.Wrap(err, "Can't unmarshal jobs response\n")
	}

	jobs = make([]models.Job, len(jListResource.Items))
	for i, e := range jListResource.Items {
		jobs[i], err = e.ToJob()
		if err != nil {
			return jobs, pagInfo, errors.Wrap(err, "Can't convert job resource to job model\n")
		}
	}

	return jobs, jListResource.GetPaginationInfo(), nil
}

// Get Job list for adapter with all details
// adapterId – Unique identifier in form of UUID representing a specific adapter.
// jobId – Unique identifier in form of UUID representing a specific job.
func (s *jobService) Get(adapterID string, jobID string) (job models.Job, err error) {
	targetUrl := s.url + s.basePath + "/" + adapterID + s.path + "/" + jobID
	resp, err := s.client.Get(targetUrl)
	if err != nil {
		return job, errors.Wrap(err, "Can't get job\n")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return job, errors.Wrap(err, "Can't convert job to byte slice\n")
	}

	err = handleHttpResponseCode(resp.StatusCode, body)
	if err != nil {
		return job, errors.Wrap(err, "Error in fetching job\n")
	}

	var jResource models.JobResource
	err = json.Unmarshal(body, &jResource)
	if err != nil {
		return job, errors.Wrap(err, "Can't unmarshal job response\n")
	}

	return jResource.ToJob()
}

// Send job with list of steps to adapter
// adapterId – Unique identifier in form of UUID representing a specific adapter.
func (s *jobService) Add(adapterID string, job models.NewJob) (event models.Job, err error) {
	reqBody, err := json.Marshal(job)
	if err != nil {
		return event, errors.Wrap(err, "Can't marshall req body for add job")
	}

	resp, err := s.client.Post(s.url+s.basePath+"/"+adapterID+s.path, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return event, errors.Wrap(err, "Can't post job\n")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return event, errors.Wrap(err, "Can't convert resp job to byte slice\n")
	}

	err = handleHttpResponseCode(resp.StatusCode, body)
	if err != nil {
		return event, errors.Wrap(err, "Error in post job\n")
	}

	var eRes models.JobResource
	err = json.Unmarshal(body, &eRes)
	if err != nil {
		return event, errors.Wrap(err, "Can't unmarshal post job response \n")
	}

	return eRes.ToJob()
}

// Delete all jobs from adapter
// adapterId – Unique identifier in form of UUID representing a specific adapter.
func (s *jobService) DeleteAll(adapterID string) (err error) {
	// Create request
	req, err := http.NewRequest("DELETE", s.url+s.basePath+"/"+adapterID+s.path, nil)
	if err != nil {
		return errors.Wrap(err, "Can't build delete all jobs request")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Can't delete all jobs\n")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Can't convert resp delete all jobs to byte slice\n")
	}

	err = handleHttpResponseCode(resp.StatusCode, body)
	if err != nil {
		return errors.Wrap(err, "Error in delete all jobs\n")
	}

	return err
}

// Delete job from adapter
// adapterId – Unique identifier in form of UUID representing a specific adapter.
// jobId – Unique identifier in form of UUID representing a specific job.
func (s *jobService) Delete(adapterID string, jobID string) (err error) {
	req, err := http.NewRequest("DELETE", s.url+s.basePath+"/"+adapterID+s.path+"/"+jobID, nil)
	if err != nil {
		return errors.Wrap(err, "Can't build delete job request")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Can't delete job\n")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Can't convert resp delete job to byte slice\n")
	}

	err = handleHttpResponseCode(resp.StatusCode, body)
	if err != nil {
		return errors.Wrap(err, "Error in delete job\n")
	}

	return err
}

// Update job status in adapter
// adapterId – Unique identifier in form of UUID representing a specific adapter.
// jobId – Unique identifier in form of UUID representing a specific job.
func (s *jobService) UpdateStatus(adapterID string, jobID string, status models.JobStatus) (job models.Job, err error) {
	reqBody, err := json.Marshal(models.JobStatusUpdate{Status: status.String()})
	if err != nil {
		return job, errors.Wrap(err, "Can't marshall req body for patch job status")
	}

	req, err := http.NewRequest("PATCH", s.url+s.basePath+"/"+adapterID+s.path+"/"+jobID, bytes.NewBuffer(reqBody))
	if err != nil {
		return job, errors.Wrap(err, "Can't build patch job status request")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return job, errors.Wrap(err, "Can't patch job status \n")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return job, errors.Wrap(err, "Can't convert resp patch job status  to byte slice\n")
	}

	err = handleHttpResponseCode(resp.StatusCode, body)
	if err != nil {
		return job, errors.Wrap(err, "Error in patch job status\n")
	}

	var eRes models.JobResource
	err = json.Unmarshal(body, &eRes)
	if err != nil {
		return job, errors.Wrap(err, "Can't unmarshal patch job status response \n")
	}

	return eRes.ToJob()
}

// Function builds jobs get query params
func buildJobsQueryParams(filter JobFilter) (queryParams string) {
	queryParams = ""

	if filter.Status != nil {
		queryParams += "&status=" + filter.Status.String()
	}

	if filter.SortBy != nil {
		queryParams += "&sortby=" + *filter.SortBy
	}

	if filter.SortDir != nil {
		queryParams += "&sortdir=" + *filter.SortDir
	}

	if filter.Offset != nil {
		queryParams += "&offset=" + strconv.Itoa(*filter.Offset)
	}

	if filter.Limit != nil {
		queryParams += "&limit=" + strconv.Itoa(*filter.Limit)
	}

	if len(queryParams) > 0 {
		// remove first & and add ?
		_, i := utf8.DecodeRuneInString(queryParams)
		return "?" + queryParams[i:]
	}

	return queryParams
}
