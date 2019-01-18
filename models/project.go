package models

type Project struct {
	Model
	Name string
	Desc string

	Team   Team
	TeamID int

	Owner   User
	OwnerID int

	Members    []User `gorm:"many2many:project_user_rel"`
	TaskStates []TaskState
	Tasks      []Task
	Public     bool
	Active     bool
}
