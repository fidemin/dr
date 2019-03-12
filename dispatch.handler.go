package main

import (
	"github.com/labstack/echo"
	"strconv"
)

// 배차 API 객체
type DispatchAPI struct {
	dispatch DispatchInterface
	config   map[string]string
}

func NewDispatchAPI(dispatch DispatchInterface, config map[string]string) *DispatchAPI {
	api := new(DispatchAPI)
	api.dispatch = dispatch
	api.config = config
	return api
}

// POST /dispatch
// 배차를 생성한다.
func (d *DispatchAPI) Create(c echo.Context) error {
	// JWT 토큰을 가져온다.
	claims := GetJWTClaims(c)

	dc := DispatchCreate{}
	// request data 바인딩
	if err := c.Bind(&dc); err != nil {
		c.Logger().Warn(err)
		return c.JSON(400, echo.Map{
			"message": "bad request data",
		})
	}

	// 100자가 넘으면 안된다. validator 라이브러리가 유니코드에 제대로 작동을 안해 직접 구현
	length := len([]rune(dc.Address))
	if length == 0 || length > 100 {
		c.Logger().Warn("address length limit exceeded")
		return c.JSON(400, echo.Map{
			"message": "bad request validation failed",
		})
	}

	dispatch := new(Dispatch)
	dispatch.PaId = claims.UserId
	dispatch.Address = dc.Address

	if err := d.dispatch.Create(dispatch); err != nil {
		c.Logger().Error(err)
		return c.JSON(500, echo.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(201, echo.Map{
		"message": "success",
		"data":    dispatch,
	})
}

// GET /dispatch
// 배차 목록 API
func (d *DispatchAPI) List(c echo.Context) error {
	GetJWTClaims(c)

	data, err := d.dispatch.List()
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(500, echo.Map{
			"message": "Internal Server Error",
		})
	}
	return c.JSON(200, echo.Map{
		"message": "success",
		"data":    data,
	})
}

// POST /dispatch/:dispatch_id/accept
// 배차 완료 API
func (d *DispatchAPI) Accept(c echo.Context) error {
	claims := GetJWTClaims(c)

	if claims.Usertype != "DR" {
		msg := "user is not a driver"
		c.Logger().Warn(msg)
		return c.JSON(400, echo.Map{
			"message": msg,
		})
	}

	id, err := strconv.ParseUint(c.Param("dispatch_id"), 10, 64)

	if err != nil {
		msg := "Not Found"
		c.Logger().Warn(msg)
		return c.JSON(404, echo.Map{
			"message": msg,
		})
	}

	preDispatch, err := d.dispatch.Get(id)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(500, echo.Map{
			"message": "Internal Server Error",
		})
	}

	if preDispatch == nil {
		msg := "the dispatch does not exist"
		c.Logger().Warn(msg)
		return c.JSON(400, echo.Map{
			"message": msg,
		})
	}

	if preDispatch.PaId == claims.UserId {
		msg := "the driver id and passenger id is same"
		c.Logger().Warn(msg)
		return c.JSON(400, echo.Map{
			"message": msg,
		})
	}

	dispatch, err := d.dispatch.Accept(id, claims.UserId)

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(500, echo.Map{
			"message": "Internal Server Error",
		})
	}

	if dispatch == nil {
		msg := "the dispatch is already accepted"
		c.Logger().Warn(msg)
		return c.JSON(400, echo.Map{
			"message": msg,
		})
	}

	return c.JSON(200, echo.Map{
		"message": "success",
		"data":    dispatch,
	})

}
