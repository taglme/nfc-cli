package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/taglme/nfc-cli/opts"
	"github.com/taglme/nfc-cli/repository"
	"testing"
)

func TestAppService_getFlagsMap(t *testing.T) {
	var rep *repository.RepositoryService
	config := opts.Config{}
	cbCliStarted := func(string) {}

	app := New(rep, cbCliStarted, config)
	flags := app.getFlagsMap()

	for i, f := range flags {
		assert.NotEmpty(t, i)
		assert.NotEmpty(t, f)
	}
}
