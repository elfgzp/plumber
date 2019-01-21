package serializers

import (
	"github.com/elfgzp/plumber/models"
	"time"
)

type TaskSerializer struct {
	ModelSerializer
	Name                 string    `json:"name"`
	Desc                 string    `json:"desc"`
	Sequence             int       `json:"sequence"`
	Deadline             time.Time `json:"deadline"`
	Doing                bool      `json:"doing"`
	Completed            bool      `json:"completed"`
	AssignedUserNickname string    `json:"assigned_user_nickname"`
	AssignedUserSlug     string    `json:"assigned_user_slug"`
	TaskListSlug         string    `json:"task_list_slug"`
	TaskListName         string    `json:"task_list_name"`
	ProjectSlug          string    `json:"project_slug"`
	ProjectName          string    `json:"project_name"`
}

func SerializeTask(t *models.Task) TaskSerializer {
	models.LoadRelatedObject(&t, &t.Assign, "Assign").
		Related(&t.TaskList, "TaskList").
		Related(&t.Project, "Project")
	ts := TaskSerializer{
		Name:                 t.Name,
		Desc:                 t.Desc,
		Sequence:             t.Sequence,
		Doing:                t.Doing,
		Deadline:             t.Deadline,
		Completed:            t.Completed,
		AssignedUserNickname: t.Assign.Nickname,
		AssignedUserSlug:     t.Assign.Slug,
		TaskListSlug:         t.TaskList.Slug,
		TaskListName:         t.TaskList.Name,
		ProjectSlug:          t.Project.Slug,
		ProjectName:          t.Project.Name,
	}
	ts.serializeBaseField(t.Slug, t.CreatedAt, t.UpdatedAt, t)
	return ts
}
