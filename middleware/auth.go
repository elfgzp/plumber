package middleware

import (
	"fmt"
	"github.com/elfgzp/plumber/helpers"
	"log"
	"net/http"
)

func JWTTokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("authorization")
		if tokenStr == "" {
			helpers.ResponseWithJSON(w, http.StatusUnauthorized, helpers.UnauthorizedResponse())
		} else {
			email, err := helpers.CheckToken(tokenStr)
			if err != nil {
				helpers.ResponseWithJSON(w, http.StatusUnauthorized, helpers.UnauthorizedResponse())
			} else {
				log.Println(fmt.Sprintf("%s %s", email, r.RemoteAddr, r.URL))
				next.ServeHTTP(w, r)
			}
		}
	})
}
