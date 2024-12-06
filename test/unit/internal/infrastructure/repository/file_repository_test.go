package repository_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/fyvri/fresh-proxy-list/internal/infrastructure/repository"
	"github.com/fyvri/fresh-proxy-list/pkg/utils"
)

var (
	column1     = "Column1"
	column2     = "Column2"
	row1Column1 = "Row1Column1"
	row1Column2 = "Row1Column2"
	row2Column1 = "Row2Column1"
	row2Column2 = "Row2Column2"
)

func TestNewFileRepository(t *testing.T) {
	mockMkdirAll := func(path string, perm fs.FileMode) error {
		if path == "" {
			return errors.New("path cannot be empty")
		}
		return nil
	}
	mockCreate := func(name string) (io.Writer, error) {
		if name == "" {
			return nil, errors.New("file name cannot be empty")
		}
		return &bytes.Buffer{}, nil
	}
	mockCSVWriterUtil := &mockCSVWriterUtil{}
	fileRepository := repository.NewFileRepository(mockMkdirAll, mockCreate, mockCSVWriterUtil)

	if fileRepository == nil {
		t.Errorf(expectedReturnNonNil, "NewFileRepository", "FileRepositoryInterface")
	}

	r, ok := fileRepository.(*repository.FileRepository)
	if !ok {
		t.Errorf(expectedTypeAssertionErrorMessage, "*FileRepository")
	}

	if r.MkdirAll == nil {
		t.Errorf("expected mkdirAll to be set")
	}

	if r.Create == nil {
		t.Errorf("expected create to be set")
	}
}

func TestSaveFile(t *testing.T) {
	type fields struct {
		mkdirAll  func(path string, perm os.FileMode) error
		create    func(name string) (io.Writer, error)
		csvWriter utils.CSVWriterUtilInterface
	}

	type args struct {
		path   string
		data   interface{}
		format string
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		want      string
		wantError error
	}{
		{
			name: "CreateDirectoryError",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return errors.New("error creating directory")
				},
				create: func(name string) (io.Writer, error) {
					return nil, nil
				},
			},
			args: args{
				path:   testClassicFilePath + "." + testTXTExtension,
				data:   strings.Join(testIPs, "\n"),
				format: testTXTExtension,
			},
			want:      "",
			wantError: fmt.Errorf("error creating directory %v: %v", testClassicFilePath+"."+testTXTExtension, "error creating directory"),
		},
		{
			name: "CreateFileError",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return nil, errors.New("error creating file")
				},
			},
			args: args{
				path:   testClassicFilePath + "." + testTXTExtension,
				data:   testIPs,
				format: testTXTExtension,
			},
			want:      "",
			wantError: fmt.Errorf("error creating file %v: %v", testClassicFilePath+"."+testTXTExtension, "error creating file"),
		},
		{
			name: "UnsupportedFormat",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
			},
			args: args{
				path:   testClassicFilePath + "." + testTXTExtension,
				data:   testIPs,
				format: "unsupported-format",
			},
			want:      "",
			wantError: fmt.Errorf("unsupported format: %v", "unsupported-format"),
		},
		{
			name: "WriteTXTSuccess",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
			},
			args: args{
				path:   testClassicFilePath + "." + testTXTExtension,
				data:   testIPs,
				format: testTXTExtension,
			},
			want:      testIP1 + testIP2,
			wantError: nil,
		},
		{
			name: "WriteTXTError",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &mockWriter{
						errWrite: errors.New(testErrorWriting),
					}, nil
				},
			},
			args: args{
				path:   testClassicFilePath + "." + testTXTExtension,
				data:   testIPs,
				format: testTXTExtension,
			},
			want:      "",
			wantError: fmt.Errorf("error writing TXT: %v", testErrorWriting),
		},
		{
			name: "EncodeJSON",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
			},
			args: args{
				path:   testClassicFilePath + "." + testJSONExtension,
				data:   string(testProxiesToString),
				format: testJSONExtension,
			},
			want:      string(testProxiesToString),
			wantError: nil,
		},
		{
			name: "EncodeJSONError",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &mockWriter{
						errWrite: errors.New(testErrorWriting),
					}, nil
				},
			},
			args: args{
				path:   testClassicFilePath + "." + testJSONExtension,
				data:   testProxies,
				format: testJSONExtension,
			},
			want:      "",
			wantError: fmt.Errorf(testErrorEncode, "JSON", testErrorWriting),
		},
		{
			name: "EncodeCSVWithStringData",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
				csvWriter: &mockCSVWriterUtil{},
			},
			args: args{
				path:   testClassicFilePath + "." + testCSVExtension,
				data:   testIPs,
				format: testCSVExtension,
			},
			want:      string(testIPsToString) + "\n",
			wantError: nil,
		},
		{
			name: "EncodeCSVWithProxyData",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
				csvWriter: &mockCSVWriterUtil{},
			},
			args: args{
				path:   testAdvancedFilePath + "." + testCSVExtension,
				data:   testProxies,
				format: testCSVExtension,
			},
			want:      string(testProxiesToString) + "\n",
			wantError: nil,
		},
		{
			name: "EncodeCSVWithAdvancedProxyData",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
				csvWriter: &mockCSVWriterUtil{},
			},
			args: args{
				path:   testAdvancedFilePath + "." + testCSVExtension,
				data:   testAdvancedProxies,
				format: testCSVExtension,
			},
			want:      string(testAdvancedProxiesToString) + "\n",
			wantError: nil,
		},
		{
			name: "EncodeCSVWithErrorDataType",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
				csvWriter: &mockCSVWriterUtil{},
			},
			args: args{
				path:   testClassicFilePath + "." + testCSVExtension,
				data:   []error{},
				format: testCSVExtension,
			},
			want:      "",
			wantError: errors.New("invalid data type for CSV encoding"),
		},
		{
			name: "EncodeXMLWithStringStruct",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
			},
			args: args{
				path:   testClassicFilePath + "." + testXMLExtension,
				data:   testIPs,
				format: testXMLExtension,
			},
			want:      string(testIPsToString),
			wantError: nil,
		},
		{
			name: "EncodeXMLWithProxyStruct",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
			},
			args: args{
				path:   testClassicFilePath + "." + testXMLExtension,
				data:   testProxies,
				format: testXMLExtension,
			},
			want:      string(testProxiesToString),
			wantError: nil,
		},
		{
			name: "EncodeXMLWithAdvancedProxyStruct",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
			},
			args: args{
				path:   testClassicFilePath + "." + testXMLExtension,
				data:   testAdvancedProxies,
				format: testXMLExtension,
			},
			want:      string(testAdvancedProxiesToString),
			wantError: nil,
		},
		{
			name: "EncodeXMLError",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &mockWriter{
						errWrite: errors.New(testErrorWriting),
					}, nil
				},
			},
			args: args{
				path:   testClassicFilePath + "." + testXMLExtension,
				data:   testProxies,
				format: testXMLExtension,
			},
			want:      "",
			wantError: fmt.Errorf(testErrorEncode, "XML", testErrorWriting),
		},
		{
			name: "EncodeYAMLWithStringStruct",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
			},
			args: args{
				path:   testClassicFilePath + "." + testYAMLExtension,
				data:   testIPs,
				format: testYAMLExtension,
			},
			want:      string(testIPsToString),
			wantError: nil,
		},
		{
			name: "EncodeYAMLWithProxyStruct",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
			},
			args: args{
				path:   testClassicFilePath + "." + testYAMLExtension,
				data:   testProxies,
				format: testYAMLExtension,
			},
			want:      string(testProxiesToString),
			wantError: nil,
		},
		{
			name: "EncodeYAMLWithAdvancedProxyStruct",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &bytes.Buffer{}, nil
				},
			},
			args: args{
				path:   testClassicFilePath + "." + testYAMLExtension,
				data:   testAdvancedProxies,
				format: testYAMLExtension,
			},
			want:      string(testAdvancedProxiesToString),
			wantError: nil,
		},
		{
			name: "EncodeYAMLError",
			fields: fields{
				mkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
				create: func(name string) (io.Writer, error) {
					return &mockWriter{
						errWrite: errors.New(testErrorWriting),
					}, nil
				},
			},
			args: args{
				path:   testClassicFilePath + "." + testYAMLExtension,
				data:   testProxies,
				format: testYAMLExtension,
			},
			want:      "",
			wantError: fmt.Errorf(testErrorEncode, "YAML", "yaml: write error: error writing"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository.FileRepository{
				MkdirAll:  tt.fields.mkdirAll,
				Create:    tt.fields.create,
				CSVWriter: tt.fields.csvWriter,
			}
			err := r.SaveFile(tt.args.path, tt.args.data, tt.args.format)
			if (err != nil && tt.wantError != nil && err.Error() != tt.wantError.Error()) ||
				(err == nil && tt.wantError != nil) ||
				(err != nil && tt.wantError == nil) {
				t.Errorf(expectedErrorButGotMessage, "SaveFile()", tt.wantError, err)
			}
		})
	}
}

