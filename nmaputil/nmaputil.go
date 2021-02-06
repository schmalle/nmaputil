package nmaputil

import (
	"encoding/xml"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

var dbGorm *gorm.DB

type NmapRun struct {
	ID       uint   `gorm:"primary_key"`
	Nmaprun  string `xml:"nmaprun"`
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

type Port struct {
	NmapRunID   string
	Hostname    string
	Protocol    string `xml:"protocol,attr"`
	Portid      string `xml:"portid,attr"`
	State       State  `xml:"state"`
	StateDetail string
	Reason      string
	Script      []Script `xml:"script"`
	Service     string   `xml:"service"` // e.g. IIS version x.y
	RedirectURL string   // if it is a web service and a redirect was found
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

func initDB() bool {
	var err error
	dataSourceName := "nmaputil:GHGHG%%%DFDDDDDffff@tcp(127.0.0.1:3306)/nmaputil?parseTime=True"
	dbGorm, err = gorm.Open("mysql", dataSourceName)

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
		return false
	}

	// Migration to create tables for NmapRun schema
	dbGorm.AutoMigrate(&NmapRun{})
	dbGorm.AutoMigrate(&Host{})
	dbGorm.AutoMigrate(&Port{})

	return true

} // initDB

/*
scan a given target with nmap (see www.insecure.org)
*/
func ScanTarget(target string, admin bool, persistance bool) (NmapRun, bool) {

	nmaprun := NmapRun{}

	prg := "nmap"
	currentTime := time.Now()

	arg0 := "-A"
	arg1 := ""

	if admin {
		arg1 = "-sS"
	}

	arg2 := target
	arg3 := "-oX"

	file, err := ioutil.TempFile(os.TempDir(), "scan.*.xml")
	if err != nil {
		log.Fatal(err)
		return nmaprun, false
	}

	fmt.Println(currentTime.Format("2006-01-02 15:04:05") + ": " + "Starting scan of target " + target + " with target scan file " + file.Name())

	cmd := exec.Command(prg, arg0, arg1, arg2, arg3, file.Name())
	stdout, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
		return nmaprun, false
	}

	print(stdout)

	return ParseXmlFile(file.Name(), persistance)
}

func ParseXmlFile(filename string, persistance bool) (NmapRun, bool) {

	nmaprun := NmapRun{}
	currentTime := time.Now()

	// Open the xmlFile
	xmlFile, err := os.Open(filename)
	// if os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return nmaprun, false
	}
	// defer the closing of xmlFile so that we can parse it.
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	// Unmarshal takes a []byte and fills the rss struct with the values found in the xmlFile
	xml.Unmarshal(byteValue, &nmaprun)
	fmt.Println("Nmap scanner : " + nmaprun.Scanner)

	if persistance {

		err := initDB()

		if !err {
			fmt.Println(currentTime.Format("2006-01-02 15:04:05") + ": " + " Database GORM setup failed")
			return nmaprun, false
		}

		// store nmaprun
		//result := dbGorm.Create(&nmaprun)

		nmaprunRun := nmaprun.Start // start start seconds as identifier

		for i := 0; i < len(nmaprun.Host); i++ {

			host := nmaprun.Host[i]

			host.NmapRunID = nmaprunRun
			host.Hostname = host.Hostnames.HostName[0].Name
			//dbGorm.Create(&host)

			for j := 0; j < len(host.Ports.Ports); j++ {

				port := host.Ports.Ports[j]
				port.Hostname = host.Hostname
				port.NmapRunID = nmaprunRun
				port.StateDetail = port.State.State
				//dbGorm.Create(&port)

			}

		}

		fmt.Println("Done: ")

	}

	return nmaprun, true
}
