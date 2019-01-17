package main

import (
	db2 "github.com/elfgzp/plumber/db"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	db := db2.ConnectToDB()
	defer db.Close()
}
