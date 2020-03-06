package client

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-client/pkg/models"
	"io/ioutil"
	"net/http"
)

type TagService interface {
	GetAll(adapterID string, tagType *models.TagType) ([]models.Tag, error)
	Get(adapterID string, tagID string) (models.Tag, error)
}

type tagService struct {
	url      string
	basePath string
	path     string
	client   *http.Client
}

func newTagService(c *http.Client, url string) TagService {
	return &tagService{
		url:      url,
		client:   c,
		path:     "/tags",
		basePath: "/adapters",
	}
}

// Get all adapter's tags
// adapterID – Unique identifier in form of UUID representing a specific adapter.
// tagType – Tags' type filter.
func (s *tagService) GetAll(adapterID string, tagType *models.TagType) (tags []models.Tag, err error) {
	targetUrl := s.url + s.basePath + "/" + adapterID + s.path
	if tagType != nil {
		targetUrl += "?type=" + tagType.String()
	}

	resp, err := s.client.Get(targetUrl)
	if err != nil {
		return tags, errors.Wrap(err, "Can't get tags\n")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tags, errors.Wrap(err, "Can't convert tags to byte slice\n")
	}

	err = handleHttpResponseCode(resp.StatusCode, body)
	if err != nil {
		return tags, errors.Wrap(err, "Error in fetching tags\n")
	}

	var tListResource models.TagListResource
	err = json.Unmarshal(body, &tListResource)
	if err != nil {
		return tags, errors.Wrap(err, "Can't unmarshal tags response\n")
	}

	tags = make([]models.Tag, len(tListResource))
	for i, t := range tListResource {
		tags[i], err = t.ToTag()
		if err != nil {
			return tags, errors.Wrap(err, "Can't convert tag resource to tag model\n")
		}
	}

	return tags, nil
}

// Get all specified tag's details in adapter
// adapterID – Unique identifier in form of UUID representing a specific adapter.
// tagID – Unique identifier in form of UUID representing a specific tag.
func (s *tagService) Get(adapterID string, tagID string) (tag models.Tag, err error) {
	targetUrl := s.url + s.basePath + "/" + adapterID + s.path + "/" + tagID

	resp, err := s.client.Get(targetUrl)
	if err != nil {
		return tag, errors.Wrap(err, "Can't get tag\n")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tag, errors.Wrap(err, "Can't convert tag to byte slice\n")
	}

	err = handleHttpResponseCode(resp.StatusCode, body)
	if err != nil {
		return tag, errors.Wrap(err, "Error in fetching tag\n")
	}

	var tRes models.TagResource
	err = json.Unmarshal(body, &tRes)
	if err != nil {
		return tag, errors.Wrap(err, "Can't unmarshal tag response\n")
	}

	return tRes.ToTag()
}
