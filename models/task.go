package models

import "time"

type Task struct {
	Model
	Name      string
	Desc      string
	Sequence  int `gorm:"AUTO_INCREMENT"`
	Deadline  *time.Time
	Doing     bool
	Completed bool

	TaskList   TaskList
	TaskListID uint

	TaskCheckPoint []TaskCheckPoint
	TaskComments   []TaskComment

	Assign   User
	AssignID uint

	Project   Project
	ProjectID Project

	StaredUsers   []Task `gorm:"many2many:stared_task_user_rel;association_jointable_foreignkey:user_id"`
	NotifiedUsers []User `gorm:"many2many:notified_task_user_rel;association_jointable_foreignkey:user_id"`
}
