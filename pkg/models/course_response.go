package models

type CourseAllResponse struct {
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	ID        int       `json:"id"`
	Score     []float64 `json:"score"`
	Institute string    `json:"institute"`
}
