package main

import (
	"time"
)

// 배차 정보 객체
type Dispatch struct {
	Id          uint64     `json:"id" gorm:"AUTO_INCREMENT"`
	PaId        uint64     `json:"pa_id" gorm:"column:pa_id"` //
	DrId        uint64     `json:"dr_id,omitempty" gorm:"column:dr_id;default:null"`
	Address     string     `json:"address" gorm:"column:address"`
	IsComplete  bool       `json:"is_complete" gorm:"is_complete" colum:is_complete"`
	CreatedAt   *time.Time `json:"created_at" gorm:"column:created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" gorm:"column:completed_at"`
}

func (Dispatch) TableName() string {
	return "dispatch"
}

// 배차 생성 API reqeust body 객체
type DispatchCreate struct {
	Address string `json:"address"`
}
