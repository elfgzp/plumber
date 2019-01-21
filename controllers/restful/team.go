package restful

import (
	"encoding/json"
	"github.com/elfgzp/plumber/helpers"
	"github.com/elfgzp/plumber/models"
	"github.com/elfgzp/plumber/serializers"
	"net/http"
)

func ListTeamHandler(w http.ResponseWriter, r *http.Request) {
	ru := getRequestUser(r)
	params := getRouteParams(r)
	query := getQuery(r)
	userSlug := params["userSlug"]
	page, limit := getPageLimitQuery(query)

	user, _ := models.GetUserBySlug(userSlug)
	if user == nil {
		helpers.Response404(w, "User not found.")
		return
	}

	if userSlug != ru.Slug {
		helpers.Response403(w)
		return
	}

	teams, total, _ := user.GetJoinedTeamLimit(page, limit)
	var tsi []interface{}
	if total != 0 {
		tsi = make([]interface{}, len(*teams))
		for i, team := range *teams {
			tsi[i] = serializers.SerializeTeam(&team)
		}
	} else {
		tsi = make([]interface{}, 0)
	}
	helpers.Response200(w, "", helpers.PagedData{Total: total, Page: page, Limit: limit, Result: tsi})

}

func RetrieveTeamHandler(w http.ResponseWriter, r *http.Request) {
	user := getRequestUser(r)
	params := getRouteParams(r)
	teamSlug := params["teamSlug"]

	team, err := models.GetTeamBySlug(teamSlug)
	if err != nil {
		helpers.Response500(w)
		return
	}

	if team == nil {
		helpers.Response404(w, "Team not found.")
		return
	}

	if !team.IsTeamMember(user.ID) {
		helpers.Response403(w)
		return
	}

	helpers.Response200(w, "", serializers.SerializeTeam(team))
}

type TeamCreate struct {
	TeamName string `json:"teamName"`
}

func checkTeamCreate(teamCreate TeamCreate) []FieldCheckError {
	var errs []FieldCheckError
	if teamCreate.TeamName == "" {
		errs = append(errs, FieldCheckError{"teamName", "Team name required"})
	}

	if teamNameExist(teamCreate.TeamName) {
		errs = append(errs, FieldCheckError{"teamName", "Team name existed"})
	}

	return errs

}

func CreateTeamHandler(w http.ResponseWriter, r *http.Request) {
	var teamCreate TeamCreate
	params := getRouteParams(r)
	u, _ := models.GetUserBySlug(params["userSlug"])

	if u == nil {
		helpers.Response404(w, "User not found.")
		return
	}

	ru := getRequestUser(r)
	if u.ID != ru.ID {
		helpers.Response403(w)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&teamCreate)
	if err != nil {
		helpers.Response400(w, "Json data invalid.", nil)
		return
	}

	if errs := checkTeamCreate(teamCreate); len(errs) > 0 {
		helpers.Response400(w, "", errs)
		return
	}

	if team, err := models.CreateTeam(teamCreate.TeamName, ru); err != nil {
		helpers.Response500(w)
		return
	} else {
		helpers.Response201(w, "", serializers.SerializeTeam(team))
	}

}
