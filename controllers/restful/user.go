package restful

import (
	"encoding/json"
	"github.com/elfgzp/plumber/helpers"
	"github.com/elfgzp/plumber/models"
	"github.com/elfgzp/plumber/serializers"
	"net/http"
)

type UserCreate struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email;"`
	Pwd1     string `json:"pwd1"`
	Pwd2     string `json:"pwd2"`
}

func checkUserCreate(userCreate UserCreate) []ErrorData {
	var errDataArr []ErrorData
	var errStr = ""
	if userCreate.Nickname == "" {
		errDataArr = append(errDataArr, ErrorData{Field: "nickname", Error: "Nick name required"})
	}

	errStr = checkEmail(userCreate.Email)
	if errStr != "" {
		errDataArr = append(errDataArr, ErrorData{Field: "email", Error: errStr})
	}

	if checkUserEmailExist(userCreate.Email) {
		errDataArr = append(errDataArr, ErrorData{Field: "email", Error: "This email has been registered"})
	}

	errStr = checkPassword(userCreate.Pwd1)
	if errStr != "" {
		errDataArr = append(errDataArr, ErrorData{Field: "pwd1", Error: errStr})
	}

	errStr = checkPwdRepeatMatch(userCreate.Pwd1, userCreate.Pwd2)
	if errStr != "" {
		errDataArr = append(errDataArr, ErrorData{Field: "pwd2", Error: errStr})
	}

	return errDataArr
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	u := getRequestUser(r)
	helpers.ResponseWithJSON(w, http.StatusOK, helpers.JSONResponse{Code: http.StatusOK, Data: serializers.SerializeUser(u)})
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var userCreate UserCreate
	err := json.NewDecoder(r.Body).Decode(&userCreate)

	if err != nil {
		helpers.ResponseWithJSON(w, http.StatusBadRequest, helpers.JSONResponse{Code: http.StatusBadRequest, Msg: "Json data invalid."})
		return
	}

	errs := checkUserCreate(userCreate)
	if len(errs) != 0 {
		helpers.ResponseWithJSON(w, http.StatusBadRequest, helpers.JSONResponse{Code: http.StatusBadRequest, Data: errs})
		return
	}

	if err := models.CreateUser(userCreate.Nickname, userCreate.Email, userCreate.Pwd1); err != nil {
		helpers.ResponseWithJSON(w, http.StatusInternalServerError, helpers.JSONResponse{Code: http.StatusInternalServerError,})
	} else {
		helpers.ResponseWithJSON(w, http.StatusOK, helpers.JSONResponse{Code: http.StatusOK})
	}

}
