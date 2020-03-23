package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/taglme/nfc-cli/mock"
	"github.com/taglme/nfc-cli/models"
	"github.com/taglme/nfc-cli/opts"
	"github.com/urfave/cli/v2"
	"os"
	"testing"
	"time"
)

func Test_eventHandler(t *testing.T) {
	rep := mock.NewRepositoryService(nil)
	config := opts.Config{}
	cbCliStarted := func(string) {}
	app := New(rep, cbCliStarted, config)
	app.exitCh = make(chan struct{})
	app.ongoingJobs.left = 2
	app.ongoingJobs.published = 2

	app.eventHandler(models.EventRunSuccess, "test data")
	assert.Equal(t, 1, app.ongoingJobs.left)

	app.eventHandler(models.EventJobActivated, "test data")
	assert.Equal(t, 1, app.ongoingJobs.left)

	app.eventHandler(models.EventJobFinished, "test data")
	assert.Equal(t, 1, app.ongoingJobs.left)

	app.eventHandler(models.EventJobDeleted, "test data")
	assert.Equal(t, 1, app.ongoingJobs.left)

	app.output = "decorators_test_file.json"
	app.eventHandler(models.EventRunSuccess, "test data")
	assert.Equal(t, 0, app.ongoingJobs.left)
	err := os.Remove(app.output)
	assert.Nil(t, err)
	app.output = ""

	c1, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		app.eventHandler(models.EventJobFinished, "test data")
	}(c1)
	assert.Equal(t, 0, app.ongoingJobs.left)
	time.Sleep(time.Duration(1000) * time.Millisecond)
	cancel()
	select {
	case x, ok := <-app.exitCh:
		assert.True(t, ok)
		assert.NotNil(t, x)
	default:
		t.Error("Exit haven't been received")
	}
	close(app.exitCh)
}

func Test_errorHandler(t *testing.T) {
	rep := mock.NewRepositoryService(nil)
	config := opts.Config{}
	cbCliStarted := func(string) {}
	app := New(rep, cbCliStarted, config)
	app.exitCh = make(chan struct{})

	c1, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		app.errorHandler(errors.New("websocket closed on server side"))
	}(c1)
	time.Sleep(time.Duration(1000) * time.Millisecond)
	cancel()
	select {
	case x, ok := <-app.exitCh:
		assert.True(t, ok)
		assert.NotNil(t, x)
	default:
		t.Error("Exit haven't been received")
	}
	close(app.exitCh)
}

func Test_withAdapter(t *testing.T) {
	rep := mock.NewRepositoryService(nil)
	config := opts.Config{}
	cbCliStarted := func(string) {}
	app := New(rep, cbCliStarted, config)

	f := func(ctx *cli.Context) error {
		fmt.Println("Has been called.")
		assert.NotNil(t, ctx)
		return nil
	}

	ctx := cli.Context{}
	app.adapter = -1
	err := app.withAdapter(&ctx, f)
	assert.EqualError(t, err, "Can't find adapter with such index")

	app.adapter = 25
	err = app.withAdapter(&ctx, f)
	assert.EqualError(t, err, "Can't find adapter with such index")

	app.adapter = 1
	err = app.withAdapter(&ctx, f)
	assert.Nil(t, err)
}
