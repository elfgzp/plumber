package models

type User struct {
	BaseModel
	Username      string
	Email         string
	Mobile        int32
	PasswordHash  string
	Team          []Team    `gorm:"many2many:team_user_rel;association_jointable_foreignkey;team_id"`
	Project       []Project `gorm:"many2many:project_user_rel;association_jointable_foreignkey:project_id"`
	Tasks         []Task
	StaredTasks   []Task `gorm:"many2many:stared_task_user_rel;association_jointable_foreignkey:task_id"`
	NotifiedTasks []User `gorm:"many2many:notified_task_user_rel;association_jointable_foreignkey:task_id"`
}
