package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-cli/models"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/signal"
	"time"
)

func (s *appService) cmdVersion(*cli.Context) error {
	s.cliStartedCb(s.host)
	_, err := s.repository.GetVersion()

	return err
}

func (s *appService) cmdAdapters(*cli.Context) error {
	s.cliStartedCb(s.host)
	_, err := s.repository.GetAdapters()

	return err
}

func (s *appService) eventHandler(e models.Event, data interface{}) {
	s.cliStartedCb(s.host)

	if e == models.EventJobFinished && len(s.output) > 0 {
		err := s.writeToFile(s.output, data)
		if err != nil {
			log.Println("Can't write to the file: ", err)
		}
	}

	if e == models.EventJobFinished || e == models.EventJobDeleted {
		s.exitCh <- struct{}{}
	}
}

func (s *appService) cmdRead(*cli.Context) error {
	fmt.Println("Available adapters:")
	adapters, err := s.repository.GetAdapters()
	if s.adapter <= 0 || s.adapter > len(adapters) {
		return errors.New("Can't find adapter with such index")
	}

	pwd, err := s.parseHexString(s.auth)
	if err != nil {
		log.Println("Can't parse password string")
	}

	_, err = s.repository.AddGenericJob(models.CommandRead, adapters[s.adapter-1].AdapterID, s.repeat, s.timeout, pwd)
	return err
}

func (s *appService) cmdDump(*cli.Context) error {
	fmt.Println("Available adapters:")
	adapters, err := s.repository.GetAdapters()
	if s.adapter <= 0 || s.adapter > len(adapters) {
		return errors.New("Can't find adapter with such index")
	}

	pwd, err := s.parseHexString(s.auth)
	if err != nil {
		log.Println("Can't parse password string")
	}

	_, err = s.repository.AddGenericJob(models.CommandDump, adapters[s.adapter-1].AdapterID, s.repeat, s.timeout, pwd)
	return err
}

func (s *appService) cmdLock(*cli.Context) error {
	fmt.Println("Available adapters:")
	adapters, err := s.repository.GetAdapters()
	if s.adapter <= 0 || s.adapter > len(adapters) {
		return errors.New("Can't find adapter with such index")
	}

	pwd, err := s.parseHexString(s.auth)
	if err != nil {
		log.Println("Can't parse password string")
	}

	_, err = s.repository.AddGenericJob(models.CommandLock, adapters[s.adapter-1].AdapterID, s.repeat, s.timeout, pwd)
	return err
}

func (s *appService) cmdFormat(*cli.Context) error {
	fmt.Println("Available adapters:")
	adapters, err := s.repository.GetAdapters()
	if s.adapter <= 0 || s.adapter > len(adapters) {
		return errors.New("Can't find adapter with such index")
	}

	pwd, err := s.parseHexString(s.auth)
	if err != nil {
		log.Println("Can't parse password string")
	}

	_, err = s.repository.AddGenericJob(models.CommandFormat, adapters[s.adapter-1].AdapterID, s.repeat, s.timeout, pwd)
	return err
}

func (s *appService) cmdRmPwd(*cli.Context) error {
	fmt.Println("Available adapters:")
	adapters, err := s.repository.GetAdapters()
	if s.adapter <= 0 || s.adapter > len(adapters) {
		return errors.New("Can't find adapter with such index")
	}

	pwd, err := s.parseHexString(s.auth)
	if err != nil {
		log.Println("Can't parse password string")
	}

	_, err = s.repository.AddGenericJob(models.CommandRmpwd, adapters[s.adapter-1].AdapterID, s.repeat, s.timeout, pwd)
	return err
}

func (s *appService) cmdSetPwd(ctx *cli.Context) error {
	fmt.Println("Available adapters:")
	adapters, err := s.repository.GetAdapters()
	if s.adapter <= 0 || s.adapter > len(adapters) {
		return errors.New("Can't find adapter with such index")
	}

	password, err := s.parseHexString(ctx.String(models.FlagPwd))
	if err != nil {
		return errors.Wrap(err, "Can't parse password arg: ")
	}

	auth, err := s.parseHexString(s.auth)
	if err != nil {
		log.Println("Can't parse auth string")
	}

	_, err = s.repository.AddSetPwdJob(adapters[s.adapter-1].AdapterID, s.repeat, s.timeout, auth, password)
	return err
}

func (s *appService) cmdTransmit(ctx *cli.Context) error {
	fmt.Println("Available adapters:")
	adapters, err := s.repository.GetAdapters()
	if s.adapter <= 0 || s.adapter > len(adapters) {
		return errors.New("Can't find adapter with such index")
	}

	target := ctx.String(models.FlagTarget)
	if target != "tag" && target != "adapter" {
		return errors.New("Wrong target flag value. Can be either \"tag\" or \"adapter\".")
	}

	txBytes, err := s.parseHexString(ctx.String(models.FlagTxBytes))
	if err != nil {
		return errors.Wrap(err, "Can't parse password arg: ")
	}

	auth, err := s.parseHexString(s.auth)
	if err != nil {
		log.Println("Can't parse auth string")
	}

	_, err = s.repository.AddTransmitJob(adapters[s.adapter-1].AdapterID, s.repeat, s.timeout, auth, txBytes, target)
	return err
}

func (s *appService) withWsConnect(ctx *cli.Context, cmdFunc func(*cli.Context) error) error {
	s.cliStartedCb(s.host)
	c1, cancel := context.WithCancel(context.Background())
	s.exitCh = make(chan struct{})
	err := s.repository.RunWsConnection(s.eventHandler)
	if err != nil {
		return errors.Wrap(err, "Can't establish the WS connection")
	}
	defer s.repository.StopWsConnection()

	err = cmdFunc(ctx)
	if err != nil {
		ctx.Done()
		return err
	}

	go func(ctx context.Context) {
		fmt.Println("Waiting for Job deleted event. Press ^C to stop.")
		for {
			select {
			case <-ctx.Done():
				fmt.Println("\nReceived done, exiting...")
				time.Sleep(50 * time.Millisecond)
				s.exitCh <- struct{}{}
				return
			default:
			}
		}
	}(c1)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		select {
		case <-signalCh:
			cancel()
			return
		}
	}()
	<-s.exitCh

	return err
}
