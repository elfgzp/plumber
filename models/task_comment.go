package models

type TaskComment struct {
	Model
	Content string

	Task   Task
	TaskID uint
}
