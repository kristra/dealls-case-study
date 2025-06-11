package dto

import "time"

type UpsertPayrollRequest struct {
	Name        *string    `json:"name,omitempty"`
	PeriodStart *time.Time `json:"period_start,omitempty"`
	PeriodEnd   *time.Time `json:"period_end,omitempty"`
}

type PayrollResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	PeriodStart *time.Time `json:"period_start,omitempty"`
	PeriodEnd   *time.Time `json:"period_end,omitempty"`
	Status      string     `json:"status"`
}

type PayrollSummaryResponse struct {
	PayrollID     uint                   `json:"payroll_id"`
	Year          int                    `json:"year"`
	Month         int                    `json:"month"`
	TotalSalaries float64                `json:"total_salaries"`
	Payslips      []EmployeePayslipBrief `json:"payslips"`
}

type EmployeePayslipBrief struct {
	UserID        uint    `json:"user_id"`
	Username      string  `json:"username"`
	BaseSalary    float64 `json:"base_salary"`
	OvertimePay   float64 `json:"overtime_pay"`
	Reimbursement float64 `json:"reimbursement"`
	TotalPay      float64 `json:"total_pay"`
}
