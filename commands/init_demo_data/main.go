package main

import (
	db2 "github.com/elfgzp/plumber/db"
	"github.com/elfgzp/plumber/models"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

/*
	createdb plumber_db -U postgres -W
*/

func main() {
	log.Println("Migrate database ...")
	db := db2.ConnectToDB()
	defer db.Close()
	models.SetDB(db)
	db.LogMode(true)

	demoUserEmail := "plumber@elfgzp.cn"
	_ = models.CreateUser("admin", demoUserEmail, "123456")
	user, _ := models.GetUserByEmail(demoUserEmail)

	_ = models.CreateTeam("plumber", user)
	_ = models.CreateTeam("plumber-front-end", user)
}
