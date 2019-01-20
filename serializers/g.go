package serializers

import (
	"time"
)

type ModelSerializer struct {
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ms *ModelSerializer) serializeBaseField(slug string, createAt, updateAt time.Time) {
	ms.Slug = slug
	ms.CreatedAt = createAt
	ms.UpdatedAt = updateAt
}
