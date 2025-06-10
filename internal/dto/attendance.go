package dto

import "time"

type AttendanceResponse struct {
	ID         uint       `json:"id"`
	Date       time.Time  `json:"date"`
	CheckInAt  *time.Time `json:"check_in_at,omitempty"`
	CheckOutAt *time.Time `json:"check_out_at,omitempty"`
}
