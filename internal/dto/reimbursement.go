package dto

type SubmitReimbursementRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description *string `json:"description,omitempty"`
}
