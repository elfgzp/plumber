package restful

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

var (
	DefaultPageLimit int
)

func init() {
	DefaultPageLimit = 10
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
