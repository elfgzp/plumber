package restful

import (
	"encoding/json"
	"github.com/elfgzp/plumber/helpers"
	"github.com/elfgzp/plumber/models"
	"net/http"
)

func ListTeamHandler(w http.ResponseWriter, r *http.Request) {

}

func RetrieveTeamHandler(w http.ResponseWriter, r *http.Request) {

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
		errs = append(errs, ErrorData{"teamName", "Team Name required"})
	}

	return errs

}

func CreateTeamHandler(w http.ResponseWriter, r *http.Request) {
	var teamCreate TeamCreate
	params := getRouteParams(r)
	u, _ := models.GetUserBySlug(params["userSlug"])

	if u == nil {
		helpers.ResponseWithJSON(w, http.StatusNotFound, helpers.JSONResponse{Code: http.StatusNotFound, Msg:"User not found."})
	}

	ru := getRequestUser(r)
	if u.ID != ru.ID {
		helpers.ResponseWithJSON(w, http.StatusForbidden, helpers.JSONResponse{Code: http.StatusForbidden, Msg:"Permission denied."})
	}


	err := json.NewDecoder(r.Body).Decode(&teamCreate)
	if err != nil {
		helpers.ResponseWithJSON(w, http.StatusBadRequest, helpers.JSONResponse{Code: http.StatusBadRequest, Msg: "Json data invalid."})
		return
	}

	if errs := checkTeamCreate(teamCreate); len(errs) > 0 {
		helpers.ResponseWithJSON(w, http.StatusBadRequest, helpers.JSONResponse{Code: http.StatusBadRequest, Data: errs})
		return
	}

	if err := models.CreateTeam(teamCreate.TeamName, ru); err != nil {
		helpers.ResponseWithJSON(w, http.StatusInternalServerError, helpers.JSONResponse{Code: http.StatusInternalServerError})
		return
	} else {
		helpers.ResponseWithJSON(w, http.StatusOK, helpers.JSONResponse{Code: http.StatusOK})
	}

}
