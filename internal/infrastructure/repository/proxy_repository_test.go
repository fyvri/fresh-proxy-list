package repository

import (
	"reflect"
	"testing"

	"github.com/fyvri/fresh-proxy-list/internal/entity"
)

func TestNewProxyRepository(t *testing.T) {
	tests := []struct {
		name string
		want ProxyRepositoryInterface
	}{
		{
			name: "Success",
			want: &ProxyRepository{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proxyRepository := NewProxyRepository()

			if proxyRepository == nil {
				t.Errorf(expectedReturnNonNil, "NewProxyRepository", "ProxyRepositoryInterface")
			}

			got, ok := proxyRepository.(*ProxyRepository)
			if !ok {
				t.Errorf(expectedTypeAssertionErrorMessage, "*ProxyRepository")
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf(expectedButGotMessage, "*ProxyRepository", tt.want, got)
			}
		})
	}
}

func TestProxyRepository(t *testing.T) {
	type fields struct {
		allClassicView     []string
		httpClassicView    []string
		httpsClassicView   []string
		socks4ClassicView  []string
		socks5ClassicView  []string
		allAdvancedView    []entity.AdvancedProxy
		httpAdvancedView   []entity.Proxy
		httpsAdvancedView  []entity.Proxy
		socks4AdvancedView []entity.Proxy
		socks5AdvancedView []entity.Proxy
	}

	type args struct {
		proxy entity.Proxy
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    fields
		wantErr error
	}{
		{
			name: "StoreHTTPProxy",
			fields: fields{
				allClassicView:     []string{},
				httpClassicView:    []string{},
				httpsClassicView:   []string{},
				socks4ClassicView:  []string{},
				socks5ClassicView:  []string{},
				allAdvancedView:    []entity.AdvancedProxy{},
				httpAdvancedView:   []entity.Proxy{},
				httpsAdvancedView:  []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{},
				socks5AdvancedView: []entity.Proxy{},
			},
			args: args{
				proxy: testProxyEntity1,
			},
			want: fields{
				allClassicView: []string{
					testProxy1,
				},
				httpClassicView: []string{
					testProxy1,
				},
				httpsClassicView:  []string{},
				socks4ClassicView: []string{},
				socks5ClassicView: []string{},
				allAdvancedView: []entity.AdvancedProxy{
					testAdvancedProxyEntity1,
				},
				httpAdvancedView: []entity.Proxy{
					testProxyEntity1,
				},
				httpsAdvancedView:  []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{},
				socks5AdvancedView: []entity.Proxy{},
			},
			wantErr: nil,
		},
		{
			name: "StoreHTTPSProxy",
			fields: fields{
				allClassicView:     []string{},
				httpClassicView:    []string{},
				httpsClassicView:   []string{},
				socks4ClassicView:  []string{},
				socks5ClassicView:  []string{},
				allAdvancedView:    []entity.AdvancedProxy{},
				httpAdvancedView:   []entity.Proxy{},
				httpsAdvancedView:  []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{},
				socks5AdvancedView: []entity.Proxy{},
			},
			args: args{
				proxy: testProxyEntity2,
			},
			want: fields{
				allClassicView: []string{
					testProxy2,
				},
				httpClassicView: []string{},
				httpsClassicView: []string{
					testProxy2,
				},
				socks4ClassicView: []string{},
				socks5ClassicView: []string{},
				allAdvancedView: []entity.AdvancedProxy{
					testAdvancedProxyEntity2,
				},
				httpAdvancedView: []entity.Proxy{},
				httpsAdvancedView: []entity.Proxy{
					testProxyEntity2,
				},
				socks4AdvancedView: []entity.Proxy{},
				socks5AdvancedView: []entity.Proxy{},
			},
			wantErr: nil,
		},
		{
			name: "StoreSOCKS4Proxy",
			fields: fields{
				allClassicView:     []string{},
				httpClassicView:    []string{},
				httpsClassicView:   []string{},
				socks4ClassicView:  []string{},
				socks5ClassicView:  []string{},
				allAdvancedView:    []entity.AdvancedProxy{},
				httpAdvancedView:   []entity.Proxy{},
				httpsAdvancedView:  []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{},
				socks5AdvancedView: []entity.Proxy{},
			},
			args: args{
				proxy: testProxyEntity3,
			},
			want: fields{
				allClassicView: []string{
					testProxy3,
				},
				httpClassicView:  []string{},
				httpsClassicView: []string{},
				socks4ClassicView: []string{
					testProxy3,
				},
				socks5ClassicView: []string{},
				allAdvancedView: []entity.AdvancedProxy{
					testAdvancedProxyEntity3,
				},
				httpAdvancedView:  []entity.Proxy{},
				httpsAdvancedView: []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{
					testProxyEntity3,
				},
				socks5AdvancedView: []entity.Proxy{},
			},
			wantErr: nil,
		},
		{
			name: "StoreSOCKS5Proxy",
			fields: fields{
				allClassicView:     []string{},
				httpClassicView:    []string{},
				httpsClassicView:   []string{},
				socks4ClassicView:  []string{},
				socks5ClassicView:  []string{},
				allAdvancedView:    []entity.AdvancedProxy{},
				httpAdvancedView:   []entity.Proxy{},
				httpsAdvancedView:  []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{},
				socks5AdvancedView: []entity.Proxy{},
			},
			args: args{
				proxy: testProxyEntity4,
			},
			want: fields{
				allClassicView: []string{
					testProxy4,
				},
				httpClassicView:   []string{},
				httpsClassicView:  []string{},
				socks4ClassicView: []string{},
				socks5ClassicView: []string{
					testProxy4,
				},
				allAdvancedView: []entity.AdvancedProxy{
					testAdvancedProxyEntity4,
				},
				httpAdvancedView:   []entity.Proxy{},
				httpsAdvancedView:  []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{},
				socks5AdvancedView: []entity.Proxy{
					testProxyEntity4,
				},
			},
			wantErr: nil,
		},
		{
			name: "DuplicatedProxyWithinHTTPCategoryAndDifferentCategory",
			fields: fields{
				allClassicView: []string{
					testProxy3,
					testProxy4,
				},
				httpClassicView:  []string{},
				httpsClassicView: []string{},
				socks4ClassicView: []string{
					testProxy3,
				},
				socks5ClassicView: []string{
					testProxy4,
				},
				allAdvancedView: []entity.AdvancedProxy{
					testAdvancedProxyEntity3,
					testAdvancedProxyEntity4,
				},
				httpAdvancedView:  []entity.Proxy{},
				httpsAdvancedView: []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{
					testProxyEntity3,
				},
				socks5AdvancedView: []entity.Proxy{
					testProxyEntity4,
				},
			},
			args: args{
				proxy: entity.Proxy{
					Category:  testHTTPCategory,
					IP:        testIP4,
					Port:      testPort4,
					Proxy:     testProxy4,
					TimeTaken: testTimeTaken,
					CheckedAt: testCheckedAt,
				},
			},
			want: fields{
				allClassicView: []string{
					testProxy3,
					testProxy4,
				},
				httpClassicView: []string{
					testProxy4,
				},
				httpsClassicView: []string{},
				socks4ClassicView: []string{
					testProxy3,
				},
				socks5ClassicView: []string{
					testProxy4,
				},
				allAdvancedView: []entity.AdvancedProxy{
					testAdvancedProxyEntity3,
					{
						Proxy:     testAdvancedProxyEntity4.Proxy,
						IP:        testAdvancedProxyEntity4.IP,
						Port:      testAdvancedProxyEntity4.Port,
						TimeTaken: testTimeTaken,
						CheckedAt: testAdvancedProxyEntity4.CheckedAt,
						Categories: []string{
							testHTTPCategory,
							testProxyEntity4.Category,
						},
					},
				},
				httpAdvancedView: []entity.Proxy{
					{
						Category:  testHTTPCategory,
						IP:        testIP4,
						Port:      testPort4,
						Proxy:     testProxy4,
						TimeTaken: testTimeTaken,
						CheckedAt: testCheckedAt,
					},
				},
				httpsAdvancedView: []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{
					testProxyEntity3,
				},
				socks5AdvancedView: []entity.Proxy{
					testProxyEntity4,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ProxyRepository{
				AllClassicView:     tt.fields.allClassicView,
				HTTPClassicView:    tt.fields.httpClassicView,
				HTTPSClassicView:   tt.fields.httpsClassicView,
				SOCKS4ClassicView:  tt.fields.socks4ClassicView,
				SOCKS5ClassicView:  tt.fields.socks5ClassicView,
				AllAdvancedView:    tt.fields.allAdvancedView,
				HTTPAdvancedView:   tt.fields.httpAdvancedView,
				HTTPSAdvancedView:  tt.fields.httpsAdvancedView,
				SOCKS4AdvancedView: tt.fields.socks4AdvancedView,
				SOCKS5AdvancedView: tt.fields.socks5AdvancedView,
			}
			r.Store(&tt.args.proxy)

			views := map[string]struct {
				got  interface{}
				want interface{}
			}{
				"GetAllClassicView()":     {r.AllClassicView, tt.want.allClassicView},
				"GetHTTPClassicView()":    {r.HTTPClassicView, tt.want.httpClassicView},
				"GetHTTPSClassicView()":   {r.HTTPSClassicView, tt.want.httpsClassicView},
				"GetSOCKS4ClassicView()":  {r.SOCKS4ClassicView, tt.want.socks4ClassicView},
				"GetSOCKS5ClassicView()":  {r.SOCKS5ClassicView, tt.want.socks5ClassicView},
				"GetAllAdvancedView()":    {r.AllAdvancedView, tt.want.allAdvancedView},
				"GetHTTPAdvancedView()":   {r.HTTPAdvancedView, tt.want.httpAdvancedView},
				"GetHTTPSAdvancedView()":  {r.HTTPSAdvancedView, tt.want.httpsAdvancedView},
				"GetSOCKS4AdvancedView()": {r.SOCKS4AdvancedView, tt.want.socks4AdvancedView},
				"GetSOCKS5AdvancedView()": {r.SOCKS5AdvancedView, tt.want.socks5AdvancedView},
			}
			for name, v := range views {
				if !reflect.DeepEqual(v.got, v.want) {
					t.Errorf(expectedButGotMessage, name, v.want, v.got)
				}
			}
		})
	}
}

func TestGetAllClassicView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []string
	}{
		{
			name: "EmptyAllProxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					AllClassicView: []string{},
				}
			},
			want: []string{},
		},
		{
			name: "WithAllProxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.AllClassicView = []string{
					testProxy1,
					testProxy2,
					testProxy3,
					testProxy4,
				}
				return r
			},
			want: []string{
				testProxy1,
				testProxy2,
				testProxy3,
				testProxy4,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetAllClassicView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(expectedButGotMessage, "GetAllClassicView()", tt.want, got)
			}
		})
	}
}

