package repository

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/fyvri/fresh-proxy-list/internal/entity"
	"github.com/fyvri/fresh-proxy-list/pkg/utils"

	"gopkg.in/yaml.v3"
)

type FileRepository struct {
	MkdirAll  func(path string, perm os.FileMode) error
	Create    func(name string) (io.Writer, error)
	CSVWriter utils.CSVWriterUtilInterface
}

type FileRepositoryInterface interface {
	SaveFile(filePath string, data interface{}, format string) error
	CreateDirectory(filePath string) error
	WriteTxt(writer io.Writer, data interface{}) error
	EncodeCSV(writer io.Writer, data interface{}) error
	WriteCSV(writer io.Writer, header []string, rows [][]string) error
	EncodeJSON(writer io.Writer, data interface{}) error
	EncodeXML(writer io.Writer, data interface{}) error
	EncodeYAML(writer io.Writer, data interface{}) error
}

type MkdirAllFunc func(path string, perm os.FileMode) error
type CreateFunc func(name string) (io.Writer, error)

func NewFileRepository(mkdirAll MkdirAllFunc, create CreateFunc, csvWriter utils.CSVWriterUtilInterface) FileRepositoryInterface {
	return &FileRepository{
		MkdirAll:  mkdirAll,
		Create:    create,
		CSVWriter: csvWriter,
	}
}

func (r *FileRepository) SaveFile(filePath string, data interface{}, format string) error {
	if err := r.CreateDirectory(filePath); err != nil {
		return err
	}

	file, err := r.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", filePath, err)
	}
	defer func() {
		if f, ok := file.(io.Closer); ok {
			f.Close()
		}
	}()

	switch format {
	case "txt":
		return r.WriteTxt(file, data)
	case "json":
		return r.EncodeJSON(file, data)
	case "csv":
		return r.EncodeCSV(file, data)
	case "xml":
		return r.EncodeXML(file, data)
	case "yaml":
		return r.EncodeYAML(file, data)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func (r *FileRepository) CreateDirectory(filePath string) error {
	err := r.MkdirAll(filepath.Dir(filePath), fs.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory %s: %v", filePath, err)
	}
	return nil
}

func (r *FileRepository) WriteTxt(writer io.Writer, data interface{}) error {
	var dataString string
	if stringData, ok := data.([]string); ok {
		dataString = strings.Join(stringData, "\n")
	}

	_, err := writer.Write([]byte(dataString))
	if err != nil {
		return fmt.Errorf("error writing TXT: %v", err)
	}
	return nil
}

func (r *FileRepository) EncodeCSV(writer io.Writer, data interface{}) error {
	switch proxyData := data.(type) {
	case []string:
		rows := make([][]string, len(proxyData))
		for i, rowElem := range proxyData {
			rows[i] = []string{rowElem}
		}
		return r.WriteCSV(writer, nil, rows)
	case []entity.Proxy:
		header := []string{"Proxy", "IP", "Port", "TimeTaken", "CheckedAt"}
		rows := make([][]string, len(proxyData))
		for i, proxy := range proxyData {
			rows[i] = []string{proxy.Proxy, proxy.IP, proxy.Port, fmt.Sprintf("%v", proxy.TimeTaken), proxy.CheckedAt}
		}
		return r.WriteCSV(writer, header, rows)
	case []entity.AdvancedProxy:
		header := []string{"Proxy", "IP", "Port", "Categories", "TimeTaken", "CheckedAt"}
		rows := make([][]string, len(proxyData))
		for i, proxy := range proxyData {
			rows[i] = []string{proxy.Proxy, proxy.IP, proxy.Port, strings.Join(proxy.Categories, ","), fmt.Sprintf("%v", proxy.TimeTaken), proxy.CheckedAt}
		}
		return r.WriteCSV(writer, header, rows)
	default:
		return fmt.Errorf("invalid data type for CSV encoding")
	}
}

func (r *FileRepository) WriteCSV(writer io.Writer, header []string, rows [][]string) error {
	csvWriter := r.CSVWriter.Init(writer)
	defer r.CSVWriter.Flush(csvWriter)

	if header != nil {
		if err := r.CSVWriter.Write(csvWriter, header); err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}
	}

	for _, row := range rows {
		if err := r.CSVWriter.Write(csvWriter, row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return nil
}

func (r *FileRepository) EncodeJSON(writer io.Writer, data interface{}) error {
	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}
	return nil
}

func (r *FileRepository) EncodeXML(writer io.Writer, data interface{}) error {
	var err error
	switch proxyData := data.(type) {
	case []string:
		view := entity.ProxyXMLClassicView{
			XMLName: xml.Name{Local: "proxies"},
			Proxies: make([]string, len(proxyData)),
		}
		copy(view.Proxies, proxyData)
		err = xml.NewEncoder(writer).Encode(view)
	case []entity.Proxy:
		view := entity.ProxyXMLAdvancedView{
			XMLName: xml.Name{Local: "proxies"},
			Proxies: make([]entity.Proxy, len(proxyData)),
		}
		copy(view.Proxies, proxyData)
		err = xml.NewEncoder(writer).Encode(view)
	case []entity.AdvancedProxy:
		view := entity.ProxyXMLAllAdvancedView{
			XMLName: xml.Name{Local: "Proxies"},
			Proxies: make([]entity.AdvancedProxy, len(proxyData)),
		}
		copy(view.Proxies, proxyData)
		err = xml.NewEncoder(writer).Encode(view)
	}

	if err != nil {
		return fmt.Errorf("error encoding XML: %v", err)
	}
	return nil
}

func (r *FileRepository) EncodeYAML(writer io.Writer, data interface{}) error {
	var err error
	switch proxyData := data.(type) {
	case []string:
		view := struct {
			Proxies []string `yaml:"proxies"`
		}{
			Proxies: make([]string, len(proxyData)),
		}
		copy(view.Proxies, proxyData)
		err = yaml.NewEncoder(writer).Encode(view)
	case []entity.Proxy:
		view := struct {
			Proxies []entity.Proxy `yaml:"proxies"`
		}{
			Proxies: make([]entity.Proxy, len(proxyData)),
		}
		copy(view.Proxies, proxyData)
		err = yaml.NewEncoder(writer).Encode(view)
	case []entity.AdvancedProxy:
		view := struct {
			Proxies []entity.AdvancedProxy `yaml:"proxies"`
		}{
			Proxies: make([]entity.AdvancedProxy, len(proxyData)),
		}
		copy(view.Proxies, proxyData)
		err = yaml.NewEncoder(writer).Encode(view)
	}

	if err != nil {
		return fmt.Errorf("error encoding YAML: %v", err)
	}
	return nil
}
