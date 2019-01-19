package models

type Project struct {
	Model
	Name string
	Desc string

	Team   Team
	TeamID uint

	Owner   User
	OwnerID uint

	Members    []User `gorm:"many2many:project_user_rel;association_jointable_foreignkey:user_id"`
	TaskStates []TaskState
	Tasks      []Task
	Public     bool
	Active     bool
}
