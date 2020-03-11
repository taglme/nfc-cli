package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/pkg/errors"
	"github.com/taglme/nfc-goclient/pkg/models"
)

//RunFilter defines filter for run list
type RunFilter struct {
	JobID   *string
	Status  *models.JobRunStatus
	SortBy  *string
	SortDir *string
	Limit   *int
	Offset  *int
}

//RunService run service interface
type RunService interface {
	GetAll(adapterID string) ([]models.JobRun, models.PageInfo, error)
	GetFiltered(adapterID string, filter RunFilter) ([]models.JobRun, models.PageInfo, error)
	Get(adapterID string, runID string) (models.JobRun, error)
}

type runService struct {
	url      string
	basePath string
	path     string
	client   *http.Client
}

func newRunService(c *http.Client, url string) RunService {
	return &runService{
		url:      url,
		client:   c,
		path:     "/runs",
		basePath: "/adapters",
	}
}

// Get Run list for adapter with all details
func (s *runService) GetAll(adapterID string) ([]models.JobRun, models.PageInfo, error) {
	return s.GetFiltered(adapterID, RunFilter{})
}

// Get Run list for adapter with all details
// adapterId – Unique identifier in form of UUID representing a specific adapter.
// filter.status – Runs' status filter.
// filter.job_id – Filter Run by specified Job
// filter.limit – Limit number of events in response.
// filter.offset – Offset from start of list.
// filter.sortBy – Sort field for list.
// filter.sortDir – Sort direction for list
func (s *runService) GetFiltered(adapterID string, filter RunFilter) (runs []models.JobRun, pagInfo models.PageInfo, err error) {
	targetUrl := s.url + s.basePath + "/" + adapterID + s.path + buildRunsQueryParams(filter)
	resp, err := s.client.Get(targetUrl)
	if err != nil {
		return runs, pagInfo, errors.Wrap(err, "Can't get runs\n")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return runs, pagInfo, errors.Wrap(err, "Can't convert runs to byte slice\n")
	}

	err = handleHttpResponseCode(resp.StatusCode, body)
	if err != nil {
		return runs, pagInfo, errors.Wrap(err, "Error in fetching runs\n")
	}

	var rListResource models.JobRunListResource
	err = json.Unmarshal(body, &rListResource)
	if err != nil {
		return runs, pagInfo, errors.Wrap(err, "Can't unmarshal runs response\n")
	}

	runs = make([]models.JobRun, len(rListResource.Items))
	for i, e := range rListResource.Items {
		runs[i], err = e.ToJobRun()
		if err != nil {
			return runs, pagInfo, errors.Wrap(err, "Can't convert job run resource to job run model\n")
		}
	}

	return runs, rListResource.GetPaginationInfo(), nil
}

// Get all specefied jobrun's details
// runId – Unique identifier in form of UUID representing a specific job run.
// adapterId – Unique identifier in form of UUID representing a specific adapter.
func (s *runService) Get(adapterID string, runID string) (run models.JobRun, err error) {
	targetUrl := s.url + s.basePath + "/" + adapterID + s.path + "/" + runID
	resp, err := s.client.Get(targetUrl)
	if err != nil {
		return run, errors.Wrap(err, "Can't get run\n")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return run, errors.Wrap(err, "Can't convert run to byte slice\n")
	}

	err = handleHttpResponseCode(resp.StatusCode, body)
	if err != nil {
		return run, errors.Wrap(err, "Error in fetching run\n")
	}

	var rResource models.JobRunResource
	err = json.Unmarshal(body, &rResource)
	if err != nil {
		return run, errors.Wrap(err, "Can't unmarshal run response\n")
	}

	return rResource.ToJobRun()
}

// Function builds runs get query params
func buildRunsQueryParams(filter RunFilter) (queryParams string) {
	queryParams = ""

	if filter.Status != nil {
		queryParams += "&status=" + filter.Status.String()
	}

	if filter.JobID != nil {
		queryParams += "&job_id=" + *filter.JobID
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
