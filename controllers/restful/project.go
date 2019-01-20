package restful

import (
	"encoding/json"
	"github.com/elfgzp/plumber/helpers"
	"github.com/elfgzp/plumber/models"
	"github.com/elfgzp/plumber/serializers"
	"net/http"
)

type ProjectCreate struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func checkProjectCreate(projectCreate ProjectCreate) []FieldCheckError {
	var errs []FieldCheckError

	if projectCreate.Name == "" {
		errs = append(errs, FieldCheckError{Field: "name", Error: "Project name required."})
	}

	if projectNameExist(projectCreate.Name) {
		errs = append(errs, FieldCheckError{Field: "name", Error: "Project name existed."})

	}

	return errs
}

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	var projectCreate ProjectCreate
	params := getRouteParams(r)
	ru := getRequestUser(r)

	team, _ := models.GetTeamBySlug(params["teamSlug"])
	if team == nil {
		helpers.Response404(w, "Team not found.")
		return
	}

	if !team.IsTeamMember(ru.ID) {
		helpers.Response403(w)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&projectCreate); err != nil {
		helpers.Response400(w, "JSON invalid.", nil)
		return
	}

	if errs := checkProjectCreate(projectCreate); len(errs) > 0 {
		helpers.Response400(w, "", errs)
		return
	}

	if err := models.CreateProject(projectCreate.Name, projectCreate.Desc, team.ID, ru, false); err != nil {
		helpers.Response500(w)
		return
	} else {
		helpers.Response201(w, "", nil)
	}

}

func ListProjectHandler(w http.ResponseWriter, r *http.Request) {
	ru := getRequestUser(r)
	params := getRouteParams(r)
	query := getQuery(r)
	teamSlug := params["teamSlug"]
	page, limit := getPageLimitQuery(query)

	team, _ := models.GetTeamBySlug(teamSlug)
	if team == nil {
		helpers.Response404(w, "Team not found.")
		return
	}

	if !team.IsTeamMember(ru.ID) {
		helpers.Response403(w)
		return
	}

	projects, total, _ := team.GetTeamProjectsLimit(page, limit)

	var psi []interface{}
	if total != 0 {
		psi = make([]interface{}, len(*projects))
		for i, project := range *projects {
			psi[i] = serializers.SerializeProject(&project)
		}
	} else {
		psi = make([]interface{}, 0)
	}

	helpers.Response200(w, "", helpers.PagedData{Total: total, Page: page, Limit: limit, Result: psi})

}

func RetrieveProjectHandler(w http.ResponseWriter, r *http.Request) {
	ru := getRequestUser(r)
	params := getRouteParams(r)
	projectSlug := params["projectSlug"]

	project, err := models.GetProjectBySlug(projectSlug)
	if err != nil {
		helpers.Response500(w)
		return
	}

	if !project.IsProjectMember(ru.ID) {
		helpers.Response403(w)
		return
	}

	helpers.Response200(w, "", serializers.SerializeProject(project))
}
