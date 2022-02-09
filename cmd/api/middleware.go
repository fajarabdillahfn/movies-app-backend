package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
)

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Acecess-Control-Allow-Headers", "Content-Type,Authorization")
		next.ServeHTTP(w, r)
	})
}

func (app *application) checkToken(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authoriztion")

		if authHeader == "" {
			
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			app.logger.Println("invalid auth header")
			app.errorJSON(w, errors.New("invalid auth header"))
			return
		}

		if headerParts[0] != "Bearer" {
			app.logger.Println("unauthorized - no bearer")
			app.errorJSON(w, errors.New("unauthorized - no bearer"))
			return
		}

		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))
		if err != nil {
			app.logger.Println("unauthorized - failed hmac check")
			app.errorJSON(w, errors.New("unauthorized - failed hmac check"))
			return
		}

		if !claims.Valid(time.Now()) {
			app.logger.Println("unauthorized - token expired")
			app.errorJSON(w, errors.New("unauthorized - token expired"))
			return
		}

		if !claims.AcceptAudience("mydomain.com") {
			app.logger.Println("unauthorized - invalid audienc")
			app.errorJSON(w, errors.New("unauthorized - invalid audience"))
			return
		}

		if claims.Issuer != "mydomain.com" {
			app.logger.Println("unauthorized - invalid issuer")
			app.errorJSON(w, errors.New("unauthorized - invalid issuer"))
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			app.logger.Println("unauthorized")
			app.errorJSON(w, errors.New("unauthorized"))
			return
		}

		app.logger.Println("valid user:", userID)

		next.ServeHTTP(w, r)
	})
}
