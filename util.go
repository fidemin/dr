package main

import (
	"encoding/json"
	"os"
)

// json 형식의 config 파일을 go object에 파싱한다.
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
