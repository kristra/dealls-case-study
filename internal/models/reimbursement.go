package models

import (
	"time"
)

type Reimbursement struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	User        User `gorm:"foreignKey:UserID"`
	Date        time.Time
	Amount      float64 `gorm:"not null"`
	Description string
	CreatedBy   uint
	CreatedAt   time.Time
	UpdatedBy   uint
	UpdatedAt   time.Time
}

func (r *Reimbursement) DateOnlyString() string {
	return r.Date.Format("2006-01-02")
}
