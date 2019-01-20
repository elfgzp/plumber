package serializers

import "github.com/elfgzp/plumber/models"

type TeamSerializer struct {
	ModelSerializer
	Name string `json:"name"`
}

func SerializeTeam(t *models.Team) TeamSerializer {
	us := TeamSerializer{
		Name: t.Name,
	}

	us.serializeBaseField(t.Slug, t.CreatedAt, t.UpdatedAt)
	return us
}
