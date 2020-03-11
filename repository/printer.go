package repository

import (
	"github.com/jedib0t/go-pretty/table"
	apiModels "github.com/taglme/nfc-client/pkg/models"
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
	PrintStepResults([]apiModels.StepResult)
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
	s.writer.AppendHeader(table.Row{
		"Name",
		"Version",
		"Commit",
		"SDK Info",
		"Platform",
		"Build time",
		"CheckSuccess",
		"Supported",
		"Have update",
		"Update version",
		"Update download",
		"Started at",
	})
	s.writer.AppendRow(table.Row{
		info.Name,
		info.Version,
		info.Commit,
		info.SDKInfo,
		info.Platform,
		info.BuildTime,
		info.CheckSuccess,
		info.Supported,
		info.HaveUpdate,
		info.UpdateVersion,
		info.UpdateDownload,
		info.StartedAt,
	})
	s.writer.Render()
}

func (s *printerService) PrintAdapters(adapters []apiModels.Adapter) {
	s.writer.AppendHeader(table.Row{"Adapter ID", "Name", "Type", "Driver"})

	for _, a := range adapters {
		s.writer.AppendRow(table.Row{a.AdapterID, a.Name, a.Type.String(), a.Driver})
	}
	s.writer.SetStyle(table.StyleLight)
	s.writer.Render()
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
		s.writer.AppendRow(table.Row{
			step.Command,
			step.Params,
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

func (s *printerService) PrintStepResults(r []apiModels.StepResult) {
	s.writer.AppendHeader(table.Row{
		"Command",
		"Params",
		"Output",
		"Status",
		"Message",
	})

	for _, sr := range r {
		s.writer.AppendRow(table.Row{
			sr.Command.String(),
			sr.Params,
			sr.Output,
			sr.Status.String(),
			sr.Message,
		})
	}

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
