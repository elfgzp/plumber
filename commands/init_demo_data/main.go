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
	tl1, _ := models.CreateTaskList("开发", p1.ID, u1)
	_, _ = models.CreateTaskList("测试", p1.ID, u1)
	_, _ = models.CreateTaskList("上线", p1.ID, u1)

	_, _ = models.CreateTask("登陆页面开发", nil, nil, tl1, u1)
	_, _ = models.CreateTask("注册页面开发", nil, nil, tl1, u1)
	task1, _ := models.CreateTask("任务列表页面开发", nil, nil, tl1, u1)

	_, _ = models.CreateTaskCheckpoint("检查项组件开发", nil, nil, task1, u1)
	_, _ = models.CreateTaskCheckpoint("评论组件开发", nil, nil, task1, u1)

}
