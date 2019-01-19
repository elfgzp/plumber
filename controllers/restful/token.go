package restful

import (
	"encoding/json"
	"github.com/elfgzp/plumber/helpers"
	"github.com/elfgzp/plumber/models"
	"net/http"
)

type UserLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JWTToken struct {
	Token string `json:"token"`
}

func Authenticate(login, pwd string) *models.User {
	user, _ := models.GetUserByLogin(login)
	if user == nil {
		return nil
	}

	if user.PasswordHash != helpers.GeneratePasswordHash(pwd) {
		return nil
	}

	return user
}

func CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var userLogin UserLogin
	err := json.NewDecoder(r.Body).Decode(&userLogin)
	if err != nil || userLogin.Login == "" || userLogin.Password == "" {
		helpers.ResponseWithJSON(w, http.StatusBadRequest, helpers.JSONResponse{Code: http.StatusBadRequest, Msg: "Login and password required."})
		return
	}

	user := Authenticate(userLogin.Login, userLogin.Password)
	if user == nil {
		helpers.ResponseWithJSON(w, http.StatusBadRequest, helpers.JSONResponse{Code: http.StatusBadRequest, Msg: "Login or password wrong."})
	} else {
		token, _ := helpers.GenerateToken(user.Email)
		helpers.ResponseWithJSON(w, http.StatusOK, helpers.JSONResponse{Code: http.StatusOK, Data: JWTToken{Token: token}})
	}
}

func TokenVerifyHandler(w http.ResponseWriter, r *http.Request) {
	var jwtToken JWTToken
	err := json.NewDecoder(r.Body).Decode(&jwtToken)

	if err != nil || jwtToken.Token == "" {
		helpers.ResponseWithJSON(w, http.StatusBadRequest, helpers.JSONResponse{Code: http.StatusBadRequest, Msg: "Token required."})
		return
	}

	email := helpers.CheckToken(jwtToken.Token)
	if email == "" {
		helpers.ResponseWithJSON(w, http.StatusUnauthorized, helpers.UnauthorizedResponse())
	} else {
		helpers.ResponseWithJSON(w, http.StatusOK, helpers.JSONResponse{Code: http.StatusOK, Data: JWTToken{Token: jwtToken.Token}})
	}

}
