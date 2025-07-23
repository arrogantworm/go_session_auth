package models

import "time"

type (
	Session struct {
		Token     string
		Username  string
		UserId    int
		ExpiresAt time.Time
	}
)
