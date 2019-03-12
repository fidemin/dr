package main

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type JWTClaims struct {
	UserId   uint64 `json:"user_id"`
	Usertype string `json:"usertype"`
	jwt.StandardClaims
}

// JWT 토큰을 가져온다.
func GetJWTClaims(c echo.Context) JWTClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	return *claims
}
