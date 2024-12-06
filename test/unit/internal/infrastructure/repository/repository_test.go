package repository_test

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"time"

	"github.com/fyvri/fresh-proxy-list/internal/entity"
)

var (
	expectedButGotMessage             = "Expected %v = %v, but got = %v"
	expectedErrorButGotMessage        = "Expected %v error = %v, but got = %v"
	expectedTypeAssertionErrorMessage = "Expected type assertion error, but got = %v"
	expectedReturnNonNil              = "Expected %v to return a non-nil %v"
	testErrorWriting                  = "error writing"
	testErrorEncode                   = "error encoding %s: %s"
	testStorageDir                    = "/tmp"
	testClassicDir                    = "/classic"
	testAdvancedDir                   = "/advanced"
	testClassicFilePath               = testStorageDir + testClassicDir + "/test_file"
	testAdvancedFilePath              = testStorageDir + testAdvancedDir + "/test_file"
	testTXTExtension                  = "txt"
	testCSVExtension                  = "csv"
	testJSONExtension                 = "json"
	testXMLExtension                  = "xml"
	testYAMLExtension                 = "yaml"
	testHTTPCategory                  = "HTTP"
	testHTTPSCategory                 = "HTTPS"
	testSOCKS4Category                = "SOCKS4"
	testSOCKS5Category                = "SOCKS5"
	testTimeTaken                     = 1.2345
	testCheckedAt                     = "2024-07-27T00:00:00Z"

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

	testIPs = []string{
		testIP1,
		testIP2,
		testIP3,
		testIP4,
	}
	testProxies = []entity.Proxy{
		testProxyEntity1,
		testProxyEntity2,
		testProxyEntity3,
		testProxyEntity4,
	}
	testAdvancedProxies = []entity.AdvancedProxy{
		testAdvancedProxyEntity1,
		testAdvancedProxyEntity2,
		testAdvancedProxyEntity3,
		testAdvancedProxyEntity4,
	}
	testIPsToString, _             = json.Marshal(testIPs)
	testProxiesToString, _         = json.Marshal(testProxies)
	testAdvancedProxiesToString, _ = json.Marshal(testAdvancedProxies)
)

type mockWriter struct {
	errWrite error
	errClose error
}

func (m *mockWriter) Write(p []byte) (int, error) {
	if m.errWrite != nil {
		return 0, m.errWrite
	}
	return len(p), nil
}

func (m *mockWriter) Close() error {
	if m.errClose != nil {
		return m.errClose
	}
	return nil
}

type mockCSVWriterUtil struct {
	errFlush error
	errWrite error
}

func (m *mockCSVWriterUtil) Init(w io.Writer) *csv.Writer {
	return csv.NewWriter(w)
}

func (m *mockCSVWriterUtil) Flush(csvWriter *csv.Writer) {
	if m.errFlush != nil {
		return
	}
}

func (m *mockCSVWriterUtil) Write(csvWriter *csv.Writer, record []string) error {
	if m.errWrite != nil {
		return m.errWrite
	}
	return nil
}
