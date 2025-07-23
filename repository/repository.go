package repository

import (
	"context"
	"session-auth/models"

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

type Hasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(hashPassword, password string) bool
}

type Session interface {
	CreateSession(username string) *models.Session
}

type Authentication struct {
	Hasher
	Session
}

type Repository struct {
	Authentication *Authentication
	Database
}

func NewRepository(postgres *Postgres) (*Repository, error) {

	hasher := NewHasher()
	session, err := NewSessionManager()
	if err != nil {
		return nil, err
	}

	return &Repository{
		Authentication: &Authentication{
			Hasher:  hasher,
			Session: session,
		},
		Database: postgres,
	}, nil
}
