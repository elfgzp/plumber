package main

import (
	"github.com/elfgzp/plumber/database"
	"github.com/elfgzp/plumber/models"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

/*
	createdb plumber_db -U postgres -W
*/

func main() {
	log.Println("Migrate database ...")
	db := database.ConnectToDB()
	defer db.Close()
	models.SetDB(db)
	db.LogMode(true)
	db.AutoMigrate(
		&models.User{},
		&models.Team{},
		&models.Project{},
		&models.Task{},
		&models.TaskList{},
		&models.TaskComment{},
		&models.TaskCheckpoint{},
	)
}
