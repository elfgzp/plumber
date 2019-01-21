package models

type TaskCheckPoint struct {
	Model
	Desc      string
	Sequence  int `gorm:"AUTO_INCREMENT"`
	Completed bool

	User   User
	UserID uint

	Task   Task
	TaskID uint
}
