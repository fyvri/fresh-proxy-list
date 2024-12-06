package repository

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fyvri/fresh-proxy-list/internal/entity"
)

type SourceRepository struct {
	ProxyResources string
}

type SourceRepositoryInterface interface {
	LoadSources() ([]entity.Source, error)
}

func NewSourceRepository(proxyResources string) SourceRepositoryInterface {
	return &SourceRepository{
		ProxyResources: proxyResources,
	}
}

func (r *SourceRepository) LoadSources() ([]entity.Source, error) {
	sourcesJSON := r.ProxyResources
	if sourcesJSON == "" {
		return nil, errors.New("PROXY_RESOURCES not found on environment")
	}

	var sources []entity.Source
	err := json.Unmarshal([]byte(sourcesJSON), &sources)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return sources, nil
}
