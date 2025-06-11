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

	// summary totals
	BaseSalary    float64
	OvertimePay   float64
	Reimbursement float64
	TotalSalary   float64

	// calculation context
	MonthlySalary       float64 `gorm:"not null"`
	ExpectedWorkingDays int     `gorm:"not null"`
	DaysAttended        int     `gorm:"not null"`
	HourlyRate          float64 `gorm:"not null"`
	OvertimeRatePerHour float64 `gorm:"not null"`

	// breakdowns
	TotalHoursWorked       float64
	TotalOvertimeHours     float64
	AttendanceBreakdown    string `gorm:"type:text"`
	OvertimeBreakdown      string `gorm:"type:text"`
	ReimbursementBreakdown string `gorm:"type:text"`

	CreatedAt time.Time
	CreatedBy uint
	UpdatedAt time.Time
	UpdatedBy uint
}
