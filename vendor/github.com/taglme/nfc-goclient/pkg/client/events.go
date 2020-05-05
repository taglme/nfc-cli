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

//EventFilter defines filter for event list
type EventFilter struct {
	Name    *models.EventName
	SortBy  *string
	SortDir *string
	Offset  *int
	Limit   *int
}

//EventService event service interface
type EventService interface {
	GetAll() ([]models.Event, models.PageInfo, error)
	GetFiltered(adapterID *string, filter EventFilter) ([]models.Event, models.PageInfo, error)
	Add(event models.NewEvent) (models.Event, error)
}

type eventService struct {
	url    string
	path   string
	client *http.Client
}

func newEventService(c *http.Client, url string) EventService {
	return &eventService{
		url:    url,
		client: c,
		path:   "/events",
	}
}

// Endpoint returns information about events.
func (s *eventService) GetAll() ([]models.Event, models.PageInfo, error) {
	return s.GetFiltered(nil, EventFilter{})
}

// Endpoint returns information about events.
// adapterId – Adapter filter for events.
// filter.limit – Limit number of events in response.
// filter.offset – Offset from start of list.
// filter.sortby – Sort field for list.
// filter.sortdir – Sort direction for list
func (s *eventService) GetFiltered(adapterID *string, filter EventFilter) (events []models.Event, pagInfo models.PageInfo, err error) {
	queryParams := buildEventsQueryParams(adapterID, filter)
	resp, err := s.client.Get(s.url + s.path + queryParams)
	if err != nil {
		return events, pagInfo, errors.Wrap(err, "Can't get events")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events, pagInfo, errors.Wrap(err, "Can't convert events to byte slice")
	}

	err = handleHttpResponseCode(resp.StatusCode, body)
	if err != nil {
		return events, pagInfo, errors.Wrap(err, "Error in fetching events")
	}

	var eListResource models.EventListResource
	err = json.Unmarshal(body, &eListResource)
	if err != nil {
		return events, pagInfo, errors.Wrap(err, "Can't unmarshal events response")
	}

	events = make([]models.Event, len(eListResource.Items))
	for i, e := range eListResource.Items {
		events[i], err = e.ToEvent()
		if err != nil {
			return events, pagInfo, errors.Wrap(err, "Can't convert event resource to event model")
		}
	}

	return events, eListResource.GetPaginationInfo(), nil
}

// Send event to service
func (s *eventService) Add(ne models.NewEvent) (event models.Event, err error) {
	reqBody, err := json.Marshal(ne)
	if err != nil {
		return event, errors.Wrap(err, "Can't marshall req body for add event")
	}

	resp, err := s.client.Post(s.url+s.path, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return event, errors.Wrap(err, "Can't post event")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return event, errors.Wrap(err, "Can't convert resp event to byte slice")
	}

	err = handleHttpResponseCode(resp.StatusCode, body)
	if err != nil {
		return event, errors.Wrap(err, "Error in post event")
	}

	var eRes models.EventResource
	err = json.Unmarshal(body, &eRes)
	if err != nil {
		return event, errors.Wrap(err, "Can't unmarshal post event response")
	}

	return eRes.ToEvent()
}

// Function builds events get query params
func buildEventsQueryParams(adapterID *string, filter EventFilter) (queryParams string) {
	queryParams = ""

	if adapterID != nil {
		queryParams += "&adapter_id=" + *adapterID
	}

	if filter.Name != nil {
		queryParams += "&name=" + filter.Name.String()
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
