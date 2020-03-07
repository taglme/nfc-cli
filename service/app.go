package service

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-cli/actions"
	"github.com/taglme/nfc-client/pkg/client"
	"github.com/urfave/cli/v2"
	"os"
)

type AppService interface {
	Start() error
	GetHost() string
	SetActionService(actions.ActionService)
}

type appService struct {
	actions   actions.ActionService
	nfcClient *client.Client
	host      string
	adapter   int
	cliApp    cli.App
}

func New() AppService {
	return &appService{
		cliApp: cli.App{
			Name:        "nfc-cli",
			Version:     "v0.0.1",
			Description: "Cross-platform CLI for reading NFC tags ",
		},
	}
}

func (s *appService) Start() error {
	s.cliApp.Commands = s.getCommands()

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
			Name:  "adapters",
			Usage: "Get adapters list",
			Flags: s.getFlags(),
			Action: func(c *cli.Context) error {
				fmt.Println(s.host)
				return nil
			},
		},
	}
}

func (s *appService) cmdVersion(ctx *cli.Context) error {
	s.nfcClient = client.New(s.host, "en")
	s.actions = actions.New(s.nfcClient)

	info, err := s.actions.GetVersion()
	if err != nil {
		return errors.Wrap(err, "Can't get application info")
	}
	fmt.Println(info)
	return nil
}
