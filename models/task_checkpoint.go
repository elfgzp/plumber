package models

import "time"

type TaskCheckpoint struct {
	Model
	Desc      string `gorm:"not null;"`
	Sequence  int    `gorm:"AUTO_INCREMENT"`
	Completed bool

	Assign   User
	AssignID *uint

	Deadline *time.Time

	Task   Task
	TaskID uint `gorm:"not null;"`
}

func GetTaskCheckpointBySlug(slug string) (*TaskCheckpoint, error) {
	var taskCP TaskCheckpoint
	if err := GetObjectByField(&taskCP, "slug", slug); err != nil {
		return nil, err
	}
	return &taskCP, nil
}

func CreateTaskCheckpoint(desc string, deadline *time.Time, assignID *uint, task *Task, user *User) (*TaskCheckpoint, error) {
	taskCP := TaskCheckpoint{
		Desc:     desc,
		AssignID: assignID,
		Deadline: deadline,
		TaskID:   task.ID,
	}
	taskCP.SetCreatedBy(user)
	if err := db.Create(&taskCP).Error; err != nil {
		return nil, err
	}
	return &taskCP, nil
}
