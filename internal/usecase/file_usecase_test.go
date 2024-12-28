package usecase

import (
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/fyvri/fresh-proxy-list/internal/entity"
)

var (
	mutex sync.Mutex
)

func TestSaveFiles(t *testing.T) {
	mockFileRepository := &mockFileRepository{}
	mockProxyRepository := &mockProxyRepository{}

	mockProxyRepository.GetAllClassicViewFunc = func() []string {
		return []string{
			testProxy1,
			testProxy2,
			testProxy3,
			testProxy4,
		}
	}
	mockProxyRepository.GetAllAdvancedViewFunc = func() []entity.AdvancedProxy {
		return []entity.AdvancedProxy{
			testAdvancedProxyEntity1,
			testAdvancedProxyEntity2,
			testAdvancedProxyEntity3,
			testAdvancedProxyEntity4,
		}
	}

	mockProxyRepository.GetHTTPClassicViewFunc = func() []string {
		return []string{testProxy1}
	}
	mockProxyRepository.GetHTTPAdvancedViewFunc = func() []entity.Proxy {
		return []entity.Proxy{
			testProxyEntity1,
		}
	}

	mockProxyRepository.GetHTTPSClassicViewFunc = func() []string {
		return []string{
			testProxy2,
		}
	}
	mockProxyRepository.GetHTTPSAdvancedViewFunc = func() []entity.Proxy {
		return []entity.Proxy{
			testProxyEntity2,
		}
	}

	mockProxyRepository.GetSOCKS4ClassicViewFunc = func() []string {
		return []string{
			testProxy3,
		}
	}
	mockProxyRepository.GetSOCKS4AdvancedViewFunc = func() []entity.Proxy {
		return []entity.Proxy{
			testProxyEntity3,
		}
	}

	mockProxyRepository.GetSOCKS5ClassicViewFunc = func() []string {
		return []string{
			testProxy4,
		}
	}
	mockProxyRepository.GetSOCKS5AdvancedViewFunc = func() []entity.Proxy {
		return []entity.Proxy{
			testProxyEntity4,
		}
	}

	got := 0
	mockFileRepository.SaveFileFunc = func(filename string, data interface{}, extension string) error {
		mutex.Lock()
		defer mutex.Unlock()

		got++
		t.Logf("SaveFile called with filename: %s, extension: %s", filename, extension)
		if !strings.HasPrefix(filename, filepath.Join(testStorageDir, testClassicDir)) &&
			!strings.HasPrefix(filename, filepath.Join(testStorageDir, testAdvancedDir)) {
			t.Errorf(unexpectedMessage, "filename", filename)
		}
		if extension != testCSVExtension && extension != testJSONExtension && extension != testXMLExtension && extension != testYAMLExtension && extension != testTXTExtension {
			t.Errorf(unexpectedMessage, "extension", extension)
		}
		return nil
	}
	uc := NewFileUsecase(mockFileRepository, mockProxyRepository, testFileOutputExtensions)
	uc.SaveFiles()

	// (5 categories * number of extensions * 2 file types (classic, advanced)) + (5 all * 1 extension txt * 1 file type classic)
	want := (5 * len(testFileOutputExtensions) * 2) + (5 * 1 * 1)
	if got != want {
		t.Errorf(expectedButGotMessage, "calls", want, got)
	}
}
