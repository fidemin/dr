package main

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type MockDispatch struct{}

func (d *MockDispatch) Create(dispatch *Dispatch) error {
	dispatch.Id = 1
	now := time.Now().UTC()
	dispatch.CreatedAt = &now
	return nil
}

func (d *MockDispatch) List() ([]*Dispatch, error) {
	now := time.Now().UTC()
	d1 := &Dispatch{
		Id:        1,
		PaId:      1,
		Address:   "서울시 강남구",
		CreatedAt: &now,
	}
	d2 := &Dispatch{
		Id:          2,
		PaId:        1,
		DrId:        2,
		Address:     "서울시 강남구",
		IsComplete:  true,
		CreatedAt:   &now,
		CompletedAt: &now,
	}

	return []*Dispatch{d1, d2}, nil
}

func (d *MockDispatch) Get(dispatchId uint64) (*Dispatch, error) {
	now := time.Now().UTC()

	if dispatchId == 3 {
		// drId == paId case
		dispatch := &Dispatch{
			Id:         dispatchId,
			PaId:       2,
			Address:    "서울시 강남구",
			IsComplete: false,
			CreatedAt:  &now,
		}
		return dispatch, nil
	}

	dispatch := &Dispatch{
		Id:         dispatchId,
		PaId:       1,
		Address:    "서울시 강남구 1번지",
		IsComplete: false,
		CreatedAt:  &now,
	}
	return dispatch, nil
}

func (d *MockDispatch) Accept(dispatchId uint64, drId uint64) (*Dispatch, error) {
	if dispatchId != 1 {
		return nil, nil
	}

	now := time.Now().UTC()
	d1 := &Dispatch{
		Id:          dispatchId,
		PaId:        1,
		DrId:        drId,
		Address:     "서울시 강남구",
		IsComplete:  true,
		CreatedAt:   &now,
		CompletedAt: &now,
	}
	return d1, nil
}

func TestDispatchCreate(t *testing.T) {
	assert := assert.New(t)

	tEcho := echo.New()
	tEcho.Logger.SetLevel(1)
	jwtConfig := middleware.JWTConfig{
		Claims:     &JWTClaims{},
		SigningKey: []byte(TEST_SECRET),
	}
	tEcho.Use(middleware.JWTWithConfig(jwtConfig))
	tDispatchAPI := NewDispatchAPI(new(MockDispatch), testConfig)

	var testdata = []struct {
		body DispatchCreate
		code int
	}{
		{
			DispatchCreate{
				Address: "서울시 마포구 성산동 11-1",
			},
			201,
		},
		// wrong address
		{
			DispatchCreate{
				Address: "그들의 인간에 있는 것이다. 열매를 인생에 청춘에서만 봄바람이다. 품으며, 몸이 하는 있다. 그들의 이 인간의 청춘을 불어 싸인 있으랴? 동산에는 찬미를 속잎나고, 청춘은 것이다. 가는 예수는 있는 것이다. 생생하며, 반짝이는 생의 불어 쓸쓸하랴? 끓는 구하지 커다란 소금이라 스며들어 사라지지 인생을 살 못할 아니다. 인도하겠다는 간에 미인을 것이다.보라, 그리하였는가? 무엇이 이상은 것은 넣는 것이 끓는 인류의 살 주는 힘있다. 눈에 장식하는 그것은 이상 것이다.",
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
		c := tEcho.NewContext(req, rec)
		c.SetPath("/dispatch")
		MockSignIn(c, UsertypePassenger)
		if assert.NoError(tDispatchAPI.Create(c)) {
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

func TestDispatchList(t *testing.T) {
	assert := assert.New(t)

	tEcho := echo.New()
	tEcho.Logger.SetLevel(1)
	jwtConfig := middleware.JWTConfig{
		Claims:     &JWTClaims{},
		SigningKey: []byte(TEST_SECRET),
	}
	tEcho.Use(middleware.JWTWithConfig(jwtConfig))
	tDispatchAPI := NewDispatchAPI(new(MockDispatch), testConfig)

	var testdata = []struct {
		usertype Usertype
		code     int
	}{
		{
			UsertypeDriver,
			200,
		},
		{
			UsertypePassenger,
			200,
		},
	}

	for _, test := range testdata {
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := tEcho.NewContext(req, rec)
		c.SetPath("/dispatch")
		MockSignIn(c, test.usertype)
		if assert.NoError(tDispatchAPI.List(c)) {
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

func TestDispatchAccept(t *testing.T) {
	assert := assert.New(t)

	tEcho := echo.New()
	tEcho.Logger.SetLevel(1)
	jwtConfig := middleware.JWTConfig{
		Claims:     &JWTClaims{},
		SigningKey: []byte(TEST_SECRET),
	}
	tEcho.Use(middleware.JWTWithConfig(jwtConfig))
	tDispatchAPI := NewDispatchAPI(new(MockDispatch), testConfig)

	var testdata = []struct {
		usertype   Usertype
		dispatchId uint64
		code       int
	}{
		{
			UsertypeDriver,
			1,
			200,
		},
		{
			UsertypeDriver,
			2,
			400,
		},
		// driver id 와 passenger is 가 같을 경우
		{
			UsertypeDriver,
			3,
			400,
		},
		{
			UsertypePassenger,
			1,
			400,
		},
	}

	for _, test := range testdata {
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := tEcho.NewContext(req, rec)
		c.SetPath("/dispatch/:dispatch_id/accept")
		c.SetParamNames("dispatch_id")
		idStr := strconv.FormatUint(test.dispatchId, 10)
		c.SetParamValues(idStr)
		MockSignIn(c, test.usertype)
		if assert.NoError(tDispatchAPI.Accept(c)) {
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
