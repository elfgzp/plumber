package main

import (
	"fmt"
	"github.com/elfgzp/plumber/controllers/ws"
	"github.com/elfgzp/plumber/database"
	"github.com/elfgzp/plumber/models"
	"github.com/elfgzp/plumber/router"
	"github.com/gorilla/context"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

func main() {
	db := database.ConnectToDB()
	defer db.Close()
	models.SetDB(db)
	//db.LogMode(true)

	r := router.NewRouter()
	r.HandleFunc("/test_ws", ws.TestWSHandler)
	r.HandleFunc("/ws", ws.WebsocketHandler)

	log.Println("Serve start at http://127.0.0.1:8868")

	err := http.ListenAndServe(":8868", context.ClearHandler(r))
	if err != nil {
		panic(fmt.Errorf("Serve stop with error : %s", err))
	}
}
