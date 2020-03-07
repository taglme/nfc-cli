package service

import (
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-cli/actions"
	"github.com/taglme/nfc-cli/models"
	"github.com/taglme/nfc-client/pkg/client"
	"github.com/urfave/cli/v2"
	"os"
)

func (s *appService) cmdVersion(*cli.Context) error {
	s.nfcClient = client.New(s.host, "en")
	s.actions = actions.New(s.nfcClient)

	info, err := s.actions.GetVersion()
	if err != nil {
		return errors.Wrap(err, "Can't get application info")
	}

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

	return nil
}

func (s *appService) cmdAdapters(*cli.Context) error {
	s.nfcClient = client.New(s.host, "en")
	s.actions = actions.New(s.nfcClient)

	adapters, err := s.actions.GetAdapters()
	if err != nil {
		return errors.Wrap(err, "Can't get adapters list")
	}

	s.writer.AppendHeader(table.Row{"Adapter ID", "Name", "Type", "Driver"})

	for _, a := range adapters {
		s.writer.AppendRow(table.Row{a.AdapterID, a.Name, a.Type.String(), a.Driver})
	}
	s.writer.SetStyle(table.StyleLight)
	s.writer.Render()

	return nil
}

func (s *appService) cmdRead(ctx *cli.Context) error {
	s.nfcClient = client.New(s.host, "en")
	s.actions = actions.New(s.nfcClient)

	// TODO: might be adapterId
	nJ, err := s.actions.AddJob(models.CommandRead, "7f1d71b6-875a-463e-a835-707cebe1bc8c", s.timeout, s.repeat)
	if err != nil {
		return errors.Wrap(err, "Can't add read job")
	}

	fmt.Println("New job has been submitted:")

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
		nJ.JobID,
		nJ.JobName,
		nJ.Status.String(),
		nJ.AdapterID,
		nJ.AdapterName,
		nJ.Repeat,
		nJ.TotalRuns,
		nJ.SuccessRuns,
		nJ.ErrorRuns,
		nJ.ExpireAfter,
		nJ.CreatedAt.String(),
	})

	s.writer.Render()

	// might be moved to some Printer serviice
	fmt.Println("\nJob steps:")
	s.writer = table.NewWriter()
	s.writer.SetOutputMirror(os.Stdout)
	s.writer.SetStyle(table.StyleRounded)

	s.writer.AppendHeader(table.Row{
		"Command",
		"Command params",
	})

	for _, step := range nJ.Steps {
		s.writer.AppendRow(table.Row{
			step.Command,
			step.Params,
		})
	}

	s.writer.Render()

	return nil
}