func TestWriteCSV(t *testing.T) {
	type fields struct {
		csvWriter utils.CSVWriterUtilInterface
	}

	type args struct {
		header []string
		rows   [][]string
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantError error
	}{
		{
			name: "Success",
			fields: fields{
				csvWriter: &mockCSVWriterUtil{},
			},
			args: args{
				header: []string{column1, column2},
				rows: [][]string{
					{row1Column1, row1Column2},
					{row2Column1, row2Column2},
				},
			},
			wantError: nil,
		},
		{
			name: "ErrorWritingHeader",
			fields: fields{
				csvWriter: &mockCSVWriterUtil{
					errWrite: errors.New("write header error"),
				},
			},
			args: args{
				header: []string{column1, column2},
				rows: [][]string{
					{row1Column1, row1Column2},
					{row2Column1, row2Column2},
				},
			},
			wantError: fmt.Errorf("failed to write header: %w", errors.New("write header error")),
		},
		{
			name: "ErrorWritingRow",
			fields: fields{
				csvWriter: &mockCSVWriterUtil{
					errWrite: errors.New("write row error"),
				},
			},
			args: args{
				header: nil,
				rows: [][]string{
					{row1Column1, row1Column2},
					{row2Column1, row2Column2},
				},
			},
			wantError: fmt.Errorf("failed to write row: %w", errors.New("write row error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			r := &repository.FileRepository{
				CSVWriter: tt.fields.csvWriter,
			}
			err := r.WriteCSV(&buf, tt.args.header, tt.args.rows)
			if (err != nil && tt.wantError != nil && err.Error() != tt.wantError.Error()) ||
				(err == nil && tt.wantError != nil) ||
				(err != nil && tt.wantError == nil) {
				t.Errorf(expectedErrorButGotMessage, "WriteCSV()", tt.wantError, err)
			}
		})
	}
}
