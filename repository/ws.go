package repository

import (
	"fmt"
	"github.com/taglme/nfc-cli/models"
	apiModels "github.com/taglme/nfc-client/pkg/models"
	"time"
)

func (s *ApiService) RunWsConnection(handler func(models.Event)) error {
	s.client.Ws.OnEvent(func(event apiModels.Event) {
		eventHandler(event)
		handler(MapApiEventNameToCliEvent[event.Name])
	})

	err := s.client.Ws.Connect()
	if err != nil || !s.client.Ws.IsConnected() {
		return err
	}

	return nil
}

func (s *ApiService) StopWsConnection() error {
	return s.client.Ws.Disconnect()
}

func eventHandler(e apiModels.Event) {
	switch e.Name {
	//case models.EventJobSubmitted:
	//	j, ok := e.Data.(*apiModels.Job)
	//	fmt.Println(j, ok)
	//	fmt.Printf("Job has been submitted: %s\n", e.Data)
	case apiModels.EventNameRunStarted:
		fmt.Println(parseJobRunStruct(e.Data))
		fmt.Println("\nJob Run has been started:")
	//case models.EventRunError:
	//	fmt.Printf("Job Run has been finished with error: %s\n", e.Data)
	//case models.EventRunSuccess:
	//	fmt.Printf("Job Run has been finished with success: %s\n", e.Data)
	case apiModels.EventNameJobFinished:
		fmt.Println("Job has been finished with success:")
	case apiModels.EventNameJobDeleted:
		fmt.Println("Job has been deleted:")
	}
}

func parseJobRunStruct(data interface{}) (jr apiModels.JobRun) {
	m := data.(map[string]interface{})

	if runId, ok := m["run_id"].(string); ok {
		jr.RunID = runId
	}

	if jobId, ok := m["job_id"].(string); ok {
		jr.JobID = jobId
	}

	if jobName, ok := m["job_name"].(string); ok {
		jr.JobName = jobName
	}

	if status, ok := m["status"].(string); ok {
		jr.Status, ok = apiModels.StringToJobRunStatus(status)
		if !ok {
			fmt.Println("Can't parse Job run status")
		}
	}

	if adapterId, ok := m["adapter_id"].(string); ok {
		jr.AdapterID = adapterId
	}

	if adapterName, ok := m["adapter_name"].(string); ok {
		jr.AdapterName = adapterName
	}

	if tag, ok := m["tag"].(apiModels.TagResource); ok {
		jr.Tag, _ = tag.ToTag()
	}

	if createdAt, ok := m["created_at"].(string); ok {
		t, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			fmt.Println("Error on converting created at string to time")
		}

		jr.CreatedAt = t
	}

	if results, ok := m["results"].([]struct{}); ok {
		for _, r := range results {
			srRes := data.(apiModels.StepResultResource)
			fmt.Println(srRes, r)
		}
		//jr.Tag, _ = tag.ToTag()
	}

	return jr
}
