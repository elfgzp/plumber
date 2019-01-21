package models

type Project struct {
	Model
	Name string `gorm:"not null;"`
	Desc string

	Team   Team
	TeamID uint `gorm:"not null"`

	Owner   User
	OwnerID uint `gorm:"not null"`

	Members   []User `gorm:"many2many:project_user_rel;association_jointable_foreignkey:user_id"`
	TaskLists []TaskList
	Tasks     []Task
	Public    bool
	Active    bool
}

func GetProjectByName(name string) (*Project, error) {
	var p Project
	err := GetObjectByField(&p, "name", name)
	if err != nil {
		return nil, err
	}
	return &p, err
}

func GetProjectBySlug(slug string) (*Project, error) {
	var p Project
	err := GetObjectByField(&p, "slug", slug)
	if err != nil {
		return nil, err
	}
	return &p, err
}

func CreateProject(name, desc string, teamID uint, user *User, public bool) (*Project, error) {
	p := Project{Name: name, Desc: desc, TeamID: teamID, OwnerID: user.ID, Public: public, Active: true}
	p.SetCreatedBy(user)
	if err := db.Create(&p).Error; err != nil {
		return nil, err
	}
	_ = p.AddMember(user)
	return &p, nil
}

func (p *Project) MemberIDs() []uint {
	return p.Many2ManyIDs("project_user_rel", "project_id", "user_id")
}

func (p *Project) IsProjectMember(uid uint) bool {
	isMember := false
	for _, memberID := range p.MemberIDs() {
		if uid == memberID {
			isMember = true
			break
		}
	}
	return isMember
}

func (p *Project) AddMember(u *User) error {
	return db.Model(p).Association("Members").Append(u).Error
}

func (p *Project) AddTaskList(tl *TaskList) error {
	return db.Model(p).Association("TaskLists").Append(tl).Error
}

func (p *Project) AddTask(t *Task) error {
	return db.Model(p).Association("Task").Append(t).Error
}

func (p *Project) GetTaskListsLimit(page, limit int) (*[]TaskList, int, error) {
	taskLists := make([]TaskList, limit)
	total, err := GetObjectsByFieldLimit(&taskLists, &TaskList{}, page, limit, "sequence asc, updated_at desc", "project_id", p.ID)
	return &taskLists, total, err
}
