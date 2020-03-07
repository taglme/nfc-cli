package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-cli/models"
	"github.com/urfave/cli/v2"
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

func (s *appService) eventHandler(e models.Event) {
	s.cliStartedCb(s.host)
	if e == models.EventJobFinished || e == models.EventJobDeleted {
		s.exitCh <- struct{}{}
	}
}

func (s *appService) cmdRead(*cli.Context) error {
	s.cliStartedCb(s.host)
	c1, cancel := context.WithCancel(context.Background())
	s.exitCh = make(chan struct{})
	err := s.repository.RunWsConnection(s.eventHandler)
	if err != nil {
		return errors.Wrap(err, "Can't establish the WS connection")
	}
	defer s.repository.StopWsConnection()

	// TODO: might be adapterId
	_, err = s.repository.AddJob(models.CommandRead, "7f1d71b6-875a-463e-a835-707cebe1bc8c", s.repeat, s.timeout)

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
