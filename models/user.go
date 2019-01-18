package models

import "github.com/elfgzp/plumber/helpers"

type User struct {
	BaseModel
	Nickname          string
	Email             string
	MobileCountryCode string
	Mobile            string
	PasswordHash      string
	Team              []Team    `gorm:"many2many:team_user_rel;association_jointable_foreignkey;team_id"`
	Project           []Project `gorm:"many2many:project_user_rel;association_jointable_foreignkey:project_id"`
	Tasks             []Task
	StaredTasks       []Task `gorm:"many2many:stared_task_user_rel;association_jointable_foreignkey:task_id"`
	NotifiedTasks     []User `gorm:"many2many:notified_task_user_rel;association_jointable_foreignkey:task_id"`
}

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

func AddUser(nickname, email, password string) error {
	user := User{Nickname: nickname, Email: nickname}
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
