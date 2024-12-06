package utils

import (
	"encoding/csv"
	"io"
)

type CSVWriterUtil struct {
}

type CSVWriterUtilInterface interface {
	Init(w io.Writer) *csv.Writer
	Flush(csvWriter *csv.Writer)
	Write(csvWriter *csv.Writer, record []string) error
}

func NewCSVWriter() CSVWriterUtilInterface {
	return &CSVWriterUtil{}
}

func (u *CSVWriterUtil) Init(w io.Writer) *csv.Writer {
	return csv.NewWriter(w)
}

func (u *CSVWriterUtil) Flush(csvWriter *csv.Writer) {
	csvWriter.Flush()
}

func (u *CSVWriterUtil) Write(csvWriter *csv.Writer, record []string) error {
	return csvWriter.Write(record)
}
