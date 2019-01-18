package models

type Team struct {
	Model
	Name string

	OwnerID int
	Owner   User

	Members []User `gorm:"many2many:team_user_rel;"`
}
