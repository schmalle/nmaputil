package nmaputil

import "encoding/xml"

type NmapRun struct {
	ID       string `gorm:"primary_key"`
	Scanner  string `xml:"scanner,attr"`
	Start    string `xml:"start,attr"`
	Startstr string `xml:"startstr,attr"`
	Host     []Host `xml:"host"`
}

type Elem struct {
	Key  string `xml:"key,attr"`
	Data string `xml:",chardata"`
}

type Script struct {
	Id   string `xml:"id,attr"`
	Elem []Elem `xml:"elem"`
}

type Service struct {
	Name    string `xml:"name,attr"`
	Product string `xml:"product,attr"`
	Version string `xml:"version,attr"`
}

type Port struct {
	NmapRunID      string
	Hostname       string
	Protocol       string `xml:"protocol,attr"`
	Portid         string `xml:"portid,attr"`
	State          State  `xml:"state"`
	StateDetail    string
	Reason         string
	Script         []Script `xml:"script"`
	Service        Service  `xml:"service"` // e.g. IIS version x.y
	RedirectURL    string   // if it is a web service and a redirect was found
	ServiceName    string
	ServiceProduct string
	ServiceVersion string
}

type State struct {
	NmapRunID string
	State     string `xml:"state,attr"`
	Reason    string `xml:"reason,attr"`
}

type Host struct {
	NmapRunID string
	Ports     Ports     `xml:"ports"`
	Hostnames HostNames `xml:"hostnames"`
	Hostname  string
}

type HostName struct {
	NmapRunID string
	Name      string `xml:"name,attr"`
	Type      string `xml:"type,attr"`
}
type HostNames struct {
	NmapRunID string
	HostName  []HostName `xml:"hostname"`
}

type Ports struct {
	NmapRunID string
	Hostname  string
	XMLName   xml.Name `xml:"ports"`
	Ports     []Port   `xml:"port"`
}
