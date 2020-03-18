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

	err = s.withAdapter(ctx, cmdFunc)
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
				fmt.Println("\nReceived done. Deleting adapter jobs...")

				err := s.repository.DeleteAdapterJobs(s.adapterId)
				if err != nil {
					log.Printf("Can't delete adapter jobs on exit: %s", err)
				}

				fmt.Println("Exiting...")
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

func (s *appService) withAdapter(ctx *cli.Context, cmdFunc func(*cli.Context) error) error {
	adapters, err := s.repository.GetAdapters()
	if err != nil {
		return err
	}
	if s.adapter <= 0 || s.adapter > len(adapters) {
		return errors.New("Can't find adapter with such index")
	}

	s.adapterId = adapters[s.adapter-1].AdapterID
	err = s.repository.DeleteAdapterJobs(s.adapterId)
	if err != nil {
		return errors.Wrap(err, "Can't delete adapter jobs: ")
	}
	log.Println("Adapter jobs where deleted")

	return cmdFunc(ctx)
}

func (s *appService) eventHandler(e models.Event, data interface{}) {
	s.cliStartedCb(s.host)

	if e == models.EventRunSuccess {
		s.ongoingJobs.left--

		if len(s.output) > 0 {
			err := s.writeToFile(s.output, data)
			if err != nil {
				log.Println("Can't write to the file: ", err)
			}
		}
	}

	if (e == models.EventJobFinished || e == models.EventJobDeleted) && (s.ongoingJobs.left < 1) {
		fmt.Println("Deleting all jobs for the adapter...")
		err := s.repository.DeleteAdapterJobs(s.adapterId)
		if err != nil {
			log.Printf("Can't delete adapter jobs on exit: %s", err)
		}

		fmt.Println("Exiting...")
		s.exitCh <- struct{}{}
	}
}
