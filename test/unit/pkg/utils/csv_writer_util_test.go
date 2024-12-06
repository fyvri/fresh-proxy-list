package util_test

import (
	"bytes"
	"encoding/csv"
	"io"
	"reflect"
	"testing"

	"github.com/fyvri/fresh-proxy-list/pkg/utils"
)

var (
	writerBufferSize = 1024
)

func TestNewCSVWriter(t *testing.T) {
	tests := []struct {
		name string
		want utils.CSVWriterUtilInterface
	}{
		{
			name: "Success",
			want: &utils.CSVWriterUtil{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csvWriterUtil := utils.NewCSVWriter()
			if csvWriterUtil == nil {
				t.Errorf(expectedReturnNonNil, "NewCSVWriter", "CSVWriterUtilInterface")
			}

			got, ok := csvWriterUtil.(*utils.CSVWriterUtil)
			if !ok {
				t.Errorf(expectedTypeAssertionErrorMessage, "*utils.CSVWriterUtil")
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf(expectedButGotMessage, "*utils.CSVWriterUtil", tt.want, got)
			}
		})
	}
}

func TestInit(t *testing.T) {
	type args struct {
		writer io.Writer
	}

	tests := []struct {
		name string
		args args
		want *csv.Writer
	}{
		{
			name: "Success",
			args: args{
				writer: bytes.NewBuffer(make([]byte, writerBufferSize)),
			},
			want: csv.NewWriter(bytes.NewBuffer(make([]byte, writerBufferSize))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := utils.NewCSVWriter()
			got := u.Init(tt.args.writer)
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf(expectedButGotMessage, "Init()", reflect.TypeOf(tt.want), reflect.TypeOf(got))
			}
		})
	}
}

func TestFlush(t *testing.T) {
	type setup struct {
		newCSVWriter func() *utils.CSVWriterUtil
	}

	type args struct {
		csvWriter *csv.Writer
	}

	tests := []struct {
		name  string
		setup setup
		args  args
	}{
		{
			name: "Success",
			setup: setup{
				newCSVWriter: func() *utils.CSVWriterUtil {
					return utils.NewCSVWriter().(*utils.CSVWriterUtil)
				},
			},
			args: args{
				csvWriter: csv.NewWriter(bytes.NewBuffer(make([]byte, writerBufferSize))),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := tt.setup.newCSVWriter()
			u.Flush(tt.args.csvWriter)
		})
	}
}

func TestWrite(t *testing.T) {
	type setup struct {
		newCSVWriter func() *utils.CSVWriterUtil
	}

	type args struct {
		writer io.Writer
		record []string
	}

	tests := []struct {
		name      string
		setup     setup
		args      args
		wantError error
	}{
		{
			name: "Success",
			setup: setup{
				newCSVWriter: func() *utils.CSVWriterUtil {
					return utils.NewCSVWriter().(*utils.CSVWriterUtil)
				},
			},
			args: args{
				writer: &bytes.Buffer{},
				record: []string{"a", "b", "c"},
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := tt.setup.newCSVWriter()
			csvWriter := u.Init(tt.args.writer)
			err := u.Write(csvWriter, tt.args.record)
			if (err != nil && tt.wantError != nil && err.Error() != tt.wantError.Error()) ||
				(err == nil && tt.wantError != nil) ||
				(err != nil && tt.wantError == nil) {
				t.Errorf(expectedErrorButGotMessage, "Write()", tt.wantError, err)
			}
		})
	}
}
