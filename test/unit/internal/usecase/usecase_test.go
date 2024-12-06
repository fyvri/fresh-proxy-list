package usecase_test

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"

	"github.com/fyvri/fresh-proxy-list/internal/entity"
)

var (
	unexpectedMessage                 = "Unexpected %v: %v"
	expectedButGotMessage             = "Expected %v = %v, but got = %v"
	expectedErrorButGotMessage        = "Expected %v error = %v, but got = %v"
	expectedReturnNonNil              = "Expected %v to return a non-nil %v"
	expectedTypeAssertionErrorMessage = "Expected type assertion error, but got = %v"
	testStorageDir                    = "storage"
	testClassicDir                    = "classic"
	testAdvancedDir                   = "advanced"
	testTXTExtension                  = "txt"
	testCSVExtension                  = "csv"
	testJSONExtension                 = "json"
	testXMLExtension                  = "xml"
	testYAMLExtension                 = "yaml"
	testFileOutputExtensions          = []string{testTXTExtension, testCSVExtension}
	testHTTPCategory                  = "HTTP"
	testHTTPSCategory                 = "HTTPS"
	testSOCKS4Category                = "SOCKS4"
	testSOCKS5Category                = "SOCKS5"
	testSpecialIPs                    = []string{"1.1.1.1", "2.2.2.2"}
	testPrivateIPs                    = []net.IPNet{
		{
			IP:   net.IP{3, 3, 3, 3},
			Mask: net.CIDRMask(8, 32),
		},
		{
			IP:   net.IP{4, 4, 4, 4},
			Mask: net.CIDRMask(12, 32),
		},
		{
			IP:   net.IP{5, 5, 5, 5},
			Mask: net.CIDRMask(16, 32),
		},
	}

	testIP1          = "13.37.0.1"
	testPort1        = "1337"
	testProxy1       = testIP1 + ":" + testPort1
	testCategory1    = testHTTPCategory
	testProxyEntity1 = entity.Proxy{
		Proxy:     testProxy1,
		IP:        testIP1,
		Port:      testPort1,
		Category:  testCategory1,
		TimeTaken: 0,
		CheckedAt: time.Now().Format(time.RFC3339),
	}
	testAdvancedProxyEntity1 = entity.AdvancedProxy{
		Proxy:     testProxyEntity1.Proxy,
		IP:        testProxyEntity1.IP,
		Port:      testProxyEntity1.Port,
		TimeTaken: testProxyEntity1.TimeTaken,
		CheckedAt: testProxyEntity1.CheckedAt,
		Categories: []string{
			testCategory1,
		},
	}

	testIP2          = "13.37.0.2"
	testPort2        = "1337"
	testProxy2       = testIP2 + ":" + testPort2
	testCategory2    = testHTTPSCategory
	testProxyEntity2 = entity.Proxy{
		Proxy:     testProxy2,
		IP:        testIP2,
		Port:      testPort2,
		Category:  testCategory2,
		TimeTaken: 0,
		CheckedAt: time.Now().Format(time.RFC3339),
	}
	testAdvancedProxyEntity2 = entity.AdvancedProxy{
		Proxy:     testProxyEntity2.Proxy,
		IP:        testProxyEntity2.IP,
		Port:      testProxyEntity2.Port,
		TimeTaken: testProxyEntity2.TimeTaken,
		CheckedAt: testProxyEntity2.CheckedAt,
		Categories: []string{
			testCategory2,
		},
	}

	testIP3          = "13.37.0.3"
	testPort3        = "1337"
	testProxy3       = testIP3 + ":" + testPort3
	testCategory3    = testSOCKS4Category
	testProxyEntity3 = entity.Proxy{
		Proxy:     testProxy3,
		IP:        testIP3,
		Port:      testPort3,
		Category:  testCategory3,
		TimeTaken: 0,
		CheckedAt: time.Now().Format(time.RFC3339),
	}
	testAdvancedProxyEntity3 = entity.AdvancedProxy{
		Proxy:     testProxyEntity3.Proxy,
		IP:        testProxyEntity3.IP,
		Port:      testProxyEntity3.Port,
		TimeTaken: testProxyEntity3.TimeTaken,
		CheckedAt: testProxyEntity3.CheckedAt,
		Categories: []string{
			testCategory3,
		},
	}

	testIP4          = "13.37.0.4"
	testPort4        = "1337"
	testProxy4       = testIP4 + ":" + testPort4
	testCategory4    = testSOCKS5Category
	testProxyEntity4 = entity.Proxy{
		Proxy:     testProxy4,
		IP:        testIP4,
		Port:      testPort4,
		Category:  testCategory4,
		TimeTaken: 0,
		CheckedAt: time.Now().Format(time.RFC3339),
	}
	testAdvancedProxyEntity4 = entity.AdvancedProxy{
		Proxy:     testProxyEntity4.Proxy,
		IP:        testProxyEntity4.IP,
		Port:      testProxyEntity4.Port,
		TimeTaken: testProxyEntity4.TimeTaken,
		CheckedAt: testProxyEntity4.CheckedAt,
		Categories: []string{
			testCategory4,
		},
	}
)

