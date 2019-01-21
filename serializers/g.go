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
	CreatedBySlug() string
	UpdatedBySlug() string
}

func (ms *ModelSerializer) serializeBaseField(slug string, createAt, updateAt time.Time, baseModel BaseModelInterface) {
	ms.Slug = slug
	ms.CreatedAt = createAt
	ms.UpdatedAt = updateAt
	baseModel.RelatedBaseField()

	ms.CreatedBySlug = baseModel.CreatedBySlug()
	ms.UpdatedBySlug = baseModel.UpdatedBySlug()
}
