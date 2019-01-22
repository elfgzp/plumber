package serializers

import (
	"github.com/elfgzp/plumber/models"
	"time"
)

type BaseTaskSerializer struct {
	ModelSerializer
	Name                 string     `json:"name"`
	Desc                 string     `json:"desc"`
	Sequence             int        `json:"sequence"`
	Deadline             *time.Time `json:"deadline"`
	Doing                bool       `json:"doing"`
	Completed            bool       `json:"completed"`
	AssignedUserNickname string     `json:"assigned_user_nickname"`
	AssignedUserSlug     string     `json:"assigned_user_slug"`
}

func (bts *BaseTaskSerializer) SerializeBaseTask(Name, Desc, AssignedUserNickname, AssignedUserSlug string, sequence int, deadline *time.Time, Doing, Completed bool) {
	bts.Name = Name
	bts.Desc = Desc
	bts.Sequence = sequence
	bts.Deadline = deadline
	bts.Doing = Doing
	bts.Completed = Completed
	bts.AssignedUserNickname = AssignedUserNickname
	bts.AssignedUserSlug = AssignedUserSlug
}

type ListTaskSerializer struct {
	ModelSerializer
	BaseTaskSerializer
}

type RetrieveTaskSerializer struct {
	ModelSerializer
	BaseTaskSerializer
	TaskListSlug string `json:"task_list_slug"`
	TaskListName string `json:"task_list_name"`
	ProjectSlug  string `json:"project_slug"`
	ProjectName  string `json:"project_name"`
}

func SerializeRetrieveListTask(t *models.Task) ListTaskSerializer {
	models.LoadRelatedObject(&t, &t.Assign, "Assign")
	models.LoadRelatedObject(&t, &t.TaskList, "TaskList")
	models.LoadRelatedObject(&t, &t.Project, "Project")

	ts := ListTaskSerializer{}
	ts.SerializeBaseTask(t.Name, t.Desc, t.Assign.Nickname, t.Assign.Slug, t.Sequence, t.Deadline, t.Doing, t.Completed)
	ts.serializeBaseField(t.Slug, t.CreatedAt, t.UpdatedAt, t)
	return ts
}

func SerializeRetrieveTask(t *models.Task) RetrieveTaskSerializer {
	models.LoadRelatedObject(&t, &t.Assign, "Assign")
	models.LoadRelatedObject(&t, &t.TaskList, "TaskList")
	models.LoadRelatedObject(&t, &t.Project, "Project")

	ts := RetrieveTaskSerializer{
		TaskListSlug: t.TaskList.Slug,
		TaskListName: t.TaskList.Name,
		ProjectSlug:  t.Project.Slug,
		ProjectName:  t.Project.Name,
	}
	ts.SerializeBaseTask(t.Name, t.Desc, t.Assign.Nickname, t.Assign.Slug, t.Sequence, t.Deadline, t.Doing, t.Completed)
	ts.serializeBaseField(t.Slug, t.CreatedAt, t.UpdatedAt, t)
	return ts
}