type mockFetcherUtil struct {
	fetchDataByte  []byte
	fetcherError   error
	NewRequestFunc func(method, url string, body io.Reader) (*http.Request, error)
	DoFunc         func(client *http.Client, req *http.Request) (*http.Response, error)
}

func (m *mockFetcherUtil) FetchData(url string) ([]byte, error) {
	if m.fetcherError != nil {
		return nil, m.fetcherError
	}
	return m.fetchDataByte, nil
}

func (m *mockFetcherUtil) Do(client *http.Client, req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(client, req)
	}
	return httptest.NewRecorder().Result(), nil
}

func (m *mockFetcherUtil) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	if m.NewRequestFunc != nil {
		return m.NewRequestFunc(method, url, body)
	}
	return http.NewRequest(method, url, body)
}

type mockProxyService struct {
	CheckFunc              func(category string, ip string, port string) (*entity.Proxy, error)
	GetTestingSiteFunc     func(category string) string
	GetRandomUserAgentFunc func() string
}

func (m *mockProxyService) Check(category string, ip string, port string) (*entity.Proxy, error) {
	if m.CheckFunc != nil {
		return m.CheckFunc(category, ip, port)
	}
	return nil, nil
}

func (m *mockProxyService) GetTestingSite(category string) string {
	if m.GetTestingSiteFunc != nil {
		return m.GetTestingSiteFunc(category)
	}
	return ""
}

func (m *mockProxyService) GetRandomUserAgent() string {
	if m.GetRandomUserAgentFunc != nil {
		return m.GetRandomUserAgentFunc()
	}
	return ""
}

type mockSourceRepository struct {
	LoadSourcesFunc func() ([]entity.Source, error)
}

func (m *mockSourceRepository) LoadSources() ([]entity.Source, error) {
	return m.LoadSourcesFunc()
}

type mockFileRepository struct {
	SaveFileFunc        func(filename string, data interface{}, format string) error
	CreateDirectoryFunc func(filePath string) error
	WriteTxtFunc        func(writer io.Writer, data interface{}) error
	EncodeCSVFunc       func(writer io.Writer, data interface{}) error
	WriteCSVFunc        func(writer io.Writer, header []string, rows [][]string) error
	EncodeJSONFunc      func(writer io.Writer, data interface{}) error
	EncodeXMLFunc       func(writer io.Writer, data interface{}) error
	EncodeYAMLFunc      func(writer io.Writer, data interface{}) error
}

func (m *mockFileRepository) SaveFile(filename string, data interface{}, ext string) error {
	if m.SaveFileFunc != nil {
		return m.SaveFileFunc(filename, data, ext)
	}
	return nil
}

func (m *mockFileRepository) CreateDirectory(filePath string) error {
	if m.CreateDirectoryFunc != nil {
		return m.CreateDirectoryFunc(filePath)
	}
	return nil
}

func (m *mockFileRepository) WriteTxt(writer io.Writer, data interface{}) error {
	if m.WriteTxtFunc != nil {
		return m.WriteTxtFunc(writer, data)
	}
	return nil
}

func (m *mockFileRepository) EncodeCSV(writer io.Writer, data interface{}) error {
	if m.EncodeCSVFunc != nil {
		return m.EncodeCSVFunc(writer, data)
	}
	return nil
}

