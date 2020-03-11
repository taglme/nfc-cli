package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/taglme/nfc-goclient/pkg/models"
)

//AboutService about service interface
type AboutService interface {
	Get() (models.AppInfo, error)
}

type aboutService struct {
	url    string
	client *http.Client
}

func newAboutService(c *http.Client, url string) AboutService {
	return &aboutService{
		url:    url,
		client: c,
	}
}

// Get fetching the about info
func (s *aboutService) Get() (info models.AppInfo, err error) {
	targetURL := s.url + "/about"

	resp, err := s.client.Get(targetURL)
	if err != nil {
		return info, errors.Wrap(err, "Can't get about info")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return info, errors.Wrap(err, "Can't convert about info to byte slice")
	}

	err = handleHttpResponseCode(resp.StatusCode, body)
	if err != nil {
		return info, errors.Wrap(err, "Error in fetching about info")
	}

	err = json.Unmarshal(body, &info)
	if err != nil {
		return info, errors.Wrap(err, "Can't unmarshal about info response response")
	}

	return info, nil
}
