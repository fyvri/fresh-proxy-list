package usecase

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fyvri/fresh-proxy-list/internal/entity"
	"github.com/fyvri/fresh-proxy-list/internal/infrastructure/repository"
	"github.com/fyvri/fresh-proxy-list/pkg/utils"
)

type SourceUsecase struct {
	SourceRepository repository.SourceRepositoryInterface
	FetcherUtil      utils.FetcherUtilInterface
}

type SourceUsecaseInterface interface {
	LoadSources() ([]entity.Source, error)
	ProcessSource(source *entity.Source) ([]string, error)
}

func NewSourceUsecase(sourceRepository repository.SourceRepositoryInterface, fetcherUtil utils.FetcherUtilInterface) SourceUsecaseInterface {
	return &SourceUsecase{
		SourceRepository: sourceRepository,
		FetcherUtil:      fetcherUtil,
	}
}

func (uc *SourceUsecase) LoadSources() ([]entity.Source, error) {
	return uc.SourceRepository.LoadSources()
}

func (uc *SourceUsecase) ProcessSource(source *entity.Source) ([]string, error) {
	body, err := uc.FetcherUtil.FetchData(source.URL)
	if err != nil {
		return nil, err
	}

	var proxies []string
	switch source.Method {
	case "LIST":
		proxies = strings.Split(strings.TrimSpace(string(body)), "\n")
	case "SCRAP":
		re := regexp.MustCompile(`[0-9]+(?:\.[0-9]+){3}:[0-9]+`)
		proxies = re.FindAllString(string(body), -1)
	default:
		return nil, fmt.Errorf("source method not found: %s", source.Method)
	}

	return proxies, nil
}
