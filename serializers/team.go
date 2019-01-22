package serializers

import "github.com/elfgzp/plumber/models"

type TeamSerializer struct {
	ModelSerializer
	Name string `json:"name"`
}

func SerializeTeam(t *models.Team) TeamSerializer {
	ts := TeamSerializer{
		Name: t.Name,
	}

	ts.serializeBaseField(t)
	return ts
}
