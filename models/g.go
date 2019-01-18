package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

var db *gorm.DB

/*
	Set Database
 */
func SetDB(database *gorm.DB) {
	db = database
}

type BaseModel struct {
	ID         uint `gorm:"primary_key"`
	CreateAt   *time.Time
	UpdateAt   *time.Time
	DeletedAt  *time.Time
	Slug       string
}

type Model struct {
	BaseModel
	CreateBy   User
	UpdateBy   User
	CreateByID int
	UpdateByID int
}
