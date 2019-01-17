package models

type User struct {
	BaseModel
	Username     string
	Email        string
	Mobile       int32
	PasswordHash string
}