func TestGetHTTPClassicView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []string
	}{
		{
			name: "EmptyHTTPProxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					HTTPClassicView: []string{},
				}
			},
			want: []string{},
		},
		{
			name: "WithHTTPProxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.HTTPClassicView = []string{
					testProxy1,
				}
				return r
			},
			want: []string{
				testProxy1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetHTTPClassicView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(expectedButGotMessage, "GetHTTPClassicView()", tt.want, got)
			}
		})
	}
}

func TestGetHTTPSClassicView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []string
	}{
		{
			name: "EmptyHTTPSProxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					HTTPSClassicView: []string{},
				}
			},
			want: []string{},
		},
		{
			name: "WithHTTPSProxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.HTTPSClassicView = []string{
					testProxy2,
				}
				return r
			},
			want: []string{
				testProxy2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetHTTPSClassicView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(expectedButGotMessage, "GetHTTPSClassicView()", tt.want, got)
			}
		})
	}
}

func TestGetSOCKS4ClassicView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []string
	}{
		{
			name: "EmptySOCKS4Proxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.SOCKS4ClassicView = []string{}
				return r
			},
			want: []string{},
		},
		{
			name: "WithSOCKS4Proxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					SOCKS4ClassicView: []string{
						testProxy3,
					},
				}
			},
			want: []string{
				testProxy3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetSOCKS4ClassicView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(expectedButGotMessage, "GetSOCKS4ClassicView()", tt.want, got)
			}
		})
	}
}

