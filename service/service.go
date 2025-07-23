package service

import (
	"context"
	"session-auth/models"
	"session-auth/repository"

	"github.com/google/uuid"
)

type Database interface {
	Close(ctx context.Context)
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByUsername(ctx context.Context, user *models.User) error
	SaveSession(ctx context.Context, session *models.Session) error
	GetSessionByID(ctx context.Context, session_id uuid.UUID) (*models.Session, error)
	DeleteSessionByID(ctx context.Context, session_id uuid.UUID) error
}

type Authentication interface {
	HashPassword(password string) (string, error)
	CheckPassword(hashPassword, password string) bool
	CreateSession(username string) *models.Session
}

type Service struct {
	Database
	Authentication
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authentication: repo.Authentication,
		Database:       repo.Database,
	}
}
