package handler

// Common structs

type (
	successRes struct {
		Message string `json:"success"`
	}
)

type (
	errorRes struct {
		Message string `json:"error"`
	}
)

// Auth

type (
	UserRes struct {
		Username string `json:"username"`
	}
)
