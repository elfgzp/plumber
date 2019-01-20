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
		helpers.Response400(w, "Login and password required.", nil)
		return
	}

	user := Authenticate(userLogin.Login, userLogin.Password)
	if user == nil {
		helpers.Response400(w, "Login or password wrong.", nil)
	} else {
		token, _ := helpers.GenerateToken(user.Email)
		helpers.Response201(w, "", JWTToken{Token: token})
	}
}

func TokenVerifyHandler(w http.ResponseWriter, r *http.Request) {
	var jwtToken JWTToken
	err := json.NewDecoder(r.Body).Decode(&jwtToken)

	if err != nil || jwtToken.Token == "" {
		helpers.Response400(w, "Token required.", nil)
		return
	}

	email := helpers.CheckToken(jwtToken.Token)
	if email == "" {
		helpers.Response401(w)
	} else {
		helpers.Response200(w, "", JWTToken{Token: jwtToken.Token})
	}

}
