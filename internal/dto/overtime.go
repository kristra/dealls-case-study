package dto

type SubmitOvertimeRequest struct {
	HoursWorked float64 `json:"hours_worked" binding:"required,gt=0,lte=3"`
}

type SubmitOvertimeResponse struct {
	ID          uint    `json:"id"`
	HoursWorked float64 `json:"hours_worked"`
}
