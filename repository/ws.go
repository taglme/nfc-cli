package repository

import (
	"fmt"
	"github.com/taglme/nfc-cli/models"
	apiModels "github.com/taglme/nfc-client/pkg/models"
)

func (s *ApiService) RunWsConnection(handler func(models.Event, interface{})) error {
	s.client.Ws.OnEvent(func(event apiModels.Event) {
		s.eventHandler(event)
		handler(MapApiEventNameToCliEvent[event.Name], event.Data)
	})

	err := s.client.Ws.Connect()
	if err != nil || !s.client.Ws.IsConnected() {
		return err
	}

	return nil
}

func (s *ApiService) StopWsConnection() error {
	if s.client.Ws.IsConnected() {
		return s.client.Ws.Disconnect()
	}

	return nil
}

func (s *ApiService) eventHandler(e apiModels.Event) {
	switch e.Name {
	case apiModels.EventNameJobSubmited:
		fmt.Println("Job has been submitted. Please move the nfc tag near the device.")
	case apiModels.EventNameRunStarted:
		fmt.Println("Job Run has been started. Please keep the nfc tag near the device.")
	case apiModels.EventNameRunError:
		jobRun := parseJobRunStruct(e.Data)
		fmt.Println("Job Run has been finished with error:")
		s.printer.PrintJobRun(jobRun)
		s.printer.Reset()
		fmt.Println("Job Run related tag:")
		s.printer.PrintTag(jobRun.Tag)
		s.printer.Reset()
	case apiModels.EventNameRunSuccess:
		jobRun := parseJobRunStruct(e.Data)
		fmt.Println("Job Run has been finished with success")
		s.printer.PrintJobRun(jobRun)
		s.printer.Reset()
		fmt.Println("Job Run related tag:")
		s.printer.PrintTag(jobRun.Tag)
		s.printer.Reset()
		fmt.Println("Job Run step results:")
		s.printer.PrintStepResults(jobRun.Results)
		s.printer.Reset()
	case apiModels.EventNameJobFinished:
		fmt.Println("Job has been finished with success")
	case apiModels.EventNameJobDeleted:
		fmt.Println("Job has been deleted")
	}
}
