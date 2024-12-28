package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/fyvri/fresh-proxy-list/internal/entity"
	"github.com/fyvri/fresh-proxy-list/pkg/utils"
)

var (
	expectedButGotMessage             = "Expected %v = %v, but got = %v"
	expectedErrorButGotMessage        = "Expected %v error = %v, but got = %v"
	expectedNonEmptyMessage           = "Expected non-empty %v from %v"
	expectedReturnNonNil              = "Expected %v to return a non-nil %v"
	expectedTypeAssertionErrorMessage = "Expected type assertion error, but got = %v"
	testIP                            = "13.37.0.1"
	testPort                          = "8080"
	testProxy                         = testIP + ":" + testPort
	testHTTPCategory                  = "HTTP"
	testHTTPSCategory                 = "HTTPS"
	testSOCKS4Category                = "SOCKS4"
	testHTTPTestingSites              = []string{"http://test1.com", "http://test2.com"}
	testHTTPSTestingSites             = []string{"https://secure1.com", "https://secure2.com"}
	testUserAgents                    = []string{"Mozilla", "Chrome", "Safari"}
)

type mockURLParserUtil struct {
	ParseFunc func(urlStr string) (*url.URL, error)
}

func (m *mockURLParserUtil) Parse(urlStr string) (*url.URL, error) {
	if m.ParseFunc != nil {
		return m.ParseFunc(urlStr)
	}
	return url.Parse(urlStr)
}

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

func TestNewProxyService(t *testing.T) {
	proxyService := NewProxyService(&mockFetcherUtil{}, &mockURLParserUtil{}, testHTTPTestingSites, testHTTPSTestingSites, testUserAgents)
	if proxyService == nil {
		t.Errorf(expectedReturnNonNil, "NewProxyService", "ProxyServiceInterface")
	}

	s, ok := proxyService.(*ProxyService)
	if !ok {
		t.Errorf(expectedTypeAssertionErrorMessage, "*ProxyService")
	}

	if !reflect.DeepEqual(s.HTTPTestingSites, testHTTPTestingSites) {
		t.Errorf(expectedButGotMessage, "HTTPTestingSites", testHTTPTestingSites, s.HTTPTestingSites)
	}

	if !reflect.DeepEqual(s.HTTPSTestingSites, testHTTPSTestingSites) {
		t.Errorf(expectedButGotMessage, "HTTPSTestingSites", testHTTPSTestingSites, s.HTTPSTestingSites)
	}

	if !reflect.DeepEqual(s.UserAgents, testUserAgents) {
		t.Errorf(expectedButGotMessage, "UserAgents", testUserAgents, s.UserAgents)
	}
}

