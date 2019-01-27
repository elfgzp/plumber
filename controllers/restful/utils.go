package restful

import (
	"fmt"
	"github.com/elfgzp/plumber/models"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type FieldValidError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func getRequestUser(r *http.Request) *models.User {
	userEmail, ok := context.GetOk(r, "userEmail")
	if !ok {
		return nil
	}

	user, _ := models.GetUserByEmail(userEmail.(string))

	return user
}

func getQuery(r *http.Request) url.Values {
	u := r.URL
	query := u.Query()
	return query
}

func getRouteParams(r *http.Request) map[string]string {
	vars := mux.Vars(r)
	return vars
}

func getPageLimitQuery(query url.Values) (int, int) {
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit := DefaultPageLimit

	if query.Get("limit") != "" {
		limit, err = strconv.Atoi(query.Get("limit"))
		if err != nil || limit <= 0 {
			limit = 1
		} else if limit > MaxPageLimit {
			limit = MaxPageLimit
		}
	}
	return page, limit
}

// Check functions
func ValidLen(fieldName, fieldValue string, minLen, maxLen int) string {
	lenField := len(fieldValue)

	if lenField < minLen {
		return fmt.Sprintf("%s field is too short, less than %d.", fieldName, minLen)
	}

	if lenField > maxLen {
		return fmt.Sprintf("%s field is too long, more than %d.", fieldName, maxLen)
	}

	return ""
}

func ValidPassword(password string) string {
	return ValidLen("Password", password, 6, 50)
}

func ValidPwdRepeatMatch(pwd1, pwd2 string) string {
	if pwd1 != pwd2 {
		return fmt.Sprintf("2 password does not match.")
	}
	return ""
}

func ValidEmail(email string) string {
	if len(email) == 0 {
		return fmt.Sprintf("Email field is required.")
	}

	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, email); !m {
		return fmt.Sprintf("Email field is not a valid email.")
	}
	return ""
}

func userEmailExist(email string) bool {
	user, _ := models.GetUserByEmail(email)

	return user != nil
}

func teamNameExist(teamName string) bool {
	team, _ := models.GetTeamByName(teamName)

	return team != nil
}

func projectNameExist(name string) bool {
	project, _ := models.GetProjectByName(name)

	return project != nil
}

func taskListNameExsit(name string) bool {
	taskList, _ := models.GetTaskListByName(name)

	return taskList != nil
}

func validIntField(contents map[string]interface{}, data map[string]interface{}, field string, errs []FieldValidError) (map[string]interface{}, []FieldValidError) {
	if val, ok := data[field]; ok && val != nil {
		if _, ok = val.(float64); !ok {
			errs = append(errs, FieldValidError{Field: field, Error: "Must be a number."})
		} else {
			contents[field] = int(val.(float64))
		}
	}
	return contents, errs
}

func validStringField(contents map[string]interface{}, data map[string]interface{}, field string, errs []FieldValidError) (map[string]interface{}, []FieldValidError) {
	if val, ok := data[field]; ok {
		if val != nil {
			if _, ok = val.(string); !ok {
				errs = append(errs, FieldValidError{Field: field, Error: "Must be a string."})
			} else {
				contents[field] = val.(string)
			}
		} else {
			contents[field] = nil
		}

	}
	return contents, errs
}

func validDatetimeField(contents map[string]interface{}, data map[string]interface{}, field string, errs []FieldValidError) (map[string]interface{}, []FieldValidError) {
	if val, ok := data[field]; ok {
		if val == nil {
			contents[field] = nil
		} else if _, ok := val.(string); !ok {
			errs = append(errs, FieldValidError{Field: field, Error: "Must be a datetime (RFC3339Nano)."})
		} else if t, err := time.Parse(time.RFC3339Nano, val.(string)); err != nil {
			errs = append(errs, FieldValidError{Field: field, Error: "Must be a datetime (RFC3339Nano)."})
		} else {
			contents[field] = t
		}
	}
	return contents, errs
}

func validBoolField(contents map[string]interface{}, data map[string]interface{}, field string, errs []FieldValidError) (map[string]interface{}, []FieldValidError) {
	if val, ok := data[field]; ok && val != nil {
		if _, ok = val.(bool); !ok {
			errs = append(errs, FieldValidError{Field: field, Error: "Must be a bool."})
		} else {
			contents[field] = val.(bool)
		}
	}
	return contents, errs
}
