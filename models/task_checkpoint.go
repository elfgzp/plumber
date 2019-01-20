package models

type TaskCheckPoint struct {
	Model
	Desc      string
	Sequence  int
	Completed bool

	User   User
	UserID uint

	Task   Task
	TaskID uint
}
