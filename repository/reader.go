package repository

import (
	"encoding/json"
	"github.com/pkg/errors"
	apiModels "github.com/taglme/nfc-goclient/pkg/models"
	"os"
)

func (s *ApiService) readFromFile(filename string) (data *apiModels.NewJob, err error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, errors.Wrap(err, "Can't open the file: ")
	}

	defer file.Close()

	encoder := json.NewDecoder(file)
	err = encoder.Decode(&data)
	if err != nil {
		return data, errors.Wrap(err, "Can't encode the data on reading from the file: ")
	}

	return data, err
}
