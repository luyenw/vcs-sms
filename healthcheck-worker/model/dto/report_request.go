package dto

type ReportRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
}
