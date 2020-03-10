package service

import (
	"github.com/taglme/nfc-cli/models"
	"github.com/taglme/nfc-cli/repository"
	"github.com/urfave/cli/v2"
	"os"
	"sort"
)

type AppService interface {
	Start() error
	SetRepository(*repository.ApiService)
}

type appService struct {
	repository *repository.ApiService
	cliApp     cli.App

	exitCh chan struct{}

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

	cliStartedCb CbCliStarted
}

type CbCliStarted = func(url string)

func New(repository *repository.ApiService, cb CbCliStarted) AppService {
	return &appService{
		cliStartedCb: cb,
		repository:   repository,
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

	return s.cliApp.Run(os.Args)
}

func (s *appService) SetRepository(r *repository.ApiService) {
	s.repository = r
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
				s.flagsMap[models.FlagAdapter],
				s.flagsMap[models.FlagRepeat],
				s.flagsMap[models.FlagOutput],
				s.flagsMap[models.FlagAppend],
				s.flagsMap[models.FlagTimeout],
				s.flagsMap[models.FlagAuth],
			},
			Action: func(ctx *cli.Context) error {
				return s.withWsConnect(ctx, s.cmdRead)
			},
		},
		{
			Name:  models.CommandDump,
			Usage: "Dump tag memory",
			Flags: []cli.Flag{
				s.flagsMap[models.FlagHost],
				s.flagsMap[models.FlagAdapter],
				s.flagsMap[models.FlagRepeat],
				s.flagsMap[models.FlagOutput],
				s.flagsMap[models.FlagAppend],
				s.flagsMap[models.FlagTimeout],
				s.flagsMap[models.FlagAuth],
			},
			Action: func(ctx *cli.Context) error {
				return s.withWsConnect(ctx, s.cmdDump)
			},
		},
		{
			Name:  models.CommandLock,
			Usage: "Lock tag memory",
			Flags: []cli.Flag{
				s.flagsMap[models.FlagHost],
				s.flagsMap[models.FlagAdapter],
				s.flagsMap[models.FlagRepeat],
				s.flagsMap[models.FlagOutput],
				s.flagsMap[models.FlagAppend],
				s.flagsMap[models.FlagTimeout],
				s.flagsMap[models.FlagAuth],
			},
			Action: func(ctx *cli.Context) error {
				return s.withWsConnect(ctx, s.cmdLock)
			},
		},
	}
}
