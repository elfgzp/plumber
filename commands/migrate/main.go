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

	db.AutoMigrate(
		&models.User{},
		&models.Team{},
		&models.Project{},
		&models.Task{},
		&models.TaskState{},
		&models.TaskComment{},
		&models.TaskCheckPoint{},
	)
}
