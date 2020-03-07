package service

import (
	"github.com/jedib0t/go-pretty/table"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-cli/actions"
	"github.com/taglme/nfc-cli/models"
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

	writer table.Writer
	cliApp cli.App

	flagsMap map[string]cli.Flag
	//  below is arguments controlled by ./flags.go
	host    string
	adapter int
	repeat  int
	output  string
	append  bool
	timeout int
	input   string
	auth    string
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
	s.flagsMap = s.getFlagsMap()
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

func (s *appService) getCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  models.CommandVersion,
			Usage: "Application version",
			Flags: []cli.Flag{
				s.flagsMap[models.FlagHost],
			},
			Action: s.cmdVersion,
		},
		{
			Name:  models.CommandAdapters,
			Usage: "Get adapters list",
			Flags: []cli.Flag{
				s.flagsMap[models.FlagHost],
			},
			Action: s.cmdAdapters,
		},
		{
			Name:  models.CommandRead,
			Usage: "Read tag data with NDEF message",
			Flags: []cli.Flag{
				s.flagsMap[models.FlagHost],
				s.flagsMap[models.FlagAdapters],
				s.flagsMap[models.FlagRepeat],
				s.flagsMap[models.FlagOutput],
				s.flagsMap[models.FlagAppend],
				s.flagsMap[models.FlagTimeout],
				s.flagsMap[models.FlagAuth],
			},
			Action: s.cmdRead,
		},
	}
}
