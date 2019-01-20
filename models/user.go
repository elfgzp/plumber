package models

import (
	"github.com/elfgzp/plumber/helpers"
	"log"
)

type User struct {
	BaseModel
	Nickname          string `gorm:"not null"`
	Email             string `gorm:"not null;unique_index"`
	MobileCountryCode string
	Mobile            string    `json:"mobile"`
	PasswordHash      string    `gorm:"not null"`
	Teams             []Team    `gorm:"many2many:team_user_rel;association_jointable_foreignkey:team_id"`
	Projects          []Project `gorm:"many2many:project_user_rel;association_jointable_foreignkey:project_id"`
	Tasks             []Task
	StaredTasks       []Task `gorm:"many2many:stared_task_user_rel;association_jointable_foreignkey:task_id"`
	NotifiedTasks     []User `gorm:"many2many:notified_task_user_rel;association_jointable_foreignkey:task_id"`
}

/*
Use to process login by email or mobile
*/
func GetUserByLogin(login string) (*User, error) {
	var user User

	err := db.Where("email = ?", login).
		//Or("mobile = ?", login).
		Find(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, err
}

func GetUserByEmail(email string) (*User, error) {
	var user User

	err := db.Where("email = ?", email).Find(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, err
}

func GetUserBySlug(slug string) (*User, error) {
	var user User

	err := db.Where("slug = ?", slug).Find(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, err
}

func CreateUser(nickname, email, password string) error {
	user := User{Nickname: nickname, Email: email}
	user.SetPassword(password)

	err := db.Create(&user).Error

	if err != nil {
		return err
	}

	return nil
}

func (u *User) SetPassword(password string) {
	u.PasswordHash = helpers.GeneratePasswordHash(password)
}

func (u *User) JoinedTeamIDs() []uint {
	return u.Many2ManyIDs("team_user_rel", "user_id", "team_id")
}

func (u *User) GetJoinedTeamLimit(page, limit int) (*[]Team, int, error) {
	var total int
	var teams []Team

	offset := (page - 1) * limit

	ids := u.JoinedTeamIDs()
	if err := db.Where("id in (?)", ids).
		Order("created_at asc").
		Offset(offset).
		Limit(limit).
		Find(&teams).Error; err != nil {
		log.Fatalln(err)
		return nil, total, err
	}

	db.Model(&Team{}).Where("id in (?)", ids).Count(&total)
	return &teams, total, nil
}
