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
	Slug string `gorm:"unique_index" json:"slug"`
}

func (baseModel *BaseModel) BeforeCreate() error {
	guid := xid.New()
	baseModel.Slug = guid.String()
	return nil
}

type Model struct {
	BaseModel
	CreatedBy   User
	CreatedByID uint

	UpdatedBy   User
	UpdatedByID uint

	DeletedBy   User
	DeletedByID uint
}

func (m *Model) SetCreatedBy(u *User) {
	m.CreatedByID = u.ID
	m.UpdatedByID = u.ID
}

func (m *Model) SetUpdatedBy(u *User) {
	m.UpdatedByID = u.ID
}

func (m *Model) SetDeletedBy(u *User) {
	m.DeletedByID = u.ID
}
