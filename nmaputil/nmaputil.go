package nmaputil

import (
	"encoding/xml"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

/*
scan a given target with nmap (see www.insecure.org)
*/
func ScanTarget(target string, admin bool, persistance bool, database string) (NmapRun, bool) {

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

	return ParseXmlFile(file.Name(), persistance, database)
}

func ParseXmlFile(filename string, persistance bool, database string) (NmapRun, bool) {

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

		err := initDB(database)

		if !err {
			fmt.Println(currentTime.Format("2006-01-02 15:04:05") + ": " + " Database GORM setup failed")
			return nmaprun, false
		}

		// store nmaprun
		nmaprun.ID = nmaprun.Start
		dbGorm.Create(&nmaprun)

		nmaprunRun := nmaprun.Start // start start seconds as identifier

		for i := 0; i < len(nmaprun.Host); i++ {

			host := nmaprun.Host[i]

			host.NmapRunID = nmaprunRun
			host.Hostname = host.Hostnames.HostName[0].Name
			dbGorm.Create(&host)

			for j := 0; j < len(host.Ports.Ports); j++ {

				port := host.Ports.Ports[j]
				port.Hostname = host.Hostname
				port.NmapRunID = nmaprunRun
				port.StateDetail = port.State.State
				port.ServiceName = port.Service.Name
				port.ServiceProduct = port.Service.Product
				port.ServiceVersion = port.Service.Version

				httpserverheader := false

				for k := 0; k < len(port.Script); k++ {

					script := port.Script[k]

					if strings.Contains(script.Id, "http-server-header") {
						httpserverheader = true
					} else {
						httpserverheader = false
					}

					for l := 0; l < len(script.Elem); l++ {

						element := script.Elem[l]
						if strings.Contains(element.Key, "redirect") {
							port.RedirectURL = element.Data
						}

						if httpserverheader {
							port.ServiceProduct = port.ServiceProduct + "/" + element.Data
						}

					}

				}

				dbGorm.Create(&port)

			} // iterate through ports

		} // iterate through hosts

		fmt.Println("Done: ")

	}

	return nmaprun, true
}
