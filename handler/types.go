package handler

// Common structs

type (
	SuccessRes struct {
		Message string `json:"success"`
	}
)

type (
	ErrorRes struct {
		Message string `json:"error"`
	}
)

// Auth

type (
	UserRes struct {
		Username string `json:"username"`
	}
)
