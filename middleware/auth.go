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
			helpers.ResponseWithJSON(w, http.StatusUnauthorized, helpers.UnauthorizedResponse())
		} else {
			email := helpers.CheckToken(tokenStr)
			if email == "" {
				helpers.ResponseWithJSON(w, http.StatusUnauthorized, helpers.UnauthorizedResponse())
			} else {
				log.Println(fmt.Sprintf("%s %s", email, r.RemoteAddr, r.URL))
				context.Set(r, "userEmail", email)
				next.ServeHTTP(w, r)
			}
		}
	})
}
