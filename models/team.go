package models

import (
	"fmt"
	"github.com/elfgzp/plumber/database"
	"log"
)

type Team struct {
	Model
	Name string `gorm:"not null;unique_index"`

	OwnerID uint
	Owner   User

	Members []User `gorm:"many2many:team_user_rel;association_jointable_foreignkey:user_id"`
}

func GetTeamByName(name string) (*Team, error) {
	return GetTeamByField("name", name)
}

func GetTeamBySlug(slug string) (*Team, error) {
	return GetTeamByField("slug", slug)
}

func GetTeamByField(field, value string) (*Team, error) {
	var team Team

	if err := db.Where(fmt.Sprintf("%s = ?", field), value).
		Find(&team).Error; err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return &team, nil
}

func (t *Team) MemberIDs() []uint {
	var ids []uint

	rows, err := db.Table(fmt.Sprintf("%s%s", database.TablePrefix, "team_user_rel")).Where("team_id = ?", t.ID).Select("team_id, user_id").Rows()

	if err != nil {
		log.Println("Counting team error: ", err)
	}

	defer rows.Close()
	for rows.Next() {
		var tid, UserID uint
		_ = rows.Scan(&tid, &UserID)
		ids = append(ids, UserID)
	}

	return ids
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
