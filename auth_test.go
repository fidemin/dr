package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

const TEST_SECRET = "SECRET"

var testConfig = map[string]string{
	"secret": TEST_SECRET,
}

type MockUser struct{}

func (u *MockUser) Create(email string, password string, usertype Usertype) (*User, error) {
	user := new(User)
	user.Id = 2
	user.Email = email
	user.Password = password
	user.Usertype = usertype
	now := time.Now().UTC()
	user.CreatedAt = &now

	return user, nil
}

func (u *MockUser) Get(email string) (*User, error) {
	if email == "pa@gmail.com" {
		hash, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		return &User{
			Id:       1,
			Email:    "pa@gmail.com",
			Password: string(hash[:len(hash)]),
			Usertype: "PA",
		}, nil
	} else if email == "dr@gmail.com" {
		hash, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		return &User{
			Id:       2,
			Email:    "dr@gmail.com",
			Password: string(hash[:len(hash)]),
			Usertype: "DR",
		}, nil
	}
	return nil, nil
}

func MockSignIn(c echo.Context, usertype Usertype) {
	// jwtToken generation
	claimsNew := &JWTClaims{
		UserId:   2,
		Usertype: string(usertype),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	tokenNew := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsNew)

	t, err := tokenNew.SignedString([]byte(TEST_SECRET))
	if err != nil {
		fmt.Printf("[test] %s\n", err.Error())
		panic(err)
	}

	authToken := "Bearer " + t

	auth := strings.Split(authToken, " ")[1]

	// config should be same as used in init()
	config := middleware.JWTConfig{
		Claims:     &JWTClaims{},
		SigningKey: []byte(TEST_SECRET),
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return config.SigningKey, nil
	}

	type_ := reflect.ValueOf(config.Claims).Type().Elem()
	claims := reflect.New(type_).Interface().(jwt.Claims)
	token, _ := jwt.ParseWithClaims(auth, claims, keyFunc)
	c.Set("user", token)
}

func TestUserSignup(t *testing.T) {
	assert := assert.New(t)
	tEchoNew := echo.New()
	tEchoNew.Logger.SetLevel(1)
	tUserAPI := NewUserAPI(new(MockUser), testConfig)

	var testdata = []struct {
		body UserSignUp
		code int
	}{
		{
			UserSignUp{
				Email:    "test@gmail.com",
				Password: "password123",
				Usertype: "PA",
			},
			201,
		},
		// wrong user type
		{
			UserSignUp{
				Email:    "test@gmail.com",
				Password: "password123",
				Usertype: "PA1",
			},
			400,
		},
		// wrong email address
		{
			UserSignUp{
				Email:    "testgmail.com",
				Password: "password123",
				Usertype: "BA",
			},
			400,
		},
		// duplicate user
		{
			UserSignUp{
				Email:    "pa@gmail.com",
				Password: "password123",
				Usertype: "PA",
			},
			400,
		},
	}

	for _, test := range testdata {
		d, _ := json.Marshal(test.body)
		byteData := bytes.NewBuffer(d)
		req := httptest.NewRequest(echo.POST, "/", byteData)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := tEchoNew.NewContext(req, rec)
		c.SetPath("/auth/signup")
		if assert.NoError(tUserAPI.SignUp(c)) {
			var respData map[string]interface{}
			assert.Equal(test.code, rec.Code)
			t.Log(c.Path())
			t.Log(rec.Body.String())
			if err := json.Unmarshal(rec.Body.Bytes(), &respData); err != nil {
				assert.Fail(err.Error())
			}
		}
	}

}

func TestUserLogin(t *testing.T) {
	assert := assert.New(t)
	tEchoNew := echo.New()
	tEchoNew.Logger.SetLevel(1)
	tUserAPI := NewUserAPI(new(MockUser), testConfig)

	var testdata = []struct {
		body UserSignIn
		code int
	}{
		{
			UserSignIn{
				Email:    "pa@gmail.com",
				Password: "123456",
			},
			200,
		},
		// user가 없는 경우
		{
			UserSignIn{
				Email:    "test@gmail.com",
				Password: "123456",
			},
			401,
		},
		// password가 틀린경우
		{
			UserSignIn{
				Email:    "pa@gmail.com",
				Password: "abcdef",
			},
			401,
		},
	}

	for _, test := range testdata {
		d, _ := json.Marshal(test.body)
		byteData := bytes.NewBuffer(d)
		req := httptest.NewRequest(echo.POST, "/", byteData)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := tEchoNew.NewContext(req, rec)
		c.SetPath("/auth/signin")
		if assert.NoError(tUserAPI.SignIn(c)) {
			var respData map[string]interface{}
			assert.Equal(test.code, rec.Code)
			t.Log(c.Path())
			t.Log(rec.Body.String())
			if err := json.Unmarshal(rec.Body.Bytes(), &respData); err != nil {
				assert.Fail(err.Error())
			}
		}
	}
}
