package models

import "log"

type Team struct {
	Model
	Name string `gorm:"not null;unique_index"`

	OwnerID uint
	Owner   User

	Members []User `gorm:"many2many:team_user_rel;association_jointable_foreignkey:user_id"`
}

func GetTeamByName(name string) (*Team, error) {
	var team Team

	if err := db.Where("name = ?", name).
		Find(&team).Error; err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return &team, nil
}

func CreateTeam(name string, user *User) error {
	team := Team{Name: name, OwnerID: user.ID}
	team.SetCreatedBy(user)
	if err := db.Create(&team).Error; err != nil {
		return err
	}
	_ = team.AddMember(user)
	return nil
}

func (t *Team) AddMember(u *User) error {
	return db.Model(t).Association("Members").Append(u).Error
}
