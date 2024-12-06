package usecase

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/fyvri/fresh-proxy-list/internal/infrastructure/repository"
)

type fileUsecase struct {
	FileRepository       repository.FileRepositoryInterface
	ProxyRepository      repository.ProxyRepositoryInterface
	FileOutputExtensions []string
	WaitGroup            sync.WaitGroup
}

type FileUsecaseInterface interface {
	SaveFiles()
}

func NewFileUsecase(fileRepository repository.FileRepositoryInterface, proxyRepository repository.ProxyRepositoryInterface, fileOutputExtensions []string) FileUsecaseInterface {
	return &fileUsecase{
		FileRepository:       fileRepository,
		ProxyRepository:      proxyRepository,
		FileOutputExtensions: fileOutputExtensions,
		WaitGroup:            sync.WaitGroup{},
	}
}

func (uc *fileUsecase) SaveFiles() {
	createFile := func(filename string, classic []string, advanced interface{}) {
		uc.WaitGroup.Add((len(uc.FileOutputExtensions) * 2) + 1)

		filename = strings.ToLower(filename)
		for _, ext := range uc.FileOutputExtensions {
			go func(ext string) {
				defer uc.WaitGroup.Done()
				uc.FileRepository.SaveFile(filepath.Join("storage", "classic", filename+"."+ext), classic, ext)
			}(ext)
			go func(ext string) {
				defer uc.WaitGroup.Done()
				uc.FileRepository.SaveFile(filepath.Join("storage", "advanced", filename+"."+ext), advanced, ext)
			}(ext)
		}

		go func() {
			defer uc.WaitGroup.Done()
			uc.FileRepository.SaveFile(filepath.Join("storage", "classic", filename+".txt"), classic, "txt")
		}()
	}

	createFile("all", uc.ProxyRepository.GetAllClassicView(), uc.ProxyRepository.GetAllAdvancedView())
	createFile("http", uc.ProxyRepository.GetHTTPClassicView(), uc.ProxyRepository.GetHTTPAdvancedView())
	createFile("https", uc.ProxyRepository.GetHTTPSClassicView(), uc.ProxyRepository.GetHTTPSAdvancedView())
	createFile("socks4", uc.ProxyRepository.GetSOCKS4ClassicView(), uc.ProxyRepository.GetSOCKS4AdvancedView())
	createFile("socks5", uc.ProxyRepository.GetSOCKS5ClassicView(), uc.ProxyRepository.GetSOCKS5AdvancedView())
	uc.WaitGroup.Wait()
}
