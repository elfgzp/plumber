package serializers

import (
	"time"
)

type ModelSerializer struct {
	Slug          string    `json:"slug"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedBySlug string    `json:"created_by_slug"`
	UpdatedBySlug string    `json:"updated_by_slug"`
}

type BaseModelInterface interface {
	RelatedBaseField()
	BaseSlug() string
	BaseCreatedAt() time.Time
	BaseUpdatedAt() time.Time
	CreatedBySlug() string
	UpdatedBySlug() string
}

func (ms *ModelSerializer) serializeBaseField(baseModel BaseModelInterface) {
	ms.Slug = baseModel.BaseSlug()
	ms.CreatedAt = baseModel.BaseCreatedAt()
	ms.UpdatedAt = baseModel.BaseUpdatedAt()
	baseModel.RelatedBaseField()

	ms.CreatedBySlug = baseModel.CreatedBySlug()
	ms.UpdatedBySlug = baseModel.UpdatedBySlug()
}
