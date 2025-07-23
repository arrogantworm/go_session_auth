package repository

import (
	"errors"
	"session-auth/models"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type Crypt struct{}

func NewHasher() *Crypt {
	return &Crypt{}
}

func (c Crypt) HashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (c Crypt) CheckPassword(hashPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}

type SessionManager struct {
	sessionTTL time.Duration
}

func NewSessionManager() (*SessionManager, error) {
	sessionTTL := viper.GetDuration("sessions.TTL")
	if sessionTTL == 0 {
		return nil, errors.New("session ttl not stated in config.yaml")
	}

	return &SessionManager{
		sessionTTL: sessionTTL,
	}, nil
}

func (s *SessionManager) CreateSession(username string) *models.Session {
	return &models.Session{
		Token:     uuid.NewString(),
		Username:  username,
		ExpiresAt: time.Now().Add(s.sessionTTL),
	}
}
