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
		panic("failed to connect database " + database)
		return false
	}

	// Migration to create tables for NmapRun schema
	dbGorm.AutoMigrate(&NmapRun{})
	dbGorm.AutoMigrate(&Host{})
	dbGorm.AutoMigrate(&Port{})

	return true

} // initDB

func GetNumberOfHosts() int64 {

	// Get all records
	result := dbGorm.Find(&Host{})

	if result.Error != nil {
		return 0
	}

	return result.RowsAffected
}

//db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users)
