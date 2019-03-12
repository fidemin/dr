package main

import (
	"github.com/jinzhu/gorm"
	"time"
)

// 배차 로직을 수행하는 인터페이스 객체
type DispatchInterface interface {
	List() ([]*Dispatch, error)
	Get(dispatchId uint64) (*Dispatch, error)
	Create(dispatch *Dispatch) error
	Accept(dispatchId uint64, drId uint64) (*Dispatch, error)
}

// DB에 연결된 Dispatch 객체
type DBDispatch struct {
	db *gorm.DB
}

// 배차 생성
func (d *DBDispatch) Create(dispatch *Dispatch) error {
	if err := d.db.Create(&dispatch).Error; err != nil {
		return err
	}
	return nil
}

// 배차 목록 부르기
func (d *DBDispatch) List() ([]*Dispatch, error) {
	var dispatches []*Dispatch
	if err := d.db.Order("created_at desc").Find(&dispatches).Error; err != nil {
		return nil, err
	}

	return dispatches, nil
}

// 배차 하나의 정보를 가져오기
func (d *DBDispatch) Get(dispatchId uint64) (*Dispatch, error) {
	dispatch := &Dispatch{Id: dispatchId}
	if r := d.db.First(&dispatch); r.Error != nil {
		if r.RecordNotFound() {
			return nil, nil
		}
		return nil, r.Error
	}
	return dispatch, nil
}

// 배차를 기사가 수용
func (d *DBDispatch) Accept(dispatchId uint64, drId uint64) (*Dispatch, error) {
	if dispatchId == 0 {
		return nil, nil
	}
	if drId == 0 {
		return nil, nil
	}

	now := time.Now().UTC()
	updates := map[string]interface{}{
		"dr_id":        drId,
		"is_complete":  true,
		"completed_at": &now,
	}

	// 완료된 상태가 아니라면 업데이트한다.
	if r := d.db.Model(&Dispatch{}).Where("id = ? AND is_complete = ?", dispatchId, false).Updates(updates); r.Error != nil {
		return nil, r.Error
	} else {
		// 이미 완료되었는데, 클라이언트에서 시간차로 accept 호출을 할 수 있으므로, 이를 방지한다.
		if r.RowsAffected == 0 {
			return nil, nil
		}
	}

	dispatch := &Dispatch{Id: dispatchId}
	if err := d.db.First(&dispatch).Error; err != nil {
		return nil, err
	}
	return dispatch, nil
}
