package service

import (
	"github.com/jedib0t/go-pretty/table"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-cli/actions"
	"github.com/taglme/nfc-client/pkg/client"
	"github.com/urfave/cli/v2"
	"os"
	"sort"
)

type AppService interface {
	Start() error
	GetHost() string
	SetActionService(actions.ActionService)
}

type appService struct {
	actions   actions.ActionService
	nfcClient *client.Client
	writer    table.Writer
	host      string
	adapter   int
	cliApp    cli.App
}

func New(writer table.Writer) AppService {
	return &appService{
		writer: writer,
		cliApp: cli.App{
			Name:        "nfc-cli",
			Version:     "v0.0.1",
			Description: "Cross-platform CLI for reading NFC tags ",
		},
	}
}

func (s *appService) Start() error {
	s.cliApp.Commands = s.getCommands()

	sort.Sort(cli.FlagsByName(s.cliApp.Flags))
	sort.Sort(cli.CommandsByName(s.cliApp.Commands))

	err := s.cliApp.Run(os.Args)

	if err != nil {
		return errors.Wrap(err, " Can't start thee cli application:\n")
	}
	return nil
}

func (s *appService) GetHost() string {
	return s.host
}

func (s *appService) SetActionService(a actions.ActionService) {
	s.actions = a
}

func (s *appService) getFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "host",
			Value:       "127.0.0.1:3011",
			Usage:       "Target host and port",
			Destination: &s.host,
		},
		&cli.IntFlag{
			Name:        "adapter",
			Value:       1,
			Usage:       "Adapter",
			Aliases:     []string{"a"},
			Destination: &s.adapter,
		},
	}
}

func (s *appService) getCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "version",
			Usage:  "Application version",
			Flags:  s.getFlags(),
			Action: s.cmdVersion,
		},
		{
			Name:   "adapters",
			Usage:  "Get adapters list",
			Flags:  s.getFlags(),
			Action: s.cmdAdapters,
		},
	}
}
