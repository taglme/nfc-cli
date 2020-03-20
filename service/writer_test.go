package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/taglme/nfc-cli/opts"
	"github.com/taglme/nfc-cli/repository"
	"os"
	"testing"
)

func TestAppService_WriteToFile(t *testing.T) {
	var rep *repository.RepositoryService
	config := opts.Config{}
	cbCliStarted := func(string) {}

	app := New(rep, cbCliStarted, config)
	filename := "writer_test_file.json"
	err := app.writeToFile(filename, "any data could be written")
	assert.Nil(t, err)
	err = os.Remove(filename)
	assert.Nil(t, err)
}
