package router

import (
	"github.com/elfgzp/plumber/controllers/restful"
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Method     string
	URI        string
	Handler    http.HandlerFunc
	MiddleWare mux.MiddlewareFunc
}

var routes []Route

func init() {
	register(http.MethodPost, "/api/token", restful.CreateToken, nil)
	register(http.MethodPost, "/api/token/verification", restful.TokenVerify, nil)
	register(http.MethodPost, "/api/users", restful.CreateUser, nil)
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	for _, route := range routes {
		r := router.Methods(route.Method).Path(route.URI)
		if route.MiddleWare != nil {
			r.Handler(route.MiddleWare(route.Handler))
		} else {
			r.Handler(route.Handler)
		}
	}
	return router
}

func register(method, uri string, handler http.HandlerFunc, middleware mux.MiddlewareFunc) {
	routes = append(routes, Route{method, uri, handler, middleware})
}
