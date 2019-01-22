package models

import "time"

type Task struct {
	Model
	Name      string `gorm:"not null;"`
	Desc      string
	Sequence  int `gorm:"AUTO_INCREMENT"`
	Deadline  *time.Time
	Doing     bool
	Completed bool

	TaskList   TaskList
	TaskListID uint

	TaskCheckpoint []TaskCheckpoint
	TaskComments   []TaskComment

	Assign   User
	AssignID *uint

	Project   Project
	ProjectID uint `gorm:"not null;"`

	StaredUsers   []Task `gorm:"many2many:stared_task_user_rel;association_jointable_foreignkey:user_id"`
	NotifiedUsers []User `gorm:"many2many:notified_task_user_rel;association_jointable_foreignkey:user_id"`
}

func GetTaskBySlug(slug string) (*Task, error) {
	var t Task
	if err := GetObjectByField(&t, "slug", slug); err != nil {
		return nil, err
	}
	return &t, nil
}

func CreateTask(name string, deadline *time.Time, assignID *uint, taskList *TaskList, user *User) (*Task, error) {
	t := Task{
		Name:       name,
		AssignID:   assignID,
		Deadline:   deadline,
		TaskListID: taskList.ID,
		ProjectID:  taskList.ProjectID,
	}
	t.SetCreatedBy(user)
	if err := db.Create(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (tl *Task) GetTaskCheckpointsLimit(page, limit int) (*[]TaskCheckpoint, int, error) {
	var taskCPs []TaskCheckpoint
	total, err := GetObjectsByFieldLimit(&taskCPs, &TaskCheckpoint{}, page, limit, "sequence asc, created_at desc", "task_id", tl.ID)
	return &taskCPs, total, err
}
