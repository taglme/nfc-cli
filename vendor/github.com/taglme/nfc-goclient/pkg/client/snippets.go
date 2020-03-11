package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"unicode/utf8"

	"github.com/pkg/errors"
	"github.com/taglme/nfc-goclient/pkg/models"
)

//SnippetFilter defines filter for snippet list
type SnippetFilter struct {
	Category *models.SnippetCategory
	UsageID  *string
}

//SnippetService snippet service interface
type SnippetService interface {
	GetAll() ([]models.Snippet, error)
	GetFiltered(filter SnippetFilter) ([]models.Snippet, error)
}

type snippetService struct {
	url    string
	path   string
	client *http.Client
}

func newSnippetService(c *http.Client, url string) SnippetService {
	return &snippetService{
		url:    url,
		client: c,
		path:   "/snippets",
	}
}

// Snippets list endpoint returns information about all snippets. The response includes array of Snippets
func (s *snippetService) GetAll() ([]models.Snippet, error) {
	return s.GetFiltered(SnippetFilter{})
}

// Snippets list endpoint returns information about all snippets. The response includes array of Snippets
// category – category filter for snippet.
// usage_id – usage_id filter for snippet
func (s *snippetService) GetFiltered(filter SnippetFilter) (snippets []models.Snippet, err error) {
	queryParams := buildSnippetsQueryParams(filter)
	resp, err := s.client.Get(s.url + s.path + queryParams)
	if err != nil {
		return snippets, errors.Wrap(err, "Can't get snippets\n")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return snippets, errors.Wrap(err, "Can't convert snippets to byte slice\n")
	}

	err = handleHttpResponseCode(resp.StatusCode, body)
	if err != nil {
		return snippets, errors.Wrap(err, "Error in fetching snippets\n")
	}

	var sListResource models.SnippetListResource
	err = json.Unmarshal(body, &sListResource)
	if err != nil {
		return snippets, errors.Wrap(err, "Can't unmarshal snippets response\n")
	}

	snippets = make([]models.Snippet, len(sListResource))
	for i, s := range sListResource {
		snippets[i], err = s.ToSnippet()
		if err != nil {
			return snippets, errors.Wrap(err, "Can't convert snippet resource\n")
		}
	}

	return snippets, nil
}

// Function builds events get query params
func buildSnippetsQueryParams(filter SnippetFilter) (queryParams string) {
	queryParams = ""

	if filter.UsageID != nil {
		queryParams += "&usage_id=" + *filter.UsageID
	}

	if filter.Category != nil {
		queryParams += "&category=" + filter.Category.String()
	}

	if len(queryParams) > 0 {
		// remove first & and add ?
		_, i := utf8.DecodeRuneInString(queryParams)
		return "?" + queryParams[i:]
	}

	return queryParams
}
