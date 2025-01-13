package opnsense

import "encoding/xml"

// Opnsense was generated 2025-01-13 00:41:41 by https://xml-to-go.github.io/ in Ukraine.
type Opnsense struct {
	XMLName xml.Name `xml:"opnsense"`
	Text    string   `xml:",chardata"`
	Dnsmasq struct {
		Text            string         `xml:",chardata"`
		Enable          string         `xml:"enable"`
		Regdhcpstatic   string         `xml:"regdhcpstatic"`
		StrictOrder     string         `xml:"strict_order"`
		Dnssec          string         `xml:"dnssec"`
		Interface       string         `xml:"interface"`
		Domainoverrides string         `xml:"domainoverrides"`
		Hosts           []OpnsenseHost `xml:"hosts"`
	} `xml:"dnsmasq"`
}

type OpnsenseHost struct {
	Text    string `xml:",chardata"`
	Host    string `xml:"host"`
	Domain  string `xml:"domain"`
	Ip      string `xml:"ip"`
	Descr   string `xml:"descr"`
	Aliases struct {
		Text string          `xml:",chardata"`
		Item []OpnsenseAlias `xml:"item"`
	} `xml:"aliases"`
}

type OpnsenseAlias struct {
	Text        string `xml:",chardata"`
	Description string `xml:"description"`
	Domain      string `xml:"domain"`
	Host        string `xml:"host"`
}
