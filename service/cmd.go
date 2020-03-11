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

func (s *appService) cmdRead(ctx *cli.Context) error {
	fmt.Println("Available adapters:")
	adapters, err := s.repository.GetAdapters()
	if err != nil {
		return err
	}
	if s.adapter <= 0 || s.adapter > len(adapters) {
		return errors.New("Can't find adapter with such index")
	}

	auth, err := s.parseHexString(s.auth)
	if err != nil {
		log.Println("Can't parse auth string")
	}

	export := ctx.Bool(models.FlagExport)

	var nj interface{}
	_, nj, err = s.repository.AddGenericJob(
		models.GenericJobParams{
			Cmd:       models.CommandRead,
			AdapterId: adapters[s.adapter-1].AdapterID,
			Repeat:    s.repeat,
			Expire:    s.timeout,
			Auth:      auth,
			Export:    export,
			JobName: s.jobName,
		},
	)

	if export && nj != nil {
		err := s.writeToFile(s.output, nj)
		if err != nil {
			return errors.Wrapf(err, "Can't write to the file: ")
		}
	}

	return err
}

func (s *appService) cmdDump(ctx *cli.Context) error {
	fmt.Println("Available adapters:")
	adapters, err := s.repository.GetAdapters()
	if err != nil {
		return err
	}
	if s.adapter <= 0 || s.adapter > len(adapters) {
		return errors.New("Can't find adapter with such index")
	}

	auth, err := s.parseHexString(s.auth)
	if err != nil {
		log.Println("Can't parse auth string")
	}

	export := ctx.Bool(models.FlagExport)
	var nj interface{}
	_, nj, err = s.repository.AddGenericJob(
		models.GenericJobParams{
			Cmd:       models.CommandDump,
			AdapterId: adapters[s.adapter-1].AdapterID,
			Repeat:    s.repeat,
			Expire:    s.timeout,
			Auth:      auth,
			Export:    export,
			JobName: s.jobName,
		},
	)
	if export && nj != nil {
		err := s.writeToFile(s.output, nj)
		if err != nil {
			return errors.Wrapf(err, "Can't write to the file: ")
		}
	}
	return err
}

func (s *appService) cmdLock(ctx *cli.Context) error {
	fmt.Println("Available adapters:")
	adapters, err := s.repository.GetAdapters()
	if err != nil {
		return err
	}
	if s.adapter <= 0 || s.adapter > len(adapters) {
		return errors.New("Can't find adapter with such index")
	}

	auth, err := s.parseHexString(s.auth)
	if err != nil {
		log.Println("Can't parse auth string")
	}

	export := ctx.Bool(models.FlagExport)

	var nj interface{}
	_, nj, err = s.repository.AddGenericJob(
		models.GenericJobParams{
			Cmd:       models.CommandLock,
			AdapterId: adapters[s.adapter-1].AdapterID,
			Repeat:    s.repeat,
			Expire:    s.timeout,
			Auth:      auth,
			JobName: s.jobName,
		},
	)
	if export && nj != nil {
		err := s.writeToFile(s.output, nj)
		if err != nil {
			return errors.Wrapf(err, "Can't write to the file: ")
		}
	}
	return err
}

func (s *appService) cmdFormat(ctx *cli.Context) error {
	fmt.Println("Available adapters:")
	adapters, err := s.repository.GetAdapters()
	if err != nil {
		return err
	}
	if s.adapter <= 0 || s.adapter > len(adapters) {
		return errors.New("Can't find adapter with such index")
	}

	auth, err := s.parseHexString(s.auth)
	if err != nil {
		log.Println("Can't parse auth string")
	}

	export := ctx.Bool(models.FlagExport)

	var nj interface{}
	_, nj, err = s.repository.AddGenericJob(
		models.GenericJobParams{
			Cmd:       models.CommandFormat,
			AdapterId: adapters[s.adapter-1].AdapterID,
			Repeat:    s.repeat,
			Expire:    s.timeout,
			Auth:      auth,
			Export:    export,
			JobName: s.jobName,
		},
	)
	if export && nj != nil {
		err := s.writeToFile(s.output, nj)
		if err != nil {
			return errors.Wrapf(err, "Can't write to the file: ")
		}
	}
	return err
}

func (s *appService) cmdRmPwd(ctx *cli.Context) error {
	fmt.Println("Available adapters:")
	adapters, err := s.repository.GetAdapters()
	if err != nil {
		return err
	}
	if s.adapter <= 0 || s.adapter > len(adapters) {
		return errors.New("Can't find adapter with such index")
	}

	auth, err := s.parseHexString(s.auth)
	if err != nil {
		log.Println("Can't parse auth string")
	}

	export := ctx.Bool(models.FlagExport)

	var nj interface{}
	_, nj, err = s.repository.AddGenericJob(
		models.GenericJobParams{
			Cmd:       models.CommandRmpwd,
			AdapterId: adapters[s.adapter-1].AdapterID,
			Repeat:    s.repeat,
			Expire:    s.timeout,
			Auth:      auth,
			Export:    export,
			JobName: s.jobName,
		},
	)
	if export && nj != nil {
		err := s.writeToFile(s.output, nj)
		if err != nil {
			return errors.Wrapf(err, "Can't write to the file: ")
		}
	}
	return err
}

