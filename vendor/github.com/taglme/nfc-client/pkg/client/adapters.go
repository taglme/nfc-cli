package client

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"

	"github.com/taglme/nfc-client/pkg/models"
)

type AdapterFilter struct {
	AdapterType models.AdapterType
}
type AdapterService interface {
	GetAll() ([]models.Adapter, error)
	GetFiltered(adapterType *models.AdapterType) ([]models.Adapter, error)
	Get(adapterID string) (models.Adapter, error)
}

type adapterService struct {
	url    string
	client *http.Client
	path   string
}

func newAdapterService(client *http.Client, url string) AdapterService {
	return &adapterService{
		url:    url,
		client: client,
		path:   "/adapters",
	}
}

// Adapters list endpoint returns information about all adapters. The response includes array of Adapters
func (s *adapterService) GetAll() (adapters []models.Adapter, err error) {
	return s.GetFiltered(nil)
}

// Adapters list endpoint returns information about all adapters. The response includes array of Adapters
// adapterType – Adapters' type filter.
func (s *adapterService) GetFiltered(adapterType *models.AdapterType) (adapters []models.Adapter, err error) {
	queryParams := ""
	if adapterType != nil {
		queryParams += "?type=" + adapterType.String()
	}

	resp, err := s.client.Get(s.url + s.path + queryParams)
	if err != nil {
		return adapters, errors.Wrap(err, "Can't get adapters\n")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return adapters, errors.Wrap(err, "Can't convert adapters to byte slice\n")
	}

	err = handleHttpResponseCode(resp.StatusCode, body)
	if err != nil {
		return adapters, errors.Wrap(err, "Error in fetching adapters\n")
	}

	var aListResource models.AdapterListResource
	err = json.Unmarshal(body, &aListResource)
	if err != nil {
		return adapters, errors.Wrap(err, "Can't unmarshal events response\n")
	}

	adapters = make([]models.Adapter, len(aListResource))
	for i, a := range aListResource {
		adapters[i] = a.ToAdapter()
	}

	return adapters, nil
}

// Get adapter with all details
// adapterID – Unique identifier in form of UUID representing a specific adapter.
func (s *adapterService) Get(adapterID string) (adapter models.Adapter, err error) {
	resp, err := s.client.Get(s.url + s.path + "/" + adapterID)
	if err != nil {
		return adapter, errors.Wrap(err, "Can't get adapter\n")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return adapter, errors.Wrap(err, "Can't convert adapter to byte slice\n")
	}

	err = handleHttpResponseCode(resp.StatusCode, body)
	if err != nil {
		return adapter, errors.Wrap(err, "Error in fetching adapters\n")
	}

	var aResource models.AdapterResource
	err = json.Unmarshal(body, &aResource)
	if err != nil {
		return adapter, errors.Wrap(err, "Can't unmarshal adapter response\n")
	}

	return aResource.ToAdapter(), nil
}
