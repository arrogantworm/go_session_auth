package models

type (
	Credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	User struct {
		Id           int
		Username     string
		PasswordHash string
	}
)
