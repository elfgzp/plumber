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
		return nil
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
		return total, err
	}

	db.Model(model).Where(fmt.Sprintf("%s = ?", field), value).Count(&total)

	return total, nil
}

func GetObjectByField(out interface{}, field, value string) error {
	if err := db.Where(fmt.Sprintf("%s = ?", field), value).
		Find(out).Error; err != nil {
		return err
	}
	return nil
}

func UpdateObject(obj interface{}, contents map[string]interface{}, user *User) error {
	contents["updated_by_id"] = user.ID
	if err := db.Model(obj).Updates(contents).Error; err != nil {
		return err
	}
	return nil
}

func FakedDestroyObject(obj interface{}, user *User) error {
	contents := map[string]interface{}{"deleted_by_id": user.ID}
	_ = db.Model(obj).Updates(contents).Error
	if err := db.Delete(obj).Error; err != nil {
		return err
	}
	return nil
}

func LoadRelatedObject(obj interface{}, relatedObj interface{}, Field string) *gorm.DB {
	return db.Model(obj).Related(relatedObj, Field)
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

func (m *Model) RelatedBaseField() {
	db.Model(&m).
		Related(&m.CreatedBy, "CreatedBy").
		Related(&m.UpdatedBy, "UpdatedBy").
		Related(&m.DeletedBy, "DeletedBy")
}

func (m *Model) CreatedBySlug() string {
	return m.CreatedBy.Slug
}

func (m *Model) UpdatedBySlug() string {
	return m.UpdatedBy.Slug
}
