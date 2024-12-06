package repository

import (
	"cmp"
	"slices"
	"sync"

	"github.com/fyvri/fresh-proxy-list/internal/entity"
)

type ProxyRepository struct {
	Mutex              sync.RWMutex
	AllClassicView     []string
	HTTPClassicView    []string
	HTTPSClassicView   []string
	SOCKS4ClassicView  []string
	SOCKS5ClassicView  []string
	AllAdvancedView    []entity.AdvancedProxy
	HTTPAdvancedView   []entity.Proxy
	HTTPSAdvancedView  []entity.Proxy
	SOCKS4AdvancedView []entity.Proxy
	SOCKS5AdvancedView []entity.Proxy
}

type ProxyRepositoryInterface interface {
	Store(proxy *entity.Proxy)
	GetAllClassicView() []string
	GetHTTPClassicView() []string
	GetHTTPSClassicView() []string
	GetSOCKS4ClassicView() []string
	GetSOCKS5ClassicView() []string
	GetAllAdvancedView() []entity.AdvancedProxy
	GetHTTPAdvancedView() []entity.Proxy
	GetHTTPSAdvancedView() []entity.Proxy
	GetSOCKS4AdvancedView() []entity.Proxy
	GetSOCKS5AdvancedView() []entity.Proxy
}

func NewProxyRepository() ProxyRepositoryInterface {
	return &ProxyRepository{
		Mutex: sync.RWMutex{},
	}
}

func (r *ProxyRepository) Store(proxy *entity.Proxy) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	updateProxyAll := func(proxy *entity.Proxy, classicList *[]string, advancedList *[]entity.AdvancedProxy) {
		n, found := slices.BinarySearchFunc(*advancedList, entity.AdvancedProxy{Proxy: proxy.Proxy}, func(a, b entity.AdvancedProxy) int {
			return cmp.Compare(a.Proxy, b.Proxy)
		})
		if found {
			if proxy.Category == "HTTP" && proxy.TimeTaken > 0 {
				(*advancedList)[n].TimeTaken = proxy.TimeTaken
			}

			if m, found := slices.BinarySearch((*advancedList)[n].Categories, proxy.Category); !found {
				(*advancedList)[n].Categories = slices.Insert((*advancedList)[n].Categories, m, proxy.Category)
			}
		} else {
			*classicList = append(*classicList, proxy.Proxy)
			*advancedList = slices.Insert(*advancedList, n, entity.AdvancedProxy{
				Proxy:     proxy.Proxy,
				IP:        proxy.IP,
				Port:      proxy.Port,
				TimeTaken: proxy.TimeTaken,
				CheckedAt: proxy.CheckedAt,
				Categories: []string{
					proxy.Category,
				},
			})
		}
	}

	switch proxy.Category {
	case "HTTP":
		r.HTTPClassicView = append(r.HTTPClassicView, proxy.Proxy)
		r.HTTPAdvancedView = append(r.HTTPAdvancedView, *proxy)
	case "HTTPS":
		r.HTTPSClassicView = append(r.HTTPSClassicView, proxy.Proxy)
		r.HTTPSAdvancedView = append(r.HTTPSAdvancedView, *proxy)
	case "SOCKS4":
		r.SOCKS4ClassicView = append(r.SOCKS4ClassicView, proxy.Proxy)
		r.SOCKS4AdvancedView = append(r.SOCKS4AdvancedView, *proxy)
	case "SOCKS5":
		r.SOCKS5ClassicView = append(r.SOCKS5ClassicView, proxy.Proxy)
		r.SOCKS5AdvancedView = append(r.SOCKS5AdvancedView, *proxy)
	}

	updateProxyAll(proxy, &r.AllClassicView, &r.AllAdvancedView)
}

func (r *ProxyRepository) GetAllClassicView() []string {
	return r.AllClassicView
}

func (r *ProxyRepository) GetHTTPClassicView() []string {
	return r.HTTPClassicView
}

func (r *ProxyRepository) GetHTTPSClassicView() []string {
	return r.HTTPSClassicView
}

func (r *ProxyRepository) GetSOCKS4ClassicView() []string {
	return r.SOCKS4ClassicView
}

func (r *ProxyRepository) GetSOCKS5ClassicView() []string {
	return r.SOCKS5ClassicView
}

func (r *ProxyRepository) GetAllAdvancedView() []entity.AdvancedProxy {
	return r.AllAdvancedView
}

func (r *ProxyRepository) GetHTTPAdvancedView() []entity.Proxy {
	return r.HTTPAdvancedView
}

func (r *ProxyRepository) GetHTTPSAdvancedView() []entity.Proxy {
	return r.HTTPSAdvancedView
}

func (r *ProxyRepository) GetSOCKS4AdvancedView() []entity.Proxy {
	return r.SOCKS4AdvancedView
}

func (r *ProxyRepository) GetSOCKS5AdvancedView() []entity.Proxy {
	return r.SOCKS5AdvancedView
}