func (s *appService) cmdSetPwd(ctx *cli.Context) error {
	fmt.Println("Available adapters:")
	adapters, err := s.repository.GetAdapters()
	if err != nil {
		return err
	}
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

	export := ctx.Bool(models.FlagExport)

	var nj interface{}
	_, nj, err = s.repository.AddSetPwdJob(
		models.GenericJobParams{
			Cmd:       models.CommandSetpwd,
			AdapterId: adapters[s.adapter-1].AdapterID,
			Repeat:    s.repeat,
			Expire:    s.timeout,
			Auth:      auth,
			Export:    export,
			JobName: s.jobName,
		},
		password,
	)
	if export && nj != nil {
		err := s.writeToFile(s.output, nj)
		if err != nil {
			return errors.Wrapf(err, "Can't write to the file: ")
		}
	}
	return err
}

func (s *appService) cmdTransmit(ctx *cli.Context) error {
	fmt.Println("Available adapters:")
	adapters, err := s.repository.GetAdapters()
	if err != nil {
		return err
	}
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

	export := ctx.Bool(models.FlagExport)

	var nj interface{}
	_, nj, err = s.repository.AddTransmitJob(
		models.GenericJobParams{
			Cmd:       models.CommandTransmit,
			AdapterId: adapters[s.adapter-1].AdapterID,
			Repeat:    s.repeat,
			Expire:    s.timeout,
			Auth:      auth,
			Export:    export,
			JobName: s.jobName,
		},
		txBytes,
		target,
	)
	if export && nj != nil {
		err := s.writeToFile(s.output, nj)
		if err != nil {
			return errors.Wrapf(err, "Can't write to the file: ")
		}
	}
	return err
}

func (s *appService) cmdWrite(ctx *cli.Context) error {
	fmt.Println("Available adapters:")
	adapters, err := s.repository.GetAdapters()
	if err != nil {
		return err
	}
	if s.adapter <= 0 || s.adapter > len(adapters) {
		return errors.New("Can't find adapter with such index")
	}

	auth, err := s.parseHexString(s.auth)
	if err != nil {
		log.Println("Can't parse auth string")
	}

	payload, err := s.parseNdefPayloadFlags(ctx)
	if err != nil {
		return err
	}

	protect := ctx.Bool(models.FlagProtect)

	export := ctx.Bool(models.FlagExport)

	var nj interface{}
	_, nj, err = s.repository.AddWriteJob(
		models.GenericJobParams{
			Cmd:       models.CommandTransmit,
			AdapterId: adapters[s.adapter-1].AdapterID,
			Repeat:    s.repeat,
			Expire:    s.timeout,
			Auth:      auth,
			Export:    export,
			JobName: s.jobName,
		},
		payload,
		protect,
	)
	if export && nj != nil {
		err := s.writeToFile(s.output, nj)
		if err != nil {
			return errors.Wrapf(err, "Can't write to the file: ")
		}
	}

	return err
}

func (s *appService) cmdRun(ctx *cli.Context) error {
	fmt.Println("Available adapters:")
	adapters, err := s.repository.GetAdapters()
	if err != nil {
		return err
	}
	if s.adapter <= 0 || s.adapter > len(adapters) {
		return errors.New("Can't find adapter with such index")
	}

	file := ctx.String(models.FlagFile)

	_, _, err = s.repository.AddJobFromFile(adapters[s.adapter-1].AdapterID, file, models.GenericJobParams{Expire: s.timeout, JobName: s.jobName})

	return err
}

func (s *appService) withWsConnect(ctx *cli.Context, cmdFunc func(*cli.Context) error) error {
	s.cliStartedCb(s.host)

	export := ctx.Bool(models.FlagExport)
	if export {
		err := cmdFunc(ctx)
		if err != nil {
			ctx.Done()
			return err
		}
		return nil
	}

	err := s.repository.RunWsConnection(s.eventHandler)
	if err != nil {
		return errors.Wrap(err, "Can't establish the WS connection")
	}
	defer func() {
		err = s.repository.StopWsConnection()
		if err != nil {
			log.Printf("Error on WS connection close: %s", err)
		}
	}()

	err = cmdFunc(ctx)
	if err != nil {
		ctx.Done()
		return err
	}

	c1, cancel := context.WithCancel(context.Background())
	s.exitCh = make(chan struct{})
	go func(ctx context.Context) {
		fmt.Println("Waiting for Job deleted event. Press ^C to stop.")
		for {
			select {
			case <-ctx.Done():
				fmt.Println("\nReceived done, exiting...")
				s.exitCh <- struct{}{}
				return
			default:
				time.Sleep(50 * time.Millisecond)
			}
		}
	}(c1)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		for range signalCh {
			cancel()
			return
		}
	}()
	<-s.exitCh

	return err
}
