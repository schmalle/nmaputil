package nmaputil

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var dbGorm *gorm.DB

func initDB(database string) bool {
	var err error
	dataSourceName := "nmaputil:GHGHG%%%DFDDDDDffff@tcp(127.0.0.1:3306)/nmaputil?parseTime=True"
	dbGorm, err = gorm.Open(database, dataSourceName)

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
