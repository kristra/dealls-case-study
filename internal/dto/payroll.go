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