func TestGetSOCKS5ClassicView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []string
	}{
		{
			name: "EmptySOCKS5Proxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					SOCKS5ClassicView: []string{},
				}
			},
			want: []string{},
		},
		{
			name: "WithSOCKS5Proxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.SOCKS5ClassicView = []string{
					testProxy4,
				}
				return r
			},
			want: []string{
				testProxy4,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetSOCKS5ClassicView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(expectedButGotMessage, "GetSOCKS5ClassicView()", tt.want, got)
			}
		})
	}
}

func TestGetAllAdvancedView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []entity.AdvancedProxy
	}{
		{
			name: "EmptyAllProxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					AllAdvancedView: []entity.AdvancedProxy{},
				}
			},
			want: []entity.AdvancedProxy{},
		},
		{
			name: "WithAllProxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.AllAdvancedView = []entity.AdvancedProxy{
					testAdvancedProxyEntity1,
					testAdvancedProxyEntity2,
					testAdvancedProxyEntity3,
					testAdvancedProxyEntity4,
				}
				return r
			},
			want: []entity.AdvancedProxy{
				testAdvancedProxyEntity1,
				testAdvancedProxyEntity2,
				testAdvancedProxyEntity3,
				testAdvancedProxyEntity4,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetAllAdvancedView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(expectedButGotMessage, "GetAllAdvancedView()", tt.want, got)
			}
		})
	}
}

