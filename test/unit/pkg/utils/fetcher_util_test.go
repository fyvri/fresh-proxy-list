package util_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/fyvri/fresh-proxy-list/pkg/utils"
)

var (
	testGETMethod  = "GET"
	testPOSTMethod = "POST"
	testClient     = http.DefaultClient
	testNewRequest = http.NewRequest
)

func TestNewFetcher(t *testing.T) {
	fetcherUtil := utils.NewFetcher(testClient, testNewRequest)

	if fetcherUtil == nil {
		t.Errorf(expectedReturnNonNil, "NewFetcher", "FetcherInterface")
	}

	fetcherUtilInstance, ok := fetcherUtil.(*utils.FetcherUtil)
	if !ok {
		t.Errorf(expectedTypeAssertionErrorMessage, "*utils.FetcherUtil")
	}

	req, err := fetcherUtilInstance.NewRequestFunc(testGETMethod, testRawURL, nil)
	if err != nil {
		t.Errorf(expectedButGotMessage, "newRequest", "no error", err)
	}

	if req.Method != testGETMethod {
		t.Errorf(expectedButGotMessage, "method", testGETMethod, req.Method)
	}

	if req.URL.String() != testRawURL {
		t.Errorf(expectedButGotMessage, "URL", testRawURL, req.URL.String())
	}

	if fetcherUtilInstance.Client != nil && fetcherUtilInstance.Client != http.DefaultClient {
		t.Errorf(expectedButGotMessage, "client", http.DefaultClient, fetcherUtilInstance.Client)
	}
}

func TestFetchData(t *testing.T) {
	type fields struct {
		transport      *mockTransport
		newRequestFunc func(method string, url string, body io.Reader) (*http.Request, error)
	}

	type args struct {
		url string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr error
	}{
		{
			name: "Success",
			fields: fields{
				transport: &mockTransport{
					response: &http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(&mockReadCloser{
							data: []byte("response data"),
						}),
					},
					err: nil,
				},
				newRequestFunc: testNewRequest,
			},
			args: args{
				url: testRawURL,
			},
			want:    []byte("response data"),
			wantErr: nil,
		},
		{
			name: "NewRequestError",
			fields: fields{
				transport: &mockTransport{
					response: nil,
					err:      nil,
				},
				newRequestFunc: func(method string, url string, body io.Reader) (*http.Request, error) {
					return nil, fmt.Errorf("new request error")
				},
			},
			args: args{
				url: testRawURL,
			},
			want:    nil,
			wantErr: errors.New("new request error"),
		},
		{
			name: "RequestError",
			fields: fields{
				transport: &mockTransport{
					response: nil,
					err:      fmt.Errorf("request error"),
				},
				newRequestFunc: testNewRequest,
			},
			args: args{
				url: testRawURL,
			},
			want:    nil,
			wantErr: fmt.Errorf("Get \"%s\": request error", testRawURL),
		},
		{
			name: "ResponseError",
			fields: fields{
				transport: &mockTransport{
					response: &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body: io.NopCloser(&mockReadCloser{
							data: []byte("error response"),
						}),
					},
					err: nil,
				},
				newRequestFunc: testNewRequest,
			},
			args: args{
				url: testRawURL,
			},
			want:    []byte("error response"),
			wantErr: errors.New("failed to fetch data: Internal Server Error"),
		},
		{
			name: "ReadBodyError",
			fields: fields{
				transport: &mockTransport{
					response: &http.Response{
						StatusCode: http.StatusOK,
						Body: &mockReadCloser{
							errRead: fmt.Errorf("body read error"),
						},
					},
					err: nil,
				},
				newRequestFunc: testNewRequest,
			},
			args: args{
				url: testRawURL,
			},
			want:    nil,
			wantErr: errors.New("body read error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetcherUtil := &utils.FetcherUtil{
				Client: &http.Client{
					Transport: tt.fields.transport,
				},
				NewRequestFunc: tt.fields.newRequestFunc,
			}
			got, err := fetcherUtil.FetchData(tt.args.url)

			if (err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error()) ||
				(err == nil && tt.wantErr != nil) ||
				(err != nil && tt.wantErr == nil) {
				t.Errorf(expectedErrorButGotMessage, "FetchData()", tt.wantErr, err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(expectedButGotMessage, "FetchData()", tt.want, got)
			}
		})
	}
}

func TestNewRequest(t *testing.T) {
	type args struct {
		method string
		url    string
		body   io.Reader
	}

	type want struct {
		url    string
		method string
	}

	tests := []struct {
		name    string
		args    args
		want    want
		wantErr error
	}{
		{
			name: testGETMethod,
			args: args{
				method: http.MethodGet,
				url:    testRawURL,
				body:   nil,
			},
			want: want{
				url:    testRawURL,
				method: http.MethodGet,
			},
			wantErr: nil,
		},
		{
			name: testPOSTMethod,
			args: args{
				method: http.MethodPost,
				url:    testRawURL,
				body:   bytes.NewReader([]byte("body")),
			},
			want: want{
				url:    testRawURL,
				method: http.MethodPost,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := utils.NewFetcher(&http.Client{
				Transport: &mockTransport{},
			}, testNewRequest)
			req, err := u.NewRequest(tt.args.method, tt.args.url, tt.args.body)

			if err != nil {
				t.Errorf(expectedErrorButGotMessage, "NewRequest()", nil, err)
			}

			if req.Method != tt.want.method {
				t.Errorf(expectedButGotMessage, "method", tt.want.method, req.Method)
			}

			if req.URL.String() != tt.want.url {
				t.Errorf(expectedButGotMessage, "URL", tt.want.url, req.URL.String())
			}
		})
	}
}
