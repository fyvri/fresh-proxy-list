package usecase_test

import (
	"errors"
	"net"
	"reflect"
	"sync"
	"testing"

	"github.com/fyvri/fresh-proxy-list/internal/entity"
	"github.com/fyvri/fresh-proxy-list/internal/infrastructure/repository"
	"github.com/fyvri/fresh-proxy-list/internal/service"
	"github.com/fyvri/fresh-proxy-list/internal/usecase"
)

func TestNewProxyUsecase(t *testing.T) {
	mockProxyRepository := &mockProxyRepository{}
	mockProxyService := &mockProxyService{}
	proxyUsecase := usecase.NewProxyUsecase(mockProxyRepository, mockProxyService, testSpecialIPs, testPrivateIPs)
	if proxyUsecase == nil {
		t.Errorf(expectedReturnNonNil, "NewProxyUsecase", "ProxyUsecaseInterface")
	}

	uc, ok := proxyUsecase.(*usecase.ProxyUsecase)
	if !ok {
		t.Errorf(expectedTypeAssertionErrorMessage, "*usecase.ProxyUsecase")
	}

	testKey := testProxyEntity1.Category + "_" + testProxyEntity1.Proxy
	testValue := true
	uc.ProxyMap.Store(testKey, testValue)

	got, ok := uc.ProxyMap.Load(testKey)
	if !ok || got != testValue {
		t.Errorf(expectedButGotMessage, "value", testValue, got)
	}

	_, loaded := uc.ProxyMap.LoadOrStore(testKey, false)
	if !loaded {
		t.Errorf("Expected LoadOrStore to return true indicating the key was loaded")
	}

	got, _ = uc.ProxyMap.Load(testKey)
	if got != testValue {
		t.Errorf(expectedButGotMessage, "value after LoadOrStore", testValue, got)
	}

	if !reflect.DeepEqual(uc.SpecialIPs, testSpecialIPs) {
		t.Errorf(expectedButGotMessage, "SpecialIPs", testSpecialIPs, uc.SpecialIPs)
	}

	if !reflect.DeepEqual(uc.PrivateIPs, testPrivateIPs) {
		t.Errorf(expectedButGotMessage, "PrivateIPs", testPrivateIPs, uc.PrivateIPs)
	}
}

