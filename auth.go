// Handles a simple SignIn with Google & JWT auth
// User presses SignIn w Google button, google sends token to /login,
// we validate it and manage the session with JWT.

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"google.golang.org/api/idtoken"
)

type UserClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func login(w http.ResponseWriter, req *http.Request) {

	// Verify the Signin with google token
	payload, err := idtoken.Validate(context.Background(), req.FormValue("credential"), os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		log.Println(err)
		// TODO error code
		fmt.Fprint(w, "Could not authenticate")
		return
	}
	email, ok := payload.Claims["email"].(string)
	if !ok {
		log.Println("login email err:", err)
		fmt.Fprint(w, "Could not authenticate")
		return
	}

	// Generate JWT, place in cookie "jwt", and redirect to prefs
	jwtString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{Email: email}).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Println("login jwt err:", err)
		fmt.Fprint(w, "Could not authenticate")
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "jwt", Value: jwtString})
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func logout(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "jwt", Value: ""})
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func loggedInEmail(req *http.Request) (bool, string) {
	jwtCookie, err := req.Cookie("jwt")
	if err != nil {
		log.Println("loggedInEmail no cookie:", err)
		return false, ""
	}
	parsedToken, err := jwt.ParseWithClaims(jwtCookie.Value, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		log.Println("loggedInEmail parse jwt:", err)
		return false, ""
	}
	claims := parsedToken.Claims.(*UserClaims)
	return true, claims.Email
}
