package models

import (
	"time"
)

type Payroll struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	PeriodStart time.Time
	PeriodEnd   time.Time
	Status      string
	GeneratedAt time.Time
	CreatedAt   time.Time
	CreatedBy   uint
	UpdatedAt   time.Time
	UpdatedBy   uint

	Payslips []Payslip `gorm:"foreignKey:PayrollID"`
}
