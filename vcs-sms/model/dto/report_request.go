package dto

type ReportRequest struct {
	StartDate string `json:"start_date" binding:"required" validate:"date"`
	EndDate   string `json:"end_date" binding:"required" validate:"date"`
	Email     string `json:"email" binding:"required,email" validate:"email"`
}
