// dto/payslip_response.go
package dto

type PayslipResponse struct {
	ID                     uint    `json:"id"`
	Month                  int     `json:"month"`
	Year                   int     `json:"year"`
	UserID                 uint    `json:"user_id"`
	BaseSalary             float64 `json:"base_salary"`
	OvertimePay            float64 `json:"overtime_pay"`
	Reimbursement          float64 `json:"reimbursement"`
	TotalSalary            float64 `json:"total_salary"`
	TotalHoursWorked       float64 `json:"total_hours_worked"`
	TotalOvertimeHours     float64 `json:"total_overtime_hours"`
	AttendanceBreakdown    string  `json:"attendance_breakdown"`
	OvertimeBreakdown      string  `json:"overtime_breakdown"`
	ReimbursementBreakdown string  `json:"reimbursement_breakdown"`
}