func TestGetHTTPAdvancedView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []entity.Proxy
	}{
		{
			name: "EmptyHTTPProxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					HTTPAdvancedView: []entity.Proxy{},
				}
			},
			want: []entity.Proxy{},
		},
		{
			name: "WithHTTPProxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.HTTPAdvancedView = []entity.Proxy{
					testProxyEntity1,
				}
				return r
			},
			want: []entity.Proxy{
				testProxyEntity1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetHTTPAdvancedView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(expectedButGotMessage, "GetHTTPAdvancedView()", tt.want, got)
			}
		})
	}
}

func TestGetHTTPSAdvancedView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []entity.Proxy
	}{
		{
			name: "EmptyHTTPSProxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					HTTPSAdvancedView: []entity.Proxy{},
				}
			},
			want: []entity.Proxy{},
		},
		{
			name: "WithHTTPSProxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.HTTPSAdvancedView = []entity.Proxy{
					testProxyEntity2,
				}
				return r
			},
			want: []entity.Proxy{
				testProxyEntity2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetHTTPSAdvancedView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(expectedButGotMessage, "GetHTTPSAdvancedView()", tt.want, got)
			}
		})
	}
}

func TestGetSOCKS4AdvancedView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []entity.Proxy
	}{
		{
			name: "EmptySOCKS4Proxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					SOCKS4AdvancedView: []entity.Proxy{},
				}
			},
			want: []entity.Proxy{},
		},
		{
			name: "WithSOCKSProxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.SOCKS4AdvancedView = []entity.Proxy{
					testProxyEntity3,
				}
				return r
			},
			want: []entity.Proxy{
				testProxyEntity3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetSOCKS4AdvancedView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(expectedButGotMessage, "GetSOCKS4AdvancedView()", tt.want, got)
			}
		})
	}
}

func TestGetSOCKS5AdvancedView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []entity.Proxy
	}{
		{
			name: "EmptySOCKS5Proxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					SOCKS5AdvancedView: []entity.Proxy{},
				}
			},
			want: []entity.Proxy{},
		},
		{
			name: "WithSOCKS5Proxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.SOCKS5AdvancedView = []entity.Proxy{
					testProxyEntity4,
				}
				return r
			},
			want: []entity.Proxy{
				testProxyEntity4,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetSOCKS5AdvancedView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(expectedButGotMessage, "GetSOCKS5AdvancedView()", tt.want, got)
			}
		})
	}
}
