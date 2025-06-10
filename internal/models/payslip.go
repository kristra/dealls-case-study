package models

import (
	"time"
)

type Payslip struct {
	ID uint `gorm:"primaryKey"`

	UserID    uint
	User      User `gorm:"foreignKey:UserID"`
	PayrollID uint
	Payroll   Payroll `gorm:"foreignKey:PayrollID"`

	Month int `gorm:"not null"`
	Year  int `gorm:"not null"`

	BaseSalary         float64
	OvertimePay        float64
	Reimbursement      float64
	TotalSalary        float64
	TotalHoursWorked   float64
	TotalOvertimeHours float64

	AttendanceBreakdown    string `gorm:"type:text"`
	OvertimeBreakdown      string `gorm:"type:text"`
	ReimbursementBreakdown string `gorm:"type:text"`

	CreatedAt time.Time
	CreatedBy uint
	UpdatedAt time.Time
	UpdatedBy uint
}
