package restful

import (
	"encoding/json"
	"github.com/elfgzp/plumber/serializers"
	"net/http"
	"strconv"

	"github.com/elfgzp/plumber/helpers"
	"github.com/elfgzp/plumber/models"
)

func ListTeamHandler(w http.ResponseWriter, r *http.Request) {
	params := getRouteParams(r)
	query := getQuery(r)
	userSlug := params["userSlug"]
	//query := getQuery(r)
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 1
	} else if limit > MaxPageLimit {
		limit = MaxPageLimit
	}

	//if _, ok := query["page"]; ok {
	//	page = query["page"][0]
	//}

	user, _ := models.GetUserBySlug(userSlug)
	if user == nil {
		helpers.Response404(w, "User not found.")
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

	isMember := false
	for _, memberID := range team.MemberIDs() {
		if user.ID == memberID {
			isMember = true
			break
		}
	}

	if !isMember {
		helpers.Response403(w)
		return
	}

	helpers.Response200(w, "", serializers.SerializeTeam(team))
}

type TeamCreate struct {
	TeamName string `json:"teamName"`
}

func checkTeamCreate(teamCreate TeamCreate) []ErrorData {
	var errs []ErrorData
	if teamCreate.TeamName == "" {
		errs = append(errs, ErrorData{"teamName", "Team Name required"})
	}

	if checkTeamNameExist(teamCreate.TeamName) {
		errs = append(errs, ErrorData{"teamName", "Team Name exist"})
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

	if err := models.CreateTeam(teamCreate.TeamName, ru); err != nil {
		helpers.Response500(w)
		return
	} else {
		helpers.Response200(w, "", nil)
	}

}
