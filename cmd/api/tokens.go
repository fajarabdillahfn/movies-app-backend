package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"movies-app-backend/models"
	"net/http"
	"time"

	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
)

var validUser = models.User{
	ID:       10,
	Email:    "me@here.com",
	Password: "$2a$12$Sxit9n9u.FeuUtS16frpLuGPyHBWvtBsPpG2rmK4XMmB3MqHRN7eO",
}

type Credential struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

func (app *application) signIn(w http.ResponseWriter, r *http.Request) {
	var cred Credential

	err := json.NewDecoder(r.Body).Decode(&cred)

	if err != nil {
		app.logger.Print(err)
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}

	hashedPassword := validUser.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(cred.Password))
	if err != nil {
		app.logger.Print(err)
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}

	var claims jwt.Claims
	 claims.Subject = fmt.Sprint(validUser.ID)
	 claims.Issued = jwt.NewNumericTime(time.Now())
	 claims.NotBefore = jwt.NewNumericTime(time.Now())
	 claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	 claims.Issuer = "mydomain.com"
	 claims.Audiences = []string{"mydomain.com"}

	 jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	 if err != nil {
		app.logger.Print(err)
		app.errorJSON(w, errors.New("error signing"))
		return
	}

	app.writeJSON(w, http.StatusOK, string(jwtBytes), "response ")
}
