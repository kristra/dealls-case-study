// dto/payslip_response.go
package dto

type AttendanceBreakdownItem struct {
	Date string `json:"date"`
}

type OvertimeBreakdownItem struct {
	Date        string  `json:"date"`
	HoursWorked float64 `json:"hours_worked"`
}

type ReimbursementBreakdownItem struct {
	Date        string  `json:"date"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

type PayslipResponse struct {
	ID     uint `json:"id"`
	Month  int  `json:"month"`
	Year   int  `json:"year"`
	UserID uint `json:"user_id"`

	// summary totals
	BaseSalary    float64 `json:"base_salary"`
	OvertimePay   float64 `json:"overtime_pay"`
	Reimbursement float64 `json:"reimbursement"`
	TotalSalary   float64 `json:"total_salary"`

	// calculation context
	MonthlySalary       float64 `json:"monthly_salary"`
	ExpectedWorkingDays int     `json:"expected_working_days"`
	DaysAttended        int     `json:"days_attended"`
	HourlyRate          float64 `json:"hourly_rate"`
	OvertimeRatePerHour float64 `json:"overtime_rate_per_hour"`

	// breakdowns
	TotalHoursWorked       float64                      `json:"total_hours_worked"`
	TotalOvertimeHours     float64                      `json:"total_overtime_hours"`
	AttendanceBreakdown    []AttendanceBreakdownItem    `json:"attendance_breakdown"`
	OvertimeBreakdown      []OvertimeBreakdownItem      `json:"overtime_breakdown"`
	ReimbursementBreakdown []ReimbursementBreakdownItem `json:"reimbursement_breakdown"`
}
