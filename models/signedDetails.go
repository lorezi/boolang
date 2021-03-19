package models

import "github.com/dgrijalva/jwt-go"

// SignedDetails struct
type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UID       string
	jwt.StandardClaims
}
