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
	ID       uint `gorm:"primary_key"`
	CreateBy User
	UpdateBy User
	CreateAt *time.Time
	UpdateAt *time.Time
	Slug     string
}
