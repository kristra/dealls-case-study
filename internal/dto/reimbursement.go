package dto

type SubmitReimbursementRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description *string `json:"description,omitempty"`
}

type SubmitReimbursementResponse struct {
	ID          uint    `json:"id"`
	Amount      float64 `json:"amount"`
	Description *string `json:"description,omitempty"`
}
