package usecase

import (
	"errors"
	"reflect"
	"testing"

	"github.com/fyvri/fresh-proxy-list/internal/entity"
	"github.com/fyvri/fresh-proxy-list/internal/infrastructure/repository"
	"github.com/fyvri/fresh-proxy-list/pkg/utils"
)

var (
	testListMethod  = "LIST"
	testScrapMethod = "SCRAP"
	testCategory    = "HTTP"
	testURL         = "http://example.com"
	testIsChecked   = true
)

func TestNewSourceUsecase(t *testing.T) {
	type fields struct {
		sourceRepository repository.SourceRepositoryInterface
		fetcherUtil      utils.FetcherUtilInterface
	}

	tests := []struct {
		name   string
		fields fields
		want   *SourceUsecase
	}{
		{
			name: "Success",
			fields: fields{
				sourceRepository: &mockSourceRepository{},
				fetcherUtil:      &mockFetcherUtil{},
			},
			want: &SourceUsecase{
				SourceRepository: &mockSourceRepository{},
				FetcherUtil:      &mockFetcherUtil{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceUsecase := NewSourceUsecase(tt.fields.sourceRepository, tt.fields.fetcherUtil)
			if sourceUsecase == nil {
				t.Errorf(expectedReturnNonNil, "NewSourceUsecase", "SourceUsecaseInterface")
			}

			got, ok := sourceUsecase.(*SourceUsecase)
			if !ok {
				t.Errorf(expectedTypeAssertionErrorMessage, "*SourceUsecase")
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf(expectedButGotMessage, "*SourceUsecase", tt.want, got)
			}
		})
	}
}

func TestLoadSourcesSuccess(t *testing.T) {
	type fields struct {
		sourceRepository repository.SourceRepositoryInterface
	}

	tests := []struct {
		name      string
		fields    fields
		want      []entity.Source
		wantError error
	}{
		{
			name: "Success",
			fields: fields{
				sourceRepository: &mockSourceRepository{
					LoadSourcesFunc: func() ([]entity.Source, error) {
						return []entity.Source{
							{
								Method:    testListMethod,
								Category:  testCategory,
								URL:       testURL,
								IsChecked: testIsChecked,
							},
						}, nil
					},
				},
			},
			want: []entity.Source{
				{
					Method:    testListMethod,
					Category:  testCategory,
					URL:       testURL,
					IsChecked: testIsChecked,
				},
			},
			wantError: nil,
		},
		{
			name: "Error",
			fields: fields{
				sourceRepository: &mockSourceRepository{
					LoadSourcesFunc: func() ([]entity.Source, error) {
						return nil, errors.New("load proxy resource error")
					},
				},
			},
			want:      nil,
			wantError: errors.New("load proxy resource error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &SourceUsecase{
				SourceRepository: tt.fields.sourceRepository,
			}
			got, err := uc.LoadSources()

			if err != nil && err.Error() != tt.wantError.Error() {
				t.Errorf(expectedErrorButGotMessage, "ProcessProxy()", tt.wantError, err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(expectedButGotMessage, "ProcessProxy()", tt.want, got)
			}
		})
	}
}

func TestProcessSource(t *testing.T) {
	type fields struct {
		fetcherUtil utils.FetcherUtilInterface
	}

	type args struct {
		source entity.Source
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []string
		wantError error
	}{
		{
			name: "TestFetcherError",
			fields: fields{
				fetcherUtil: &mockFetcherUtil{
					fetcherError: errors.New("error creating request"),
				},
			},
			args: args{
				source: entity.Source{
					Method:    testListMethod,
					Category:  testCategory,
					URL:       testURL,
					IsChecked: testIsChecked,
				},
			},
			want:      nil,
			wantError: errors.New("error creating request"),
		},
		{
			name: "TestFetcherWithListMethod",
			fields: fields{
				fetcherUtil: &mockFetcherUtil{
					fetchDataByte: []byte(testProxy1 + "\n" + testProxy2 + "\n" + testProxy3 + "\n" + testProxy4),
				},
			},
			args: args{
				source: entity.Source{
					Method:    testListMethod,
					Category:  testCategory,
					URL:       testURL,
					IsChecked: testIsChecked,
				},
			},
			want: []string{
				testProxy1,
				testProxy2,
				testProxy3,
				testProxy4,
			},
			wantError: nil,
		},
		{
			name: "TestFetcherWithScrapMethod",
			fields: fields{
				fetcherUtil: &mockFetcherUtil{
					fetchDataByte: []byte(testProxy1 + "\n" + testProxy2 + "\n" + testProxy3 + "\n" + testProxy4),
				},
			},
			args: args{
				source: entity.Source{
					Method:    testScrapMethod,
					Category:  testCategory,
					URL:       testURL,
					IsChecked: testIsChecked,
				},
			},
			want: []string{
				testProxy1,
				testProxy2,
				testProxy3,
				testProxy4,
			},
			wantError: nil,
		},
		{
			name: "TestFetcherWithUndefinedMethod",
			fields: fields{
				fetcherUtil: &mockFetcherUtil{},
			},
			args: args{
				source: entity.Source{
					Method:    "NO_METHOD",
					Category:  testCategory,
					URL:       testURL,
					IsChecked: testIsChecked,
				},
			},
			want:      nil,
			wantError: errors.New("source method not found: NO_METHOD"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &SourceUsecase{
				FetcherUtil: tt.fields.fetcherUtil,
			}
			got, err := uc.ProcessSource(&tt.args.source)

			if (err != nil && tt.wantError != nil && err.Error() != tt.wantError.Error()) ||
				(err == nil && tt.wantError != nil) ||
				(err != nil && tt.wantError == nil) {
				t.Errorf(expectedErrorButGotMessage, "SourceUsecase.ProcessSource()", tt.wantError, err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(expectedButGotMessage, "SourceUsecase.ProcessSource()", tt.want, got)
			}
		})
	}
}
