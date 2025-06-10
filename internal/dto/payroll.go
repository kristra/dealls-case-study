package dto

import "time"

type UpsertPayrollRequest struct {
	Name        *string    `json:"name,omitempty"`
	PeriodStart *time.Time `json:"period_start,omitempty"`
	PeriodEnd   *time.Time `json:"period_end,omitempty"`
}
