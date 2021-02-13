package nmaputil

import (
	"fmt"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

/**
collect targets from a given Excel file (1 colum prefix, 2. column domain name) and writes it to a second plaintext file
*/
func CollectTargets(source string, destination string) bool {

	counter := 0
	x := []string{}

	f, err := os.Create(destination)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fIn, err := excelize.OpenFile(source)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Get all the rows in the dns.
	rows, err := fIn.GetRows("dns")

	for _, row := range rows {
		entry := row[0] + "." + row[1]

		if false == strings.HasPrefix(entry, "@") {

			// if entry is new, add it
			if false == contains(x, entry) {
				x = append(x, entry)
				counter++
				f.WriteString(entry + "\n")
			}
		}
	}

	f.Close()
	fmt.Printf("Info: Inserted %d  entries...", counter)
	return true
}
