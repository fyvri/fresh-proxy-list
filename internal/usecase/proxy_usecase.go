package usecase

import (
	"fmt"
	"net"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/fyvri/fresh-proxy-list/internal/entity"
	"github.com/fyvri/fresh-proxy-list/internal/infrastructure/repository"
	"github.com/fyvri/fresh-proxy-list/internal/service"
)

type ProxyUsecase struct {
	ProxyRepository repository.ProxyRepositoryInterface
	ProxyService    service.ProxyServiceInterface
	ProxyMap        sync.Map
	SpecialIPs      []string
	PrivateIPs      []net.IPNet
}

type ProxyUsecaseInterface interface {
	ProcessProxy(category string, proxy string, isChecked bool) (*entity.Proxy, error)
	IsSpecialIP(ip string) bool
	GetAllAdvancedView() []entity.AdvancedProxy
}

func NewProxyUsecase(
	proxyRepository repository.ProxyRepositoryInterface,
	proxyService service.ProxyServiceInterface,
	specialIPs []string,
	privateIPs []net.IPNet,
) ProxyUsecaseInterface {
	return &ProxyUsecase{
		ProxyRepository: proxyRepository,
		ProxyService:    proxyService,
		SpecialIPs:      specialIPs,
		PrivateIPs:      privateIPs,
		ProxyMap:        sync.Map{},
	}
}

func (uc *ProxyUsecase) ProcessProxy(category string, proxy string, isChecked bool) (*entity.Proxy, error) {
	proxy = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(proxy, "\r", ""), "\n", ""))
	if proxy == "" {
		return nil, fmt.Errorf("proxy not found")
	}

	proxyParts := strings.Split(proxy, ":")
	if len(proxyParts) != 2 {
		return nil, fmt.Errorf("proxy format incorrect")
	}

	pattern := `^((25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\:(0|[1-9][0-9]{0,4})$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(proxy) {
		return nil, fmt.Errorf("proxy format not match")
	}

	if uc.IsSpecialIP(proxyParts[0]) {
		return nil, fmt.Errorf("proxy belongs to special ip")
	}

	port, err := strconv.Atoi(proxyParts[1])
	if err != nil || port < 0 || port > 65535 {
		return nil, fmt.Errorf("proxy port format incorrect")
	}

	_, loaded := uc.ProxyMap.LoadOrStore(category+"_"+proxy, true)
	if loaded {
		return nil, fmt.Errorf("proxy has been processed")
	}

	var (
		data               *entity.Proxy
		proxyIP, proxyPort = proxyParts[0], proxyParts[1]
	)
	if isChecked {
		data, err = uc.ProxyService.Check(category, proxyIP, proxyPort)
		if err != nil {
			return nil, err
		}
	} else {
		data = &entity.Proxy{
			Proxy:     proxy,
			IP:        proxyIP,
			Port:      proxyPort,
			Category:  category,
			TimeTaken: 0,
			CheckedAt: "",
		}
	}
	uc.ProxyRepository.Store(data)

	return data, nil
}

func (uc *ProxyUsecase) IsSpecialIP(ip string) bool {
	if _, found := slices.BinarySearch(uc.SpecialIPs, ip); found {
		return true
	}

	ipAddress := net.ParseIP(ip)
	if ipAddress == nil {
		return true
	}

	if ipAddress.IsLoopback() || ipAddress.IsMulticast() || ipAddress.IsUnspecified() {
		return true
	}

	for _, r := range uc.PrivateIPs {
		if r.Contains(ipAddress) {
			return true
		}
	}

	return false
}

func (uc *ProxyUsecase) GetAllAdvancedView() []entity.AdvancedProxy {
	return uc.ProxyRepository.GetAllAdvancedView()
}
