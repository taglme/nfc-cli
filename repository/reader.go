package repository

import (
	"bufio"
	"encoding/json"
	"github.com/pkg/errors"
	apiModels "github.com/taglme/nfc-goclient/pkg/models"
	"os"
)

func (s *RepositoryService) readFromFile(filename string) (data []apiModels.NewJob, err error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, errors.Wrap(err, "Can't open the file: ")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		nj := apiModels.NewJob{}
		err := json.Unmarshal([]byte(scanner.Text()), &nj)
		if err != nil {
			return data, errors.Wrap(err, "Can't unmarshall the data on reading from the file")
		}
		data = append(data, nj)
	}

	return data, err
}
