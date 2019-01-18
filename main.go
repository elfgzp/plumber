package main

import (
	"fmt"
	db2 "github.com/elfgzp/plumber/db"
	"github.com/elfgzp/plumber/models"
	"github.com/elfgzp/plumber/router"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

func main() {
	db := db2.ConnectToDB()
	defer db.Close()
	models.SetDB(db)

	r := router.NewRouter()

	log.Println("Serve start at http://127.0.0.1:8068")
	err := http.ListenAndServe(":8868", r)
	if err != nil {
		panic(fmt.Errorf("Serve stop with error : %s", err))
	}
}
