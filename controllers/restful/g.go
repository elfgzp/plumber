package restful

import (
	"fmt"
	"github.com/elfgzp/plumber/models"
	"regexp"
)

type ErrorData struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// Check functions
func checkLen(fieldName, fieldValue string, minLen, maxLen int) string {
	lenField := len(fieldValue)

	if lenField < minLen {
		return fmt.Sprintf("%s field is too short, less than %d.", fieldName, minLen)
	}

	if lenField > maxLen {
		return fmt.Sprintf("%s field is too long, more than %d.", fieldName, maxLen)
	}

	return ""
}

func checkPassword(password string) string {
	return checkLen("Password", password, 6, 50)
}

func checkPwdRepeatMatch(pwd1, pwd2 string) string {
	if pwd1 != pwd2 {
		return fmt.Sprintf("2 password does not match.")
	}
	return ""
}

func checkEmail(email string) string {
	if len(email) == 0 {
		return fmt.Sprintf("Email field is required.")
	}

	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, email); !m {
		return fmt.Sprintf("Email field is not a valid email.")
	}
	return ""
}

func checkUserExistByEmail(email string) bool {
	user, _ := models.GetUserByEmail(email)

	return user != nil
}
