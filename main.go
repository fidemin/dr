package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"time"
)

type Config struct {
	Port   int64  `json:"port"`
	Debug  bool   `json:"debug"`
	Secret string `json:"secret"`

	DB DB `json:"db"`
}

type DB struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Name string `json:"name"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// config 파일 위치 받기
	var configFile string
	flag.StringVar(&configFile, "config", "", "config file path")
	flag.Parse()

	// config 객체에 설정 파일 파싱
	config := Config{}

	if err := ParseJSONFile(configFile, &config); err != nil {
		log.Fatal("parse config file error: ", err.Error())
	}

	// API DB 연결 및DB 풀생성
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=UTC",
		config.DB.User,
		config.DB.Pass,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name)

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal("database open error: ", err.Error())
	}

	defer db.Close()

	e := echo.New()
	// debugging 설정
	if config.Debug {
		e.Logger.SetLevel(1)
		db.SetLogger(e.Logger)
	} else {
		e.Logger.SetLevel(2)
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	apiConfig := map[string]string{
		"secret": config.Secret,
	}

	// auth API
	user := &DBUser{db: db}
	userAPI := NewUserAPI(user, apiConfig)
	authR := e.Group("/auth")
	authR.POST("/signup", userAPI.SignUp)
	authR.POST("/signin", userAPI.SignIn)

	// 배차 API
	dispatch := &DBDispatch{db: db}
	dispatchAPI := NewDispatchAPI(dispatch, apiConfig)

	jwtConfig := middleware.JWTConfig{
		Claims:     &JWTClaims{},
		SigningKey: []byte(config.Secret),
	}

	dispatchR := e.Group("/dispatch")
	dispatchR.Use(middleware.JWTWithConfig(jwtConfig))
	dispatchR.POST("", dispatchAPI.Create)
	dispatchR.GET("", dispatchAPI.List)
	dispatchR.POST("/:dispatch_id/accept", dispatchAPI.Accept)

	// run
	port := strconv.FormatInt(config.Port, 10)
	go func() {
		if err := e.Start(":" + port); err != nil {
			e.Logger.Info("shutting down server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
