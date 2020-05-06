package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	uuid "github.com/nu7hatch/gouuid"
)

type NewEvent struct {
	Name      string      `json:"name" binding:"required"`
	AdapterID string      `json:"adapter_id"`
	Data      interface{} `json:"data"`
}

type Event struct {
	EventID     string
	Name        EventName
	AdapterID   string
	AdapterName string
	Data        interface{}
	CreatedAt   time.Time
}

type EventResource struct {
	EventID     string      `json:"event_id"`
	Name        string      `json:"name"`
	AdapterID   string      `json:"adapter_id"`
	AdapterName string      `json:"adapter_name"`
	Data        interface{} `json:"data"`
	CreatedAt   string      `json:"created_at"`
}

type EventListResource struct {
	Total  int             `json:"total"`
	Length int             `json:"length"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
	Items  []EventResource `json:"items"`
}

func (e EventResource) ToEvent() (event Event, err error) {
	eName, ok := StringToEventName(e.Name)
	if !ok {
		return Event{}, errors.New("Can't convert event resource name to event name")
	}

	t, err := time.Parse(time.RFC3339, e.CreatedAt)
	if err != nil {
		return Event{}, errors.Wrap(err, "Can't parse event resource created at")
	}

	return Event{
		EventID:     e.EventID,
		Name:        eName,
		AdapterID:   e.AdapterID,
		AdapterName: e.AdapterName,
		Data:        e.Data,
		CreatedAt:   t,
	}, nil
}

func (list EventListResource) GetPaginationInfo() PageInfo {
	return PageInfo{
		Total:  list.Total,
		Length: list.Length,
		Limit:  list.Limit,
		Offset: list.Offset,
	}
}

func (e Event) ToResource() EventResource {
	resource := EventResource{
		EventID:     e.EventID,
		Name:        e.Name.String(),
		AdapterID:   e.AdapterID,
		AdapterName: e.AdapterName,
		Data:        e.Data,
		CreatedAt:   e.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
	}
	return resource
}

func (ne NewEvent) ToEvent(adapterName string) Event {
	id, _ := uuid.NewV4()
	createdAt := time.Now().UTC()
	name, _ := StringToEventName(ne.Name)
	e := Event{
		EventID:     id.String(),
		Name:        name,
		AdapterID:   ne.AdapterID,
		AdapterName: adapterName,
		Data:        ne.Data,
		CreatedAt:   createdAt,
	}

	return e
}

type EventName int

const (
	EventNameTagDiscovery EventName = iota + 1
	EventNameTagRelease
	EventNameAdapterDiscovery
	EventNameAdapterRelease
	EventNameJobSubmited
	EventNameJobActivated
	EventNameJobPended
	EventNameJobDeleted
	EventNameJobFinished
	EventNameRunStarted
	EventNameRunSuccess
	EventNameRunError
	EventNameServerStarted
	EventNameServerStopped
)

func StringToEventName(s string) (EventName, bool) {
	switch s {
	case EventNameTagDiscovery.String():
		return EventNameTagDiscovery, true
	case EventNameTagRelease.String():
		return EventNameTagRelease, true
	case EventNameAdapterDiscovery.String():
		return EventNameAdapterDiscovery, true
	case EventNameAdapterRelease.String():
		return EventNameAdapterRelease, true
	case EventNameJobSubmited.String():
		return EventNameJobSubmited, true
	case EventNameJobActivated.String():
		return EventNameJobActivated, true
	case EventNameJobPended.String():
		return EventNameJobPended, true
	case EventNameJobDeleted.String():
		return EventNameJobDeleted, true
	case EventNameJobFinished.String():
		return EventNameJobFinished, true
	case EventNameRunStarted.String():
		return EventNameRunStarted, true
	case EventNameRunSuccess.String():
		return EventNameRunSuccess, true
	case EventNameRunError.String():
		return EventNameRunError, true
	case EventNameServerStarted.String():
		return EventNameServerStarted, true
	case EventNameServerStopped.String():
		return EventNameServerStopped, true
	}
	return 0, false
}

func (eventName EventName) String() string {
	names := [...]string{
		"unknown",
		"tag_discovery",
		"tag_release",
		"adapter_discovery",
		"adapter_release",
		"job_submited",
		"job_activated",
		"job_pended",
		"job_deleted",
		"job_finished",
		"run_started",
		"run_success",
		"run_error",
		"server_started",
		"server_stopped",
	}

	if eventName < EventNameTagDiscovery || eventName > EventNameServerStopped {
		return names[0]
	}
	return names[eventName]
}

func (e Event) GetAdapter() (adapter AdapterResource, ok bool) {
	jsonBytes, err := json.Marshal(e.Data)
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonBytes, &adapter)
	if err != nil {
		return
	}
	ok = true
	return
}
func (e Event) GetTag() (tag TagResource, ok bool) {
	jsonBytes, err := json.Marshal(e.Data)
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonBytes, &tag)
	if err != nil {
		return
	}
	ok = true
	return
}
func (e Event) GetJob() (job JobResource, ok bool) {
	jsonBytes, err := json.Marshal(e.Data)
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonBytes, &job)
	if err != nil {
		return
	}
	ok = true
	return
}
func (e Event) GetRun() (run JobRunResource, ok bool) {
	jsonBytes, err := json.Marshal(e.Data)
	if err != nil {
		fmt.Printf("Error marshal run resource %s", err)
		return
	}
	err = json.Unmarshal(jsonBytes, &run)
	if err != nil {
		//Todo fix unmarshal error:
		//Error unmarshal run resource json: cannot unmarshal object into Go struct field StepResultResource.results.params of type models.CommandParamsResourc
		fmt.Printf("Error unmarshal run resource %s", err)
		return
	}
	ok = true
	return
}
