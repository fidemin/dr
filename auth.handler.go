package main

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"
	"time"
)

// UserAPI 객체.
type UserAPI struct {
	user   UserInterface
	config map[string]string
}

func NewUserAPI(user UserInterface, config map[string]string) *UserAPI {
	uAPI := new(UserAPI)
	uAPI.user = user
	uAPI.config = config

	return uAPI
}

// POST /auth/signup
func (u *UserAPI) SignUp(c echo.Context) error {
	var us UserSignUp

	// body data biding
	if err := c.Bind(&us); err != nil {
		c.Logger().Warn(err)
		return c.JSON(400, echo.Map{
			"message": "bad request data",
		})
	}

	// ruquest data validation
	if err := validator.Validate(us); err != nil {
		c.Logger().Warn(err)
		return c.JSON(400, echo.Map{
			"message": "request data validation failed",
		})
	}

	// 이메일 중복 체크
	dupUser, err := u.user.Get(us.Email)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(500, echo.Map{
			"message": "Internal Server Error",
		})
	}

	if dupUser != nil {
		c.Logger().Warn(err)
		return c.JSON(400, echo.Map{
			"message": "user already exists",
		})
	}

	// 패스워드 생성
	hash, err := bcrypt.GenerateFromPassword([]byte(us.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(500, echo.Map{
			"message": "Internal Server Error",
		})
	}

	// 사용자 생성
	if _, err := u.user.Create(us.Email, string(hash[:len(hash)]), Usertype(us.Usertype)); err != nil {
		c.Logger().Error(err)
		return c.JSON(500, echo.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(201, echo.Map{
		"message": "suscess",
	})
}

// POST /auth/signin
func (u *UserAPI) SignIn(c echo.Context) error {

	// body data biding
	us := UserSignIn{}
	if err := c.Bind(&us); err != nil {
		c.Logger().Warn(err)
		return c.JSON(400, echo.Map{
			"message": "bad request data",
		})
	}

	// request data validation
	if err := validator.Validate(us); err != nil {
		c.Logger().Warn(err)
		return c.JSON(400, echo.Map{
			"message": "request data validation failed",
		})
	}

	// find user with email address
	user, err := u.user.Get(us.Email)

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(500, echo.Map{
			"message": "Internal Server Error",
		})
	}

	if user == nil {
		c.Logger().Warn(fmt.Sprintf("user with %s not found", us.Email))
		return c.JSON(401, echo.Map{
			"message": "Wrong username or password",
		})
	}

	// 패스워드 체크
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(us.Password)); err != nil {
		c.Logger().Warn(err)
		return c.JSON(401, echo.Map{
			"message": "Wrong username or password",
		})
	}

	// JWT Token 생성
	claims := &JWTClaims{
		UserId:   user.Id,
		Usertype: string(user.Usertype),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// JWT 토큰 암호화 (해싱)
	t, err := token.SignedString([]byte(u.config["secret"]))

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(500, echo.Map{
			"message": "Internal Server Error",
		})
	}

	// JWT token 응답으로 보냄
	return c.JSON(200, echo.Map{
		"message": "success",
		"data": map[string]string{
			"token": t,
		},
	})
}
