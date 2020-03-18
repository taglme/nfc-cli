package repository

import (
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	apiModels "github.com/taglme/nfc-goclient/pkg/models"
	"log"
	"os"
)

type PrinterService interface {
	Reset()
	PrintAppInfo(apiModels.AppInfo)
	PrintAdapters([]apiModels.Adapter)
	PrintJob(apiModels.Job)
	PrintNewJob(apiModels.NewJob)
	PrintJobSteps([]apiModels.JobStep)
	PrintTag(apiModels.Tag)
	PrintJobRun(apiModels.JobRun)
}

type printerService struct {
	writer table.Writer
}

func newPrinter() PrinterService {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)

	return &printerService{
		writer: t,
	}
}

func (s *printerService) Reset() {
	s.writer = table.NewWriter()
	s.writer.SetOutputMirror(os.Stdout)
	s.writer.SetStyle(table.StyleRounded)
}

func (s *printerService) PrintAppInfo(info apiModels.AppInfo) {
	fmt.Println("Server:")
	if len(info.Version) > 0 {
		fmt.Printf("   Version: %s\n", info.Version)
	}
	if len(info.Commit) > 0 {
		fmt.Printf("   Commit: %s\n", info.Commit)
	}
	if len(info.SDKInfo) > 0 {
		fmt.Printf("   SDK: %s\n", info.SDKInfo)
	}
	if len(info.Platform) > 0 {
		fmt.Printf("   Platform: %s\n", info.Platform)
	}
	if len(info.BuildTime) > 0 {
		fmt.Printf("   Build time: %s\n", info.BuildTime)
	}
}

func (s *printerService) PrintAdapters(adapters []apiModels.Adapter) {
	if len(adapters) == 0 {
		fmt.Println("Adapters not found")
		return
	}

	fmt.Println("Adapters:")

	for i, a := range adapters {
		fmt.Printf("[%d] %s\n", i+1, a.Name)
	}

	fmt.Println()
}

func (s *printerService) PrintNewJob(job apiModels.NewJob) {
	s.writer.AppendHeader(table.Row{
		"Job Name",
		"Repeat",
		"ExpireAfter",
	})

	s.writer.AppendRow(table.Row{
		job.JobName,
		job.Repeat,
		job.ExpireAfter,
	})

	s.writer.Render()
}

func (s *printerService) PrintJob(job apiModels.Job) {
	s.writer.AppendHeader(table.Row{
		"Job ID",
		"Job name",
		"Status",
		"Adapter id",
		"Adapter name",
		"Repeat",
		"Total runs",
		"Success runs",
		"Error runs",
		"Expire after",
		"Created at",
	})

	s.writer.AppendRow(table.Row{
		job.JobID,
		job.JobName,
		job.Status.String(),
		job.AdapterID,
		job.AdapterName,
		job.Repeat,
		job.TotalRuns,
		job.SuccessRuns,
		job.ErrorRuns,
		job.ExpireAfter,
		job.CreatedAt.String(),
	})

	s.writer.Render()
}

func (s *printerService) PrintJobSteps(js []apiModels.JobStep) {
	s.writer.AppendHeader(table.Row{
		"Command",
		"Command params",
	})

	for _, step := range js {
		p, err := json.MarshalIndent(step.Params.ToResource(), "", "  ")
		if err != nil {
			log.Printf("Can't marshall job step params: %s\n", err)
		}

		s.writer.AppendRow(table.Row{
			step.Command,
			string(p),
		})
	}

	s.writer.Render()
}

func (s *printerService) PrintJobRun(j apiModels.JobRun) {
	s.writer.AppendHeader(table.Row{
		"Run ID",
		"Job ID",
		"Job Name",
		"Status",
		"Adapter ID",
		"Adapter Name",
		"Created at",
	})

	s.writer.AppendRow(table.Row{
		j.RunID,
		j.JobID,
		j.JobName,
		j.Status.String(),
		j.AdapterID,
		j.AdapterName,
		j.CreatedAt.String(),
	})

	s.writer.Render()
}

func (s *printerService) PrintTag(t apiModels.Tag) {
	s.writer.AppendHeader(table.Row{
		"Tag ID",
		"Type",
		"Adapter ID",
		"Adapter Name",
		"Uid",
		"Atr",
		"Product",
		"Vendor",
	})

	s.writer.AppendRow(table.Row{
		t.TagID,
		t.Type.String(),
		t.AdapterID,
		t.AdapterName,
		t.Uid,
		t.Atr,
		t.Product,
		t.Vendor,
	})

	s.writer.Render()
}