func TestProcessProxy(t *testing.T) {
	type fields struct {
		proxyRepository repository.ProxyRepositoryInterface
		proxyService    service.ProxyServiceInterface
	}

	type args struct {
		category  string
		proxy     string
		isChecked bool
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *entity.Proxy
		wantError error
	}{
		{
			name: "ProxyNotFound",
			fields: fields{
				proxyRepository: &mockProxyRepository{},
				proxyService:    &mockProxyService{},
			},
			args: args{
				category:  testHTTPCategory,
				proxy:     "   ",
				isChecked: false,
			},
			want:      nil,
			wantError: errors.New("proxy not found"),
		},
		{
			name: "ProxyFormatIncorrect",
			fields: fields{
				proxyRepository: &mockProxyRepository{},
				proxyService:    &mockProxyService{},
			},
			args: args{
				category:  testHTTPCategory,
				proxy:     "invalid-proxy",
				isChecked: false,
			},
			want:      nil,
			wantError: errors.New("proxy format incorrect"),
		},
		{
			name: "ProxyFormatNotMatch",
			fields: fields{
				proxyRepository: &mockProxyRepository{},
				proxyService:    &mockProxyService{},
			},
			args: args{
				category:  testHTTPCategory,
				proxy:     "invalid-proxy:1337",
				isChecked: false,
			},
			want:      nil,
			wantError: errors.New("proxy format not match"),
		},
		{
			name: "ProxyIsSpecialIP",
			fields: fields{
				proxyRepository: &mockProxyRepository{},
				proxyService:    &mockProxyService{},
			},
			args: args{
				category:  testHTTPCategory,
				proxy:     "1.1.1.1:1337",
				isChecked: false,
			},
			want:      nil,
			wantError: errors.New("proxy belongs to special ip"),
		},
		{
			name: "ProxyPortIsMoreThan65535",
			fields: fields{
				proxyRepository: &mockProxyRepository{},
				proxyService:    &mockProxyService{},
			},
			args: args{
				category:  testHTTPCategory,
				proxy:     testIP1 + ":65540",
				isChecked: false,
			},
			want:      nil,
			wantError: errors.New("proxy port format incorrect"),
		},
		{
			name: "ProxyHasBeenProcessed",
			fields: fields{
				proxyRepository: &mockProxyRepository{},
			},
			args: args{
				category:  testProxyEntity1.Category,
				proxy:     testProxyEntity1.Proxy,
				isChecked: false,
			},
			want:      nil,
			wantError: errors.New("proxy has been processed"),
		},
		{
			name: "ValidProxy",
			fields: fields{
				proxyRepository: &mockProxyRepository{},
				proxyService: &mockProxyService{
					CheckFunc: func(category string, ip string, port string) (*entity.Proxy, error) {
						return &testProxyEntity1, nil
					},
				},
			},
			args: args{
				category:  testProxyEntity1.Category,
				proxy:     testProxyEntity1.Proxy,
				isChecked: true,
			},
			want: &entity.Proxy{
				Category:  testProxyEntity1.Category,
				Proxy:     testProxyEntity1.Proxy,
				IP:        testProxyEntity1.IP,
				Port:      testProxyEntity1.Port,
				TimeTaken: testProxyEntity1.TimeTaken,
				CheckedAt: testProxyEntity1.CheckedAt,
			},
			wantError: nil,
		},
		{
			name: "NotValidProxy",
			fields: fields{
				proxyRepository: &mockProxyRepository{},
				proxyService: &mockProxyService{
					CheckFunc: func(category string, ip string, port string) (*entity.Proxy, error) {
						return nil, errors.New("proxy not valid")
					},
				},
			},
			args: args{
				category:  testProxyEntity1.Category,
				proxy:     testProxyEntity1.Proxy,
				isChecked: true,
			},
			want:      nil,
			wantError: errors.New("proxy not valid"),
		},
		{
			name: "ValidProxyWithNotChecked",
			fields: fields{
				proxyRepository: &mockProxyRepository{},
				proxyService: &mockProxyService{
					CheckFunc: func(category string, ip string, port string) (*entity.Proxy, error) {
						return &testProxyEntity1, nil
					},
				},
			},
			args: args{
				category:  testProxyEntity1.Category,
				proxy:     testProxyEntity1.Proxy,
				isChecked: false,
			},
			want: &entity.Proxy{
				Category:  testProxyEntity1.Category,
				Proxy:     testProxyEntity1.Proxy,
				IP:        testProxyEntity1.IP,
				Port:      testProxyEntity1.Port,
				TimeTaken: 0,
				CheckedAt: "",
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &usecase.ProxyUsecase{
				ProxyRepository: tt.fields.proxyRepository,
				ProxyService:    tt.fields.proxyService,
				ProxyMap:        sync.Map{},
				SpecialIPs:      testSpecialIPs,
				PrivateIPs:      testPrivateIPs,
			}

			if tt.name == "ProxyHasBeenProcessed" {
				uc.ProxyMap.Store(tt.args.category+"_"+tt.args.proxy, true)
			}

			got, err := uc.ProcessProxy(tt.args.category, tt.args.proxy, tt.args.isChecked)

			if err != nil && err.Error() != tt.wantError.Error() {
				t.Errorf(expectedErrorButGotMessage, "ProcessProxy()", tt.wantError, err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(expectedButGotMessage, "ProcessProxy()", tt.want, got)
			}
		})
	}
}

func TestIsSpecialIP(t *testing.T) {
	type args struct {
		ip string
	}

	type fields struct {
		specialIPs []string
		privateIPs []net.IPNet
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "ItIsSpecialIP",
			fields: fields{
				specialIPs: testSpecialIPs,
				privateIPs: testPrivateIPs,
			},
			args: args{
				ip: "1.1.1.1",
			},
			want: true,
		},
		{
			name: "ErrorParseIP",
			fields: fields{
				specialIPs: testSpecialIPs,
				privateIPs: testPrivateIPs,
			},
			args: args{
				ip: "13.37.1",
			},
			want: true,
		},
		{
			name: "ItIsUnspecified",
			fields: fields{
				specialIPs: testSpecialIPs,
				privateIPs: testPrivateIPs,
			},
			args: args{
				ip: "::1",
			},
			want: true,
		},
		{
			name: "ItIsPrivateIP",
			fields: fields{
				specialIPs: testSpecialIPs,
				privateIPs: testPrivateIPs,
			},
			args: args{
				ip: "5.5.5.5",
			},
			want: true,
		},
		{
			name: "ItIsNotSpecialIP",
			fields: fields{
				specialIPs: testSpecialIPs,
				privateIPs: testPrivateIPs,
			},
			args: args{
				ip: "13.37.0.1",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &usecase.ProxyUsecase{
				SpecialIPs: testSpecialIPs,
				PrivateIPs: testPrivateIPs,
			}
			got := uc.IsSpecialIP(tt.args.ip)
			if got != tt.want {
				t.Errorf(expectedButGotMessage, "IP: "+tt.args.ip, tt.want, got)
			}
		})
	}
}

func TestGetAllAdvancedView(t *testing.T) {
	type fields struct {
		proxyRepository repository.ProxyRepositoryInterface
	}

	tests := []struct {
		name   string
		fields fields
		want   []entity.AdvancedProxy
	}{
		{
			name: "Should return all advanced view proxies",
			fields: fields{
				proxyRepository: &mockProxyRepository{
					GetAllAdvancedViewFunc: func() []entity.AdvancedProxy {
						return []entity.AdvancedProxy{
							testAdvancedProxyEntity1,
						}
					},
				},
			},
			want: []entity.AdvancedProxy{
				testAdvancedProxyEntity1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &usecase.ProxyUsecase{
				ProxyRepository: tt.fields.proxyRepository,
			}
			got := uc.GetAllAdvancedView()
			if len(got) != len(tt.want) {
				t.Errorf(expectedButGotMessage, "GetAllAdvancedView()", tt.want, got)
			}

			for i, v := range got {
				if !reflect.DeepEqual(v, tt.want[i]) {
					t.Errorf(expectedButGotMessage, "GetAllAdvancedView()", tt.want[i], v)
				}
			}
		})
	}
}
