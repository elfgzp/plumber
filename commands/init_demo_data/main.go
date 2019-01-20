package main

import (
	"github.com/elfgzp/plumber/database"
	"github.com/elfgzp/plumber/models"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

/*
	createdb plumber_db -U postgres -W
*/

func main() {
	db := database.ConnectToDB()
	defer db.Close()
	models.SetDB(db)
	db.LogMode(true)

	u1e := "plumber@elfgzp.cn"
	_ = models.CreateUser("admin", u1e, "123456")
	user, _ := models.GetUserByEmail(u1e)

	u2e := "demo@elfgzp.cm"
	_ = models.CreateUser("demo", u2e, "123456")

	_ = models.CreateTeam("plumber", user)
	_ = models.CreateTeam("plumber-front-end", user)
}
