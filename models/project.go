package models

import "log"

type Project struct {
	Model
	Name string `gorm:"not null;unique_index"`
	Desc string

	Team   Team
	TeamID uint `gorm:"not null"`

	Owner   User
	OwnerID uint `gorm:"not null"`

	Members    []User `gorm:"many2many:project_user_rel;association_jointable_foreignkey:user_id"`
	TaskStates []TaskState
	Tasks      []Task
	Public     bool
	Active     bool
}

func GetProjectByName(name string) (*Project, error) {
	var p Project
	if err := db.Where("name = ?", name).
		Find(&p).Error; err != nil {
		log.Fatalln("%s", err)
		return nil, err
	}
	return &p, nil
}

func CreateProject(name, desc string, teamID uint, user *User, public bool) error {
	p := Project{Name: name, Desc: desc, TeamID: teamID, OwnerID: user.ID, Public: public, Active: true}
	p.SetCreatedBy(user)
	if err := db.Create(&p).Error; err != nil {
		return err
	}
	_ = p.AddMember(user)
	return nil
}

func (p *Project) AddMember(u *User) error {
	return db.Model(p).Association("Members").Append(u).Error
}

func (p *Project) AddTaskState(ts *TaskState) error {
	return db.Model(p).Association("TaskStates").Append(ts).Error
}

func (p *Project) AddTask(t *Task) error {
	return db.Model(p).Association("Task").Append(t).Error
}
