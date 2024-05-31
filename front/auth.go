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
	Id    int64  `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func verifySignin(credential string) (string, error) {

	payload, err := idtoken.Validate(context.Background(), credential, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		return "", err
	}
	email, ok := payload.Claims["email"].(string)
	if !ok {
		return "", fmt.Errorf("error: email not in payload claims")
	}
	return email, nil
}

func jwtEncodeEmailId(email string, id int64) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{Email: email, Id: id}).SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func loggedInEmail(req *http.Request) (bool, string, int64) {
	jwtCookie, err := req.Cookie("jwt")
	if err != nil {
		log.Println("loggedInEmail no cookie:", err)
		return false, "", 0
	}
	parsedToken, err := jwt.ParseWithClaims(jwtCookie.Value, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		log.Println("loggedInEmail parse jwt:", err)
		return false, "", 0
	}
	claims := parsedToken.Claims.(*UserClaims)
	return true, claims.Email, claims.Id
}