func (m *mockFileRepository) WriteCSV(writer io.Writer, header []string, rows [][]string) error {
	if m.WriteCSVFunc != nil {
		return m.WriteCSVFunc(writer, header, rows)
	}
	return nil
}

func (m *mockFileRepository) EncodeJSON(writer io.Writer, data interface{}) error {
	if m.EncodeJSONFunc != nil {
		return m.EncodeJSONFunc(writer, data)
	}
	return nil
}

func (m *mockFileRepository) EncodeXML(writer io.Writer, data interface{}) error {
	if m.EncodeXMLFunc != nil {
		return m.EncodeXMLFunc(writer, data)
	}
	return nil
}

func (m *mockFileRepository) EncodeYAML(writer io.Writer, data interface{}) error {
	if m.EncodeYAMLFunc != nil {
		return m.EncodeYAMLFunc(writer, data)
	}
	return nil
}

type mockProxyRepository struct {
	StoreFunc                 func(proxy *entity.Proxy)
	GetAllClassicViewFunc     func() []string
	GetHTTPClassicViewFunc    func() []string
	GetHTTPSClassicViewFunc   func() []string
	GetSOCKS4ClassicViewFunc  func() []string
	GetSOCKS5ClassicViewFunc  func() []string
	GetAllAdvancedViewFunc    func() []entity.AdvancedProxy
	GetHTTPAdvancedViewFunc   func() []entity.Proxy
	GetHTTPSAdvancedViewFunc  func() []entity.Proxy
	GetSOCKS4AdvancedViewFunc func() []entity.Proxy
	GetSOCKS5AdvancedViewFunc func() []entity.Proxy

	StoredProxies []entity.Proxy
	Mutex         sync.Mutex
}

func (m *mockProxyRepository) Store(proxy *entity.Proxy) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.StoredProxies = append(m.StoredProxies, *proxy)
}

func (m *mockProxyRepository) GetStoredProxies() []entity.Proxy {
	return m.StoredProxies
}

func (m *mockProxyRepository) GetAllClassicView() []string {
	if m.GetAllClassicViewFunc != nil {
		return m.GetAllClassicViewFunc()
	}
	return nil
}

func (m *mockProxyRepository) GetHTTPClassicView() []string {
	if m.GetHTTPClassicViewFunc != nil {
		return m.GetHTTPClassicViewFunc()
	}
	return nil
}

func (m *mockProxyRepository) GetHTTPSClassicView() []string {
	if m.GetHTTPSClassicViewFunc != nil {
		return m.GetHTTPSClassicViewFunc()
	}
	return nil
}

func (m *mockProxyRepository) GetSOCKS4ClassicView() []string {
	if m.GetSOCKS4ClassicViewFunc != nil {
		return m.GetSOCKS4ClassicViewFunc()
	}
	return nil
}

func (m *mockProxyRepository) GetSOCKS5ClassicView() []string {
	if m.GetSOCKS5ClassicViewFunc != nil {
		return m.GetSOCKS5ClassicViewFunc()
	}
	return nil
}

func (m *mockProxyRepository) GetAllAdvancedView() []entity.AdvancedProxy {
	if m.GetAllAdvancedViewFunc != nil {
		return m.GetAllAdvancedViewFunc()
	}
	return nil
}

func (m *mockProxyRepository) GetHTTPAdvancedView() []entity.Proxy {
	if m.GetHTTPAdvancedViewFunc != nil {
		return m.GetHTTPAdvancedViewFunc()
	}
	return nil
}

func (m *mockProxyRepository) GetHTTPSAdvancedView() []entity.Proxy {
	if m.GetHTTPSAdvancedViewFunc != nil {
		return m.GetHTTPSAdvancedViewFunc()
	}
	return nil
}

func (m *mockProxyRepository) GetSOCKS4AdvancedView() []entity.Proxy {
	if m.GetSOCKS4AdvancedViewFunc != nil {
		return m.GetSOCKS4AdvancedViewFunc()
	}
	return nil
}

func (m *mockProxyRepository) GetSOCKS5AdvancedView() []entity.Proxy {
	if m.GetSOCKS5AdvancedViewFunc != nil {
		return m.GetSOCKS5AdvancedViewFunc()
	}
	return nil
}
