package restful

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
)

func getQuery(r *http.Request) url.Values {
	u := r.URL
	query := u.Query()
	return query
}

func getRouteParams(r *http.Request) map[string]string {
	vars := mux.Vars(r)
	return vars
}
