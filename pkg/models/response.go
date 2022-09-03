package models

import "time"

type OKResponse struct {
	Error bool        `json:"error" example:"false"`
	Data  interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error       bool      `json:"error" example:"true"`
	Errno       string    `json:"code"  example:"123"`
	Message     string    `json:"msg"`
	Timestamp   time.Time `json:"timestamp"`
	FullMessage string    `json:"full_msg"`
}
