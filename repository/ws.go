package repository

import (
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"

	"github.com/taglme/nfc-cli/models"
	apiModels "github.com/taglme/nfc-goclient/pkg/models"
)

func (s *RepositoryService) RunWsConnection(eHandler func(models.Event, interface{}), errHandler func(error)) error {
	s.client.Ws.OnEvent(func(event apiModels.Event) {
		s.eventHandler(event)
		eHandler(MapApiEventNameToCliEvent[event.Name], event.Data)
	})

	s.client.Ws.OnError(func(err error) {
		errHandler(err)
	})

	err := s.client.Ws.Connect()
	if err != nil || !s.client.Ws.IsConnected() {
		return err
	}

	return nil
}

func (s *RepositoryService) StopWsConnection() error {
	if s.client.Ws.IsConnected() {
		return s.client.Ws.Disconnect()
	}

	return nil
}

func (s *RepositoryService) eventHandler(e apiModels.Event) {
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
	case apiModels.EventNameRunSuccess, apiModels.EventNameRunError:
		j, ok := e.GetJob()
		if !ok {
			log.Println("Can't get Job from Event.")
			return
		}

		if e.Name.String() == "run_success" {
			color.Green("Job %s: run finished successfully.\n", j.JobName)
		} else {
			color.Red("Job %s: run finished unsuccessfully.\n", j.JobName)
		}

		jobRun := parseJobRunStruct(e.Data)
		fmt.Printf("Job %s: -----run results start-----\n", j.JobName)

		for i, s := range jobRun.Results {
			endStr := ""
			if len(s.Message) > 0 {
				endStr = "(" + s.Message + ")"
			}

			if s.Status == apiModels.CommandStatusSuccess {
				color.Green("[Step %d] %s – %s %s", i+1, MapRunStepCmdToString[s.Command], s.Status.String(), endStr)
			} else {
				color.Red("[Step %d] %s – %s %s", i+1, MapRunStepCmdToString[s.Command], s.Status.String(), endStr)
			}

			if s.Params != nil {
				pStr := s.Params.String()
				if len(pStr) > 0 {
					fmt.Printf("Params:\n%s\n", pStr)
				}
			}

			if s.Output != nil {
				oStr := s.Output.String()
				if len(oStr) > 0 {
					if strings.Contains(oStr, "Record") {
						if strings.Contains(oStr, "Empty") {
							oStr = color.CyanString(oStr)
						} else {
							oStr = color.MagentaString(oStr)
						}
					}
					fmt.Printf("Output:\n%s\n", oStr)

				}
			}
		}

		fmt.Printf("Job %s: -----run results end-----\n", j.JobName)

		job, err := s.GetJob(j.AdapterID, j.JobID)
		if err == nil {
			// we are not handling this error as job simply can be deleted at this point so request will always fail at last iteration
			fmt.Printf("Job %s: total %d runs (%d success, %d failed). Remain %d runs\n", job.JobName, job.TotalRuns, job.SuccessRuns, job.ErrorRuns, job.Repeat-job.SuccessRuns)

			if job.Repeat-job.SuccessRuns > 0 {
				fmt.Printf("Job %s: waiting for NFC tag...\n", j.JobName)
			}
		}
	case apiModels.EventNameJobFinished:
		job, ok := e.GetJob()
		if !ok {
			log.Println("Can't get Job from Event.")
			return
		}

		fmt.Printf("Job %s: total %d runs (%d success, %d failed). Remain %d runs\n", job.JobName, job.TotalRuns, job.SuccessRuns, job.ErrorRuns, job.Repeat-job.SuccessRuns)
		fmt.Printf("Job %s: finished successfully by adapter %s\n", job.JobName, job.AdapterName)
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
