package repository

import (
	"fmt"
	"github.com/taglme/nfc-cli/models"
	apiModels "github.com/taglme/nfc-goclient/pkg/models"
	"log"
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

//Остальные команды (без флага exported)
//Job {job.Name}: submitted to adapter {adapter.Name} (job_submitted event)
//Job {job.Name}: activated. Waiting for NFC tag... (job_activated event)
//Job {job.Name}: execution started. Hold NFC tag steady... (run_started event) Job {job.Name}: run finished successfully.
//Total {total} runs ({success} success, {failed} failed). Remain {repeat- success} runs. Run results:
//[1] {stepResult.Command} – {stepResult.Status} ({stepResult.Message})
//Params: {stepResult.Params}
//Output: {stepResult.Output}
//[2] {stepResult.Command} – { stepResult.Status} ({stepResult.Message})
//Params: { stepResult.Params}
//Output: { stepResult.Output}
//(run_success or run_error event)
//Job {job.Name}: finished successfully by adapter {adapter.Name} (job_finished)

type CommandParams interface {
	apiModels.CommandParams
	Print(string)
}

type TransmitAdapterParams apiModels.TransmitAdapterParams

func (m TransmitAdapterParams) Print() string {
	return string(m.TxBytes)
}

func (s *ApiService) eventHandler(e apiModels.Event) {
	switch e.Name {
	case apiModels.EventNameJobSubmited:
		j, ok := e.GetJob()
		if !ok {
			log.Println("Can't get Job from Event.")
			return
		}
		fmt.Printf("Job %s: submitted to adapter %s\n", j.JobName, j.AdapterName)
	case apiModels.EventNameJobActivated:
		j, ok := e.GetJob()
		if !ok {
			log.Println("Can't get Job from Event.")
			return
		}
		fmt.Printf("Job %s: activated. Waiting for NFC tag...\n", j.JobName)
	case apiModels.EventNameRunStarted:
		j, ok := e.GetJob()
		if !ok {
			log.Println("Can't get Job from Event.")
			return
		}
		fmt.Printf("Job %s: execution started. Hold NFC tag steady...\n", j.JobName)
	case apiModels.EventNameRunError:
	case apiModels.EventNameRunSuccess:
		j, ok := e.GetJob()
		if !ok {
			log.Println("Can't get Job from Event.")
			return
		}
		if e.Name.String() == "run_success" {
			fmt.Printf("Job %s: run finished successfully.\n", j.JobName)
		} else {
			fmt.Printf("Job %s: run finished unsuccessfully.\n", j.JobName)
		}

		job, err := s.GetJob(j.AdapterID, j.JobID)
		if err == nil {
			// we are not handling this error as job simply can be deleted at this point so request will always fail at last iteration
			fmt.Printf("Total %d runs (%d success, %d failed). Remain %d runs\n", job.TotalRuns, job.SuccessRuns, job.ErrorRuns, job.Repeat-job.SuccessRuns)
		}

		jobRun := parseJobRunStruct(e.Data)
		fmt.Println("Run results:")
		//s.printer.PrintStepResults(jobS)
		for i, s := range jobRun.Results {
			fmt.Printf("[%d] %s – %s (%s)\n", i+1, MapRunStepCmdToString[s.Command], s.Status.String(), s.Message)
			fmt.Printf("Params: %s\n", s.Params.String())
			fmt.Printf("Output: %s\n", s.Output.String())
		}

		if job.Repeat-job.SuccessRuns > 0 {
			fmt.Printf("Job %s: waiting for NFC tag...\n", j.JobName)
		}
	case apiModels.EventNameJobFinished:
		j, ok := e.GetJob()
		if !ok {
			log.Println("Can't get Job from Event.")
			return
		}

		fmt.Printf("Job %s: finished successfully by adapter %s\n", j.JobName, j.AdapterName)
	case apiModels.EventNameJobDeleted:
		fmt.Println("Job has been deleted")
	}
}

var MapRunStepCmdToString = map[apiModels.Command]string{
	apiModels.CommandGetTags:         "Get tags",
	apiModels.CommandTransmitAdapter: "Transmit adapter",
	apiModels.CommandTransmitTag:     "Transmit tag",
	apiModels.CommandWriteNdef:       "Write NDEF",
	apiModels.CommandReadNdef:        "Read NDEF",
	apiModels.CommandFormatDefault:   "Format default",
	apiModels.CommandLockPermanent:   "Lock permanent",
	apiModels.CommandSetPassword:     "Set password",
	apiModels.CommandRemovePassword:  "Remove password",
	apiModels.CommandAuthPassword:    "Authorize password",
	apiModels.CommandGetDump:         "Get dump",
}