func TestCheck(t *testing.T) {
	type fields struct {
		fetcherUtil   utils.FetcherUtilInterface
		urlParserUtil utils.URLParserUtilInterface
	}

	type args struct {
		category string
		ip       string
		port     string
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *entity.Proxy
		wantError error
	}{
		{
			name: "TestValid",
			fields: fields{
				fetcherUtil:   &mockFetcherUtil{},
				urlParserUtil: &mockURLParserUtil{},
			},
			args: args{
				category: testHTTPSCategory,
				ip:       testIP,
				port:     testPort,
			},
			want: &entity.Proxy{
				Category:  testHTTPSCategory,
				Proxy:     testProxy,
				IP:        testIP,
				Port:      testPort,
				TimeTaken: 123.45,
				CheckedAt: time.Now().Format(time.RFC3339),
			},
			wantError: nil,
		},
		{
			name: "TestErrorParseURL",
			fields: fields{
				fetcherUtil: &mockFetcherUtil{},
				urlParserUtil: &mockURLParserUtil{
					ParseFunc: func(urlStr string) (*url.URL, error) {
						return nil, errors.New("parse error")
					},
				},
			},
			args: args{
				category: testHTTPCategory,
				ip:       testIP,
				port:     testPort,
			},
			want:      nil,
			wantError: errors.New("error parsing proxy URL: parse error"),
		},
		{
			name: "TestCreatingRequest",
			fields: fields{
				fetcherUtil: &mockFetcherUtil{
					NewRequestFunc: func(method, url string, body io.Reader) (*http.Request, error) {
						return nil, errors.New("error creating request")
					},
				},
			},
			args: args{
				category: testSOCKS4Category,
				ip:       testIP,
				port:     testPort,
			},
			want:      nil,
			wantError: errors.New("error creating request: error creating request"),
		},
		{
			name:   "TestUnsupportedProxyCategory",
			fields: fields{},
			args: args{
				category: "FTP",
				ip:       testIP,
				port:     testPort,
			},
			want:      nil,
			wantError: errors.New("proxy category FTP not supported"),
		},
		{
			name: "TestRequestError",
			fields: fields{
				fetcherUtil: &mockFetcherUtil{
					DoFunc: func(client *http.Client, req *http.Request) (*http.Response, error) {
						return nil, fmt.Errorf("network error")
					},
				},
				urlParserUtil: &mockURLParserUtil{},
			},
			args: args{
				category: testHTTPCategory,
				ip:       testIP,
				port:     testPort,
			},
			want:      nil,
			wantError: errors.New("request error: network error"),
		},
		{
			name: "TestUnexpectedStatusCode",
			fields: fields{
				fetcherUtil: &mockFetcherUtil{
					DoFunc: func(client *http.Client, req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusInternalServerError,
							Body:       http.NoBody,
						}, nil
					},
				},
				urlParserUtil: &mockURLParserUtil{},
			},
			args: args{
				category: testHTTPCategory,
				ip:       testIP,
				port:     testPort,
			},
			want:      nil,
			wantError: errors.New("unexpected status code 500: Internal Server Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ProxyService{
				FetcherUtil:       tt.fields.fetcherUtil,
				URLParserUtil:     tt.fields.urlParserUtil,
				HTTPTestingSites:  testHTTPTestingSites,
				HTTPSTestingSites: testHTTPSTestingSites,
				UserAgents:        testUserAgents,
				Semaphore:         make(chan struct{}, 10),
			}
			got, err := s.Check(tt.args.category, tt.args.ip, tt.args.port)

			if (err != nil && tt.wantError != nil && err.Error() != tt.wantError.Error()) ||
				(err == nil && tt.wantError != nil) ||
				(err != nil && tt.wantError == nil) {
				t.Errorf(expectedErrorButGotMessage, "ProxyService.Check()", tt.wantError, err)
			}

			if tt.want != nil &&
				(!reflect.DeepEqual(got.Category, tt.want.Category) ||
					!reflect.DeepEqual(got.Proxy, tt.want.Proxy) ||
					!reflect.DeepEqual(got.IP, tt.want.IP) ||
					!reflect.DeepEqual(got.Port, tt.want.Port)) {
				t.Errorf(expectedButGotMessage, "ProxyService.Check()", tt.want, got)
			}
		})
	}
}

func TestGetTestingSite(t *testing.T) {
	type fields struct {
		httpTestingSites  []string
		httpsTestingSites []string
	}

	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "HTTP",
			fields: fields{
				httpTestingSites: testHTTPTestingSites,
			},
			want: testHTTPTestingSites,
		},
		{
			name: "HTTPS",
			fields: fields{
				httpsTestingSites: testHTTPSTestingSites,
			},
			want: testHTTPSTestingSites,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ProxyService{
				HTTPTestingSites:  tt.fields.httpTestingSites,
				HTTPSTestingSites: tt.fields.httpsTestingSites,
			}

			site := s.GetTestingSite(tt.name)
			if len(site) == 0 {
				t.Errorf(expectedNonEmptyMessage, "site", tt.name+" sites")
			}

			found := false
			for _, expectedSite := range tt.want {
				if expectedSite == site {
					found = true
					break
				}
			}
			if !found {
				t.Errorf(expectedButGotMessage, "site", tt.want, site)
			}
		})
	}
}

func TestGetRandomUserAgent(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "RandomUserAgent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ProxyService{
				UserAgents: testUserAgents,
			}
			site := s.GetRandomUserAgent()
			found := false
			for _, ua := range s.UserAgents {
				if ua == site {
					found = true
					break
				}
			}
			if !found {
				t.Errorf(expectedButGotMessage, "user agent", s.UserAgents, site)
			}
		})
	}
}
