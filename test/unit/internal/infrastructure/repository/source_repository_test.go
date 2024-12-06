package repository_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/fyvri/fresh-proxy-list/internal/entity"
	"github.com/fyvri/fresh-proxy-list/internal/infrastructure/repository"
)

func TestNewSourceRepository(t *testing.T) {
	type args struct {
		proxy_resources string
	}

	tests := []struct {
		name string
		args args
		want repository.SourceRepositoryInterface
	}{
		{
			name: "Success",
			args: args{
				proxy_resources: "",
			},
			want: &repository.SourceRepository{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceRepository := repository.NewSourceRepository(tt.args.proxy_resources)

			if sourceRepository == nil {
				t.Errorf(expectedReturnNonNil, "NewSourceRepository", "SourceRepositoryInterface")
			}

			got, ok := sourceRepository.(*repository.SourceRepository)
			if !ok {
				t.Errorf(expectedTypeAssertionErrorMessage, "*repository.SourceRepository")
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf(expectedButGotMessage, "*repository.SourceRepository", tt.want, got)
			}
		})
	}
}

func TestLoadSources(t *testing.T) {
	type args struct {
		proxy_resources string
	}

	tests := []struct {
		name    string
		args    args
		want    []entity.Source
		wantErr error
	}{
		{
			name: "EmptyResources",
			args: args{
				proxy_resources: "",
			},
			want:    nil,
			wantErr: errors.New("PROXY_RESOURCES not found on environment"),
		},
		{
			name: "InvalidJSON",
			args: args{
				proxy_resources: `{"invalid": "json"`,
			},
			want:    nil,
			wantErr: errors.New("error parsing JSON: unexpected end of JSON input"),
		},
		{
			name: "ValidJSON",
			args: args{
				proxy_resources: `[{"method": "GET", "category": "general", "url": "http://example.com", "is_checked": true}]`,
			},
			want: []entity.Source{
				{
					Method:    "GET",
					Category:  "general",
					URL:       "http://example.com",
					IsChecked: true,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository.SourceRepository{
				ProxyResources: tt.args.proxy_resources,
			}
			got, err := r.LoadSources()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(expectedButGotMessage, "LoadSources()", tt.want, got)
			}

			if (err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error()) ||
				(err != nil && tt.wantErr == nil) ||
				(err == nil && tt.wantErr != nil) {
				t.Errorf(expectedErrorButGotMessage, "LoadSources()", tt.wantErr, err)
			}
		})
	}
}
