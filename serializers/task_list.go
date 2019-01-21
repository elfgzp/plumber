package serializers

import "github.com/elfgzp/plumber/models"

type TaskListSerializer struct {
	ModelSerializer
	Name     string `json:"name"`
	Sequence int    `json:"sequence"`
}

func SerializeTaskList(tl *models.TaskList) TaskListSerializer {
	tls := TaskListSerializer{
		Name:     tl.Name,
		Sequence: tl.Sequence,
	}

	tls.serializeBaseField(tl.Slug, tl.CreatedAt, tl.UpdatedAt)
	return tls
}
