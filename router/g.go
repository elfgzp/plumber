package router

import (
	"fmt"
	"github.com/elfgzp/plumber/controllers/restful"
	"github.com/elfgzp/plumber/helpers"
	"github.com/elfgzp/plumber/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name       string
	Method     string
	URI        string
	Handler    http.HandlerFunc
	MiddleWare mux.MiddlewareFunc
}

var routers []Route

func init() {
	register("apiURL", http.MethodGet, "/api", func(w http.ResponseWriter, r *http.Request) {
		apiURL := make(map[string]string)
		for _, route := range routers {
			if _, ok := apiURL[route.Name]; !ok {
				apiURL[route.Name] = fmt.Sprintf("http://%s%s", r.Host, route.URI)
			}
		}
		helpers.ResponseWithJSON(w, http.StatusOK, apiURL)
	}, nil)

	register("tokenURL", http.MethodPost, "/api/token", restful.CreateTokenHandler, nil)
	register("tokenVerificationURL", http.MethodPost, "/api/token/verification", restful.TokenVerifyHandler, nil)

	register("currentUserUrl", http.MethodGet, "/api/users/current", restful.GetCurrentUser, middleware.JWTTokenAuthMiddleware)
	register("usersURL", http.MethodPost, "/api/users", restful.CreateUserHandler, nil)

	register("teamsURL", http.MethodPost, "/api/users/{userSlug}/teams", restful.CreateTeamHandler, middleware.JWTTokenAuthMiddleware)
	register("teamsURL", http.MethodPost, "/api/users/{userSlug}/teams", restful.ListTeamHandler, middleware.JWTTokenAuthMiddleware)
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	for _, route := range routers {
		r := router.Methods(route.Method).Path(route.URI)
		if route.MiddleWare != nil {
			r.Handler(route.MiddleWare(route.Handler))
		} else {
			r.Handler(route.Handler)
		}
	}
	return router
}

func register(name, method, uri string, handler http.HandlerFunc, middleware mux.MiddlewareFunc) {
	routers = append(routers, Route{name, method, uri, handler, middleware})
}

func GetRouters() []Route {
	return routers
}
