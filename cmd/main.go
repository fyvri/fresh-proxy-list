package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/fyvri/fresh-proxy-list/internal/entity"
	"github.com/fyvri/fresh-proxy-list/internal/infrastructure/config"
	"github.com/fyvri/fresh-proxy-list/internal/infrastructure/repository"
	"github.com/fyvri/fresh-proxy-list/internal/service"
	"github.com/fyvri/fresh-proxy-list/internal/usecase"
	"github.com/fyvri/fresh-proxy-list/pkg/utils"

	"github.com/joho/godotenv"
)

type Runners struct {
	fetcherUtil      utils.FetcherUtilInterface
	urlParserUtil    utils.URLParserUtilInterface
	proxyService     service.ProxyServiceInterface
	sourceRepository repository.SourceRepositoryInterface
	proxyRepository  repository.ProxyRepositoryInterface
	fileRepository   repository.FileRepositoryInterface
}

func main() {
	if err := runApplication(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func runApplication() error {
	loadEnv()

	httpTestingSites := config.HTTPTestingSites
	httpsTestingSites := config.HTTPSTestingSites
	userAgents := config.UserAgents

	mkdirAll := func(path string, perm os.FileMode) error {
		return os.MkdirAll(path, perm)
	}
	create := func(name string) (io.Writer, error) {
		file, err := os.Create(name)
		if err != nil {
			return nil, err
		}
		return file, nil
	}

	fetcherUtil := utils.NewFetcher(http.DefaultClient, http.NewRequest)
	urlParserUtil := utils.NewURLParser()
	csvWriterUtil := utils.NewCSVWriter()
	proxyService := service.NewProxyService(fetcherUtil, urlParserUtil, httpTestingSites, httpsTestingSites, userAgents)
	sourceRepository := repository.NewSourceRepository(os.Getenv("PROXY_RESOURCES"))
	proxyRepository := repository.NewProxyRepository()
	fileRepository := repository.NewFileRepository(mkdirAll, create, csvWriterUtil)

	runners := Runners{
		fetcherUtil:      fetcherUtil,
		urlParserUtil:    urlParserUtil,
		proxyService:     proxyService,
		sourceRepository: sourceRepository,
		proxyRepository:  proxyRepository,
		fileRepository:   fileRepository,
	}

	return run(runners)
}

func loadEnv() error {
	return godotenv.Load()
}

func run(runners Runners) error {
	startTime := time.Now()

	sourceUsecase := usecase.NewSourceUsecase(runners.sourceRepository, runners.fetcherUtil)
	sources, err := sourceUsecase.LoadSources()
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	proxyCategories := config.ProxyCategories
	specialIPs := config.SpecialIPs
	privateIPs := config.PrivateIPs
	proxyUsecase := usecase.NewProxyUsecase(runners.proxyRepository, runners.proxyService, specialIPs, privateIPs)
	for i, source := range sources {
		if _, found := slices.BinarySearch(proxyCategories, source.Category); found {
			wg.Add(1)
			go func(source entity.Source) {
				defer wg.Done()

				innerWG := sync.WaitGroup{}
				proxies, err := sourceUsecase.ProcessSource(&source)
				if err != nil {
					return
				}

				for _, proxy := range proxies {
					innerWG.Add(1)
					go func(source entity.Source, proxy string) {
						defer innerWG.Done()
						proxyUsecase.ProcessProxy(source.Category, proxy, source.IsChecked)
					}(source, proxy)
				}
				innerWG.Wait()
			}(source)
		} else {
			log.Printf("Index %v: proxy category not found", i)
		}
	}
	wg.Wait()

	fileOutputExtensions := config.FileOutputExtensions
	fileUsecase := usecase.NewFileUsecase(runners.fileRepository, runners.proxyRepository, fileOutputExtensions)
	fileUsecase.SaveFiles()

	log.Printf("Number of proxies     : %v", len(proxyUsecase.GetAllAdvancedView()))
	log.Printf("Time-consuming process: %v", time.Since(startTime))
	return nil
}
