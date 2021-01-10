package nmaputil

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)



type NmapRun struct {
	Nmaprun  xml.Name `xml:"nmaprun"`
	Scanner  string   `xml:"scanner,attr"`
	Start    string   `xml:"start,attr"`
	Startstr string   `xml:"startstr,attr"`
	Host     []Host     `xml:"host"`
}

type Port struct {
	Protocol string `xml:"protocol,attr"`
	Portid   string `xml:"portid,attr"`
	State    State  `xml:"state"`
}

type State struct {
	State  string `xml:"state,attr"`
	Reason string `xml:"reason,attr"`
}

type Host struct {
	Ports Ports `xml:"ports"`
}

type Ports struct {
	XMLName xml.Name `xml:"ports"`
	Ports   []Port   `xml:"port"`
}


func ParseWinnersXml(filename string) {
	// Open the xmlFile
	xmlFile, err := os.Open(filename)
	// if os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return
	}
	// defer the closing of xmlFile so that we can parse it.
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	// Unmarshal takes a []byte and fills the rss struct with the values found in the xmlFile
	nmaprun := NmapRun{}
	xml.Unmarshal(byteValue, &nmaprun)
	fmt.Println("Nmap scanner : " + nmaprun.Scanner)
	//	for i, item := range rss.Channel.Items {
	//		fmt.Println(item.Description)
	//	}
}

func main() {
	ParseWinnersXml("./testdata/targets.xml")

}

