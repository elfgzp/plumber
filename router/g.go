package router

import (
	"fmt"
	"net/http"

	"github.com/elfgzp/plumber/controllers/restful"
	"github.com/elfgzp/plumber/helpers"
	"github.com/elfgzp/plumber/middleware"
	"github.com/gorilla/mux"
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
		helpers.Response200(w, "", apiURL)
	}, nil)

	register("token", http.MethodPost, "/api/token", restful.CreateTokenHandler, nil)
	register("tokenVerification", http.MethodPost, "/api/token/verification", restful.TokenVerifyHandler, nil)

	register("currentUser", http.MethodGet, "/api/users/current", restful.GetCurrentUser, middleware.JWTTokenAuthMiddleware)
	register("users", http.MethodPost, "/api/users", restful.CreateUserHandler, nil)

	register("userTeams", http.MethodPost, "/api/users/{userSlug}/teams", restful.CreateTeamHandler, middleware.JWTTokenAuthMiddleware)
	register("userTeams", http.MethodGet, "/api/users/{userSlug}/teams", restful.ListTeamHandler, middleware.JWTTokenAuthMiddleware)
	register("team", http.MethodGet, "/api/teams/{teamSlug}", restful.RetrieveTeamHandler, middleware.JWTTokenAuthMiddleware)

	register("teamProjects", http.MethodPost, "/api/teams/{teamSlug}/projects", restful.CreateProjectHandler, middleware.JWTTokenAuthMiddleware)
	register("teamProjects", http.MethodGet, "/api/teams/{teamSlug}/projects", restful.ListProjectHandler, middleware.JWTTokenAuthMiddleware)
	register("project", http.MethodGet, "/api/projects/{projectSlug}", restful.RetrieveProjectHandler, middleware.JWTTokenAuthMiddleware)

	register("projectTaskLists", http.MethodPost, "/api/projects/{projectSlug}/task-lists", restful.CreateTaskListHandler, middleware.JWTTokenAuthMiddleware)
	register("projectTaskLists", http.MethodGet, "/api/projects/{projectSlug}/task-lists", restful.ListTaskListHandler, middleware.JWTTokenAuthMiddleware)
	register("projectTaskList", http.MethodGet, "/api/projects/{projectSlug}/task-lists/{taskSlug}", restful.RetrieveTaskListHandler, middleware.JWTTokenAuthMiddleware)
	register("projectTaskList", http.MethodPut, "/api/projects/{projectSlug}/task-lists/{taskSlug}", restful.UpdateTaskListHandler, middleware.JWTTokenAuthMiddleware)

	register("projectTaskList", http.MethodDelete, "/api/projects/{projectSlug}/task-lists/{taskSlug}", restful.DestroyTaskListHandler, middleware.JWTTokenAuthMiddleware)

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
