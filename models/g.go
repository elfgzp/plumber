package models

import (
	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
)

var db *gorm.DB

/*
	Set Database
 */
func SetDB(database *gorm.DB) {
	db = database
}

type BaseModel struct {
	gorm.Model
	Slug string `gorm:"unique_index"`
}

func (basModel *BaseModel) BeforeCreate() error {
	guid := xid.New()
	basModel.Slug = guid.String()
	return nil
}

type Model struct {
	gorm.Model
	CreatedBy   User
	CreatedByID int

	UpdatedBy   User
	UpdatedByID int

	DeletedBy   User
	DeletedByID int
}
