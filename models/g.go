package models

import (
	"fmt"
	"github.com/elfgzp/plumber/database"
	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
	"log"
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

func (baseModel *BaseModel) Many2ManyIDs(relTableName, IDField, targetIDField string) []uint {
	var ids []uint

	rows, err := db.Table(fmt.Sprintf("%s%s", database.TablePrefix, relTableName)).
		Where(fmt.Sprintf("%s = ?", IDField), baseModel.ID).
		Select(fmt.Sprintf("%s, %s", IDField, targetIDField)).Rows()

	if err != nil {
		log.Fatalln("%s %s %s get many2many ids error:", relTableName, IDField, targetIDField, err)
	}

	defer rows.Close()
	for rows.Next() {
		var id, targetID uint
		_ = rows.Scan(&id, &targetID)
		ids = append(ids, targetID)
	}

	return ids
}

func GetObjectByField(out interface{}, field, value string) error {
	if err := db.Where(fmt.Sprintf("%s = ?", field), value).
		Find(out).Error; err != nil {
		log.Fatalln(err)
		return err
	}
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
