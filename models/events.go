package models

type Event = string

const (
	EventTagDiscovery     Event = "tag_discovery"
	EventTagRelease       Event = "tag_release"
	EventAdapterDiscovery Event = "adapter_discovery"
	EventAdapterRelease   Event = "adapter_release"
	EventJobSubmitted     Event = "job_submited"
	EventJobActivated     Event = "job_activated"
	EventJobPended        Event = "job_pended"
	EventJobDeleted       Event = "job_deleted"
	EventJobFinished      Event = "job_finished"
	EventRunStarted       Event = "run_started"
	EventRunSuccess       Event = "run_success"
	EventRunError         Event = "run_error"
)
