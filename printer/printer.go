package printer

import (
	"github.com/jedib0t/go-pretty/table"
)

type PrinterService interface {

}

type printerService struct {
	writer    table.Writer
}

func New(writer table.Writer) PrinterService {
	return &printerService{
		writer: writer,
	}
}