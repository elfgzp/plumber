package models

type TaskCheckPoint struct {
	Model
	Desc     string
	Sequence int

	User   User
	UserID uint

	Task   Task
	TaskID uint
}
