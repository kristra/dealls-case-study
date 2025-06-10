package models

import (
	"time"

	"gorm.io/gorm"
)

type Attendance struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	User        User `gorm:"foreignKey:UserID"`
	Date        time.Time
	CheckInAt   *time.Time `gorm:"default:null"`
	CheckOutAt  *time.Time `gorm:"default:null"`
	HoursWorked float64    `gorm:"not null"`
	CreatedAt   time.Time
	CreatedBy   uint
	UpdatedAt   time.Time
	UpdatedBy   uint
}

func (a *Attendance) DateOnlyString() string {
	return a.Date.Format("2006-01-02")
}

func (a *Attendance) BeforeSave(tx *gorm.DB) (err error) {
	if a.CheckInAt != nil && a.CheckOutAt != nil {
		duration := a.CheckOutAt.Sub(*a.CheckInAt).Hours()
		if duration < 0 {
			duration = 0
		}
		a.HoursWorked = duration
	}
	return
}
