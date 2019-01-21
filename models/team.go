package models

type Team struct {
	Model
	Name string `gorm:"not null;"`

	OwnerID uint `gorm:"not null"`
	Owner   User

	Members []User `gorm:"many2many:team_user_rel;association_jointable_foreignkey:user_id"`
}

func GetTeamByName(name string) (*Team, error) {
	var t Team
	err := GetObjectByField(&t, "name", name)
	if err != nil {
		return nil, err
	}
	return &t, err
}

func GetTeamBySlug(slug string) (*Team, error) {
	var t Team
	err := GetObjectByField(&t, "slug", slug)
	if err != nil {
		return nil, err
	}
	return &t, err
}

func (t *Team) MemberIDs() []uint {
	return t.Many2ManyIDs("team_user_rel", "team_id", "user_id")
}

func CreateTeam(name string, user *User) (*Team, error) {
	t := Team{Name: name, OwnerID: user.ID}
	t.SetCreatedBy(user)
	if err := db.Create(&t).Error; err != nil {
		return nil, err
	}
	_ = t.AddMember(user)
	return &t, nil
}

func (t *Team) AddMember(u *User) error {
	return db.Model(t).Association("Members").Append(u).Error
}

func (t *Team) IsTeamMember(uid uint) bool {
	isMember := false
	for _, memberID := range t.MemberIDs() {
		if uid == memberID {
			isMember = true
			break
		}
	}
	return isMember
}

func (t *Team) GetTeamProjectsLimit(page, limit int) (*[]Project, int, error) {
	projects := make([]Project, limit)
	total, err := GetObjectsByFieldLimit(&projects, &Project{}, page, limit, "created_at desc", "team_id", t.ID)
	return &projects, total, err
}
