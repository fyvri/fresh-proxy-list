package service

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/fyvri/fresh-proxy-list/internal/entity"
	"github.com/fyvri/fresh-proxy-list/pkg/utils"

	"h12.io/socks"
)

type ProxyService struct {
	FetcherUtil       utils.FetcherUtilInterface
	URLParserUtil     utils.URLParserUtilInterface
	HTTPTestingSites  []string
	HTTPSTestingSites []string
	UserAgents        []string
	Semaphore         chan struct{}
}

type ProxyServiceInterface interface {
	Check(category string, ip string, port string) (*entity.Proxy, error)
	GetTestingSite(category string) string
	GetRandomUserAgent() string
}

func NewProxyService(
	fetcherUtil utils.FetcherUtilInterface,
	urlParserUtil utils.URLParserUtilInterface,
	httpTestingSites []string,
	httpsTestingSites []string,
	userAgents []string,
) ProxyServiceInterface {
	return &ProxyService{
		FetcherUtil:       fetcherUtil,
		URLParserUtil:     urlParserUtil,
		HTTPTestingSites:  httpTestingSites,
		HTTPSTestingSites: httpsTestingSites,
		UserAgents:        userAgents,
		Semaphore:         make(chan struct{}, 500),
	}
}

func (s *ProxyService) Check(category string, ip string, port string) (*entity.Proxy, error) {
	s.Semaphore <- struct{}{}
	defer func() { <-s.Semaphore }()

	var (
		transport   *http.Transport
		proxy       = ip + ":" + port
		proxyURI    = strings.ToLower(category + "://" + proxy)
		testingSite = s.GetTestingSite(category)
		timeout     = 60 * time.Second
	)

	if category == "HTTP" || category == "HTTPS" {
		proxyURL, err := s.URLParserUtil.Parse(proxyURI)
		if err != nil {
			return nil, fmt.Errorf("error parsing proxy URL: %v", err)
		}

		transport = &http.Transport{
			Proxy:             http.ProxyURL(proxyURL),
			DisableKeepAlives: true,
			DialContext: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: timeout,
			}).DialContext,
			TLSHandshakeTimeout: timeout,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: category == "HTTPS",
			},
		}
	} else if category == "SOCKS4" || category == "SOCKS5" {
		proxyURL := socks.Dial(proxyURI)
		transport = &http.Transport{
			Dial:              proxyURL,
			DisableKeepAlives: true,
			DialContext: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: timeout,
			}).DialContext,
		}
	} else {
		return nil, fmt.Errorf("proxy category %s not supported", category)
	}

	req, err := s.FetcherUtil.NewRequest("GET", testingSite, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}
	req.Header.Set("User-Agent", s.GetRandomUserAgent())

	startTime := time.Now()
	resp, err := s.FetcherUtil.Do(&http.Client{
		Transport: transport,
		Timeout:   timeout,
	}, req)

	// statusCode := ""
	// if err == nil {
	// 	statusCode = http.StatusText(resp.StatusCode)
	// }
	// log.Printf("Check %s: %s ~> %s ~> %v", fmt.Sprintf("%-25s", proxy), fmt.Sprintf("%-30s", statusCode), testingSite, err)

	if err != nil {
		return nil, fmt.Errorf("request error: %s", err)
	}
	defer resp.Body.Close()
	endTime := time.Now()
	timeTaken := endTime.Sub(startTime).Seconds()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return &entity.Proxy{
		Proxy:     proxy,
		IP:        ip,
		Port:      port,
		Category:  category,
		CheckedAt: endTime.Format(time.RFC3339),
		TimeTaken: timeTaken,
	}, nil
}

func (s *ProxyService) GetTestingSite(category string) string {
	if category == "HTTPS" {
		return s.HTTPSTestingSites[rand.Intn(len(s.HTTPSTestingSites))]
	}
	return s.HTTPTestingSites[rand.Intn(len(s.HTTPTestingSites))]
}

func (s *ProxyService) GetRandomUserAgent() string {
	return s.UserAgents[rand.Intn(len(s.UserAgents))]
}
