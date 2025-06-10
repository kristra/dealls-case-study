package models

import (
	"time"
)

type Overtime struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	User        User `gorm:"foreignKey:UserID"`
	Date        time.Time
	HoursWorked float64 `gorm:"not null"`
	CreatedAt   time.Time
	CreatedBy   uint
	UpdatedAt   time.Time
	UpdatedBy   uint
}

func (o *Overtime) DateOnlyString() string {
	return o.Date.Format("2006-01-02")
}
