package serializers

import "github.com/elfgzp/plumber/models"

type ProjectSerializer struct {
	ModelSerializer
	Name string
	Desc string
}

func SerializeProject(p *models.Project) ProjectSerializer {
	ps := ProjectSerializer{Name: p.Name, Desc: p.Desc}
	ps.serializeBaseField(p.Slug, p.CreatedAt, p.UpdatedAt, p)
	return ps
}
