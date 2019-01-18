package models

type TaskCheckPoint struct {
	Model
	Desc     string
	Sequence int

	User   User
	UserID int

	Task   Task
	TaskID int
}
