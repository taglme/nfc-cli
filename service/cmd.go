package service

import (
	"github.com/jedib0t/go-pretty/table"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-cli/actions"
	"github.com/taglme/nfc-client/pkg/client"
	"github.com/urfave/cli/v2"
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
