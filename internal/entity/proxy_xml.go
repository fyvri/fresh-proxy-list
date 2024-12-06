package entity

import "encoding/xml"

type ProxyXMLClassicView struct {
	XMLName xml.Name `xml:"Proxies"`
	Proxies []string `xml:"Proxy"`
}

type ProxyXMLAdvancedView struct {
	XMLName xml.Name `xml:"Proxies"`
	Proxies []Proxy  `xml:"Proxy"`
}

type ProxyXMLAllAdvancedView struct {
	XMLName xml.Name        `xml:"Proxies"`
	Proxies []AdvancedProxy `xml:"Proxy"`
}
