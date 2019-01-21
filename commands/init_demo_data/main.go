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
	u1, _ := models.CreateUser("admin", u1e, "123456")

	u2e := "demo@elfgzp.cm"
	_, _ = models.CreateUser("demo", u2e, "123456")

	t1, _ := models.CreateTeam("plumber", u1)
	_, _ = models.CreateTeam("plumber-front-end", u1)

	p1, _ := models.CreateProject("plumber-project", "Project management tool develop by golang", t1.ID, u1, false)
	_, _ = models.CreateProject("plumber-demo-project", "Project management tool develop by golang (Demo)", t1.ID, u1, false)

	_, _ = models.CreateTaskList("需求", p1.ID, u1)
	_, _ = models.CreateTaskList("设计", p1.ID, u1)
	_, _ = models.CreateTaskList("开发", p1.ID, u1)
	_, _ = models.CreateTaskList("测试", p1.ID, u1)
	_, _ = models.CreateTaskList("上线", p1.ID, u1)
}
