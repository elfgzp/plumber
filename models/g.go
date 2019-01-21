package models

import (
	"fmt"
	"github.com/elfgzp/plumber/database"
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

func (baseModel *BaseModel) Many2ManyIDs(relTableName, IDField, targetIDField string) []uint {
	var ids []uint

	rows, err := db.Table(fmt.Sprintf("%s%s", database.TablePrefix, relTableName)).
		Where(fmt.Sprintf("%s = ?", IDField), baseModel.ID).
		Select(fmt.Sprintf("%s, %s", IDField, targetIDField)).Rows()

	if err != nil {
		fmt.Println("%s %s %s get many2many ids error:", relTableName, IDField, targetIDField, err)
	}

	defer rows.Close()
	for rows.Next() {
		var id, targetID uint
		_ = rows.Scan(&id, &targetID)
		ids = append(ids, targetID)
	}

	return ids
}

func GetObjectsByFieldLimit(objs interface{}, model interface{}, page, limit int, order, field string, value interface{}) (int, error) {
	var total int

	if order == "" {
		order = "id asc"
	}

	offset := (page - 1) * limit

	if err := db.Where(fmt.Sprintf("%s = ?", field), value).
		Order(order).
		Offset(offset).
		Limit(limit).
		Find(objs).Error; err != nil {
		fmt.Println(err)
		return total, err
	}

	db.Model(model).Where(fmt.Sprintf("%s = ?", field), value).Count(&total)

	return total, nil
}

func GetObjectByField(out interface{}, field, value string) error {
	if err := db.Where(fmt.Sprintf("%s = ?", field), value).
		Find(out).Error; err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func UpdateObject(obj interface{}, contents map[string]interface{}) error {
	return db.Model(obj).Updates(contents).Error
}

func FakedDestroyObject(obj interface{}) error {
	if err := db.Delete(obj).Error; err != nil {
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
