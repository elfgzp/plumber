package db

import (
	"fmt"
	"github.com/elfgzp/plumber/config"
	"github.com/jinzhu/gorm"
	"log"
)

/*
	Connect to database
 */
func ConnectToDB() *gorm.DB {
	log.Println("Connect to postgresql database ...")
	db, err := gorm.Open("postgres", config.PSQLConString)
	if err != nil {
		panic(fmt.Errorf("Failed to connect database %s", err))
	}
	db.SingularTable(true)

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "plumber_" + defaultTableName
	}

	return db
}
