package main

import (
	"github.com/jinzhu/gorm"
)

// User 로직을 실행하는 인터페이스
type UserInterface interface {
	Create(email string, password string, usertype Usertype) (*User, error)
	Get(email string) (*User, error)
}

// DB User 객체. 실제 DB와 연결되어 UserInterface의 비즈니스 로직을 구현
type DBUser struct {
	db *gorm.DB
}

// User 생성
func (u *DBUser) Create(email string, password string, usertype Usertype) (*User, error) {
	user := &User{
		Email:    email,
		Password: password,
		Usertype: usertype,
	}

	if err := u.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// User 데이터 불러오기
func (u *DBUser) Get(email string) (*User, error) {
	user := new(User)
	if r := u.db.Where("email = ?", email).First(&user); r.Error != nil {
		if r.RecordNotFound() {
			return nil, nil
		}
		return nil, r.Error
	}
	return user, nil
}
