package serializers

import (
	"github.com/elfgzp/plumber/models"
	"time"
)

type TaskCheckpointSerializer struct {
	ModelSerializer
	Desc                 string     `json:"desc"`
	Sequence             int        `json:"sequence"`
	Deadline             *time.Time `json:"deadline"`
	Completed            bool       `json:"completed"`
	AssignedUserNickname string     `json:"assigned_user_nickname"`
	AssignedUserSlug     string     `json:"assigned_user_slug"`
}

func SerializeTaskCheckpoint(taskCP *models.TaskCheckpoint) TaskCheckpointSerializer {
	models.LoadRelatedObject(&taskCP, &taskCP.Assign, "Assign")
	taskCPS := TaskCheckpointSerializer{
		Desc:                 taskCP.Desc,
		Sequence:             taskCP.Sequence,
		Deadline:             taskCP.Deadline,
		Completed:            taskCP.Completed,
		AssignedUserNickname: taskCP.Assign.Nickname,
		AssignedUserSlug:     taskCP.Assign.Slug,
	}
	taskCPS.serializeBaseField(taskCP)
	return taskCPS
}
