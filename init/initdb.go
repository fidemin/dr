package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func ParseJSONFile(filepath string, config interface{}) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return err
	}
	return nil
}

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

	// config 파일 위치 받기
	var configFile string
	flag.StringVar(&configFile, "config", "", "config file path")
	flag.Parse()

	// config 객체에 설정 파일 파싱
	config := Config{}

	if err := ParseJSONFile(configFile, &config); err != nil {
		log.Fatal("parse config file error: ", err.Error())
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=UTC",
		config.DB.User,
		config.DB.Pass,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name)

	db, err := sql.Open("mysql", dsn)

	file, err := os.Open("init/init.sql")
	if err != nil {
		log.Fatal("[init db]", err.Error())
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("[init db]", err.Error())
	}

	queries := strings.Split(string(b), ";")

	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		if _, err := db.Exec(query); err != nil {
			log.Fatal("[init db]", err.Error())
		}
	}
}
