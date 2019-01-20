package middleware

import (
	"fmt"
	"github.com/elfgzp/plumber/helpers"
	"github.com/gorilla/context"
	"log"
	"net/http"
)

func JWTTokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("authorization")
		if tokenStr == "" {
			helpers.Response401(w)
		} else {
			email := helpers.CheckToken(tokenStr)
			if email == "" {
				helpers.Response401(w)
			} else {
				log.Println(fmt.Sprintf("%s %s %s", email, r.RemoteAddr, r.URL.Path))
				context.Set(r, "userEmail", email)
				next.ServeHTTP(w, r)
			}
		}
	})
}
