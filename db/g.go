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
	connectingStr := config.GetPostgreSQLConnectingString()
	log.Println("Connect to postgresql database ...")
	db, err := gorm.Open("postgres", connectingStr)
	if err != nil {
		panic(fmt.Errorf("Failed to connect database %s", err))
	}
	db.SingularTable(true)
	return db
}