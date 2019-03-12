package main

import (
	"time"
)

// 유저 타입 PA: 승객 DR: 택시기사
type Usertype string

const (
	UsertypePassenger Usertype = "PA"
	UsertypeDriver             = "DR"
)

// User ORM 객체
type User struct {
	Id        uint64     `gorm:"AUTO_INCREMENT"`
	Email     string     `gorm:"column:email" validate:"regexp=^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$"`
	Password  string     `gorm:"column:password"`
	Usertype  Usertype   `gorm:"column:usertype" validate:"regexp=^(DR|PA)$"`
	CreatedAt *time.Time `gorm:"column:created_at"`
}

func (User) TableName() string {
	return "user"
}

// User SignUp request data 객체
type UserSignUp struct {
	Email    string `json:"email" validate:"regexp=^[a-z0-9_%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$"`
	Password string `json:"password" validate:"nonzero"`
	Usertype string `json:"usertype" validate:"regexp=^(DR|PA)$"`
}

// User SignIn request data 객체
type UserSignIn struct {
	Email    string `json:"email" validate:"regexp=^[a-z0-9_%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$"`
	Password string `json:"password" validate:"nonzero"`
}
