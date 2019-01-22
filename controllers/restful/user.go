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

func ValidUserCreate(userCreate UserCreate) []FieldValidError {
	var errDataArr []FieldValidError
	var errStr = ""
	if userCreate.Nickname == "" {
		errDataArr = append(errDataArr, FieldValidError{Field: "nickname", Error: "Nick name required"})
	}

	errStr = ValidEmail(userCreate.Email)
	if errStr != "" {
		errDataArr = append(errDataArr, FieldValidError{Field: "email", Error: errStr})
	}

	if userEmailExist(userCreate.Email) {
		errDataArr = append(errDataArr, FieldValidError{Field: "email", Error: "This email has been registered"})
	}

	errStr = ValidPassword(userCreate.Pwd1)
	if errStr != "" {
		errDataArr = append(errDataArr, FieldValidError{Field: "pwd1", Error: errStr})
	}

	errStr = ValidPwdRepeatMatch(userCreate.Pwd1, userCreate.Pwd2)
	if errStr != "" {
		errDataArr = append(errDataArr, FieldValidError{Field: "pwd2", Error: errStr})
	}

	return errDataArr
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	u := getRequestUser(r)
	helpers.Response200(w, "", serializers.SerializeUser(u))
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var userCreate UserCreate
	err := json.NewDecoder(r.Body).Decode(&userCreate)

	if err != nil {
		helpers.Response400(w, "Json data invalid.", nil)
		return
	}

	errs := ValidUserCreate(userCreate)
	if len(errs) != 0 {
		helpers.Response400(w, "", errs)
		return
	}

	if user, err := models.CreateUser(userCreate.Nickname, userCreate.Email, userCreate.Pwd1); err != nil {
		helpers.Response500(w)
	} else {
		helpers.Response201(w, "", serializers.SerializeUser(user))
	}

}
