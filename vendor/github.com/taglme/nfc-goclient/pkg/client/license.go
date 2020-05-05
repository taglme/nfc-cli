package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/taglme/nfc-goclient/pkg/models"
)

//license service interface
type LicenseService interface {
	GetLicense() (models.License, error)
	GetAppLicense(appID string) (models.AppLicense, error)
}

type licenseService struct {
	url      string
	basePath string
	path     string
	client   *http.Client
}

func newLicenseService(c *http.Client, url string) LicenseService {
	return &licenseService{
		url:      url,
		client:   c,
		path:     "/apps",
		basePath: "/licenses",
	}
}

// Get license with all details
func (s *licenseService) GetLicense() (license models.License, err error) {
	resp, err := s.client.Get(s.url + s.basePath)
	if err != nil {
		return license, errors.Wrap(err, "Can't get license\n")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return license, errors.Wrap(err, "Can't read response body\n")
	}
	err = handleHttpResponseCode(resp.StatusCode, body)
	if err != nil {
		return license, errors.Wrap(err, "Error in fetching license\n")
	}
	var licenseResource models.LicenseResource
	err = json.Unmarshal(body, &licenseResource)
	if err != nil {
		return license, errors.Wrap(err, "Can't unmarshal license response\n")
	}
	lic, err := licenseResource.ToLicense()
	if err != nil {
		return license, errors.Wrap(err, "Can't convert license resource\n")
	}
	return lic, nil
}

// Get license for specific application with all details
func (s *licenseService) GetAppLicense(appID string) (appLicense models.AppLicense, err error) {
	resp, err := s.client.Get(s.url + s.basePath + s.path + "/" + appID)
	if err != nil {
		return appLicense, errors.Wrap(err, "Can't get app license\n")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return appLicense, errors.Wrap(err, "Can't read response body\n")
	}
	err = handleHttpResponseCode(resp.StatusCode, body)
	if err != nil {
		return appLicense, errors.Wrap(err, "Error in fetching application license\n")
	}
	var appLicenseResource models.AppLicenseResource
	err = json.Unmarshal(body, &appLicenseResource)
	if err != nil {
		return appLicense, errors.Wrap(err, "Can't unmarshal application license response\n")
	}
	appLic, err := appLicenseResource.ToAppLicense()
	if err != nil {
		return appLicense, errors.Wrap(err, "Can't convert application license resource\n")
	}
	return appLic, nil
}
