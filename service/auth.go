package service

import (
	"context"
	"errors"
	"fmt"
	custom_errors "session-auth/errors"
	"session-auth/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Service) RegisterUser(ctx context.Context, credentials *models.Credentials) error {

	if credentials.Username == "" || credentials.Password == "" {
		return &custom_errors.BadRequestError{Message: "all fields must be filled"}
	}

	var DBUser models.User
	DBUser.Username = credentials.Username
	hashedPassword, err := s.HashPassword(credentials.Password)
	DBUser.PasswordHash = hashedPassword
	if err != nil {
		return &custom_errors.InternalError{Message: fmt.Sprintf("error hashing password: %v", err)}
	}

	if err := s.CreateUser(ctx, &DBUser); err != nil {
		var badRequestError *custom_errors.BadRequestError
		if errors.As(err, &badRequestError) {
			return badRequestError
		} else {
			return &custom_errors.InternalError{Message: err.Error()}
		}
	}
	return nil
}

func (s *Service) Authenticate(ctx context.Context, credentials *models.Credentials) (*models.Session, error) {
	if credentials.Username == "" || credentials.Password == "" {
		return nil, &custom_errors.BadRequestError{Message: "all fields must be filled"}
	}

	var user models.User
	user.Username = credentials.Username
	if err := s.Database.GetUserByUsername(ctx, &user); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &custom_errors.BadRequestError{Message: "user is not registered"}
		}
		return nil, &custom_errors.InternalError{Message: err.Error()}
	}

	if !s.CheckPassword(user.PasswordHash, credentials.Password) {
		return nil, &custom_errors.BadRequestError{Message: "wrong password"}
	}

	session := s.CreateSession(user.Username)
	session.UserId = user.Id

	if err := s.SaveSession(ctx, session); err != nil {
		return nil, &custom_errors.InternalError{Message: err.Error()}
	}

	return session, nil
}

func (s *Service) ValidateSession(ctx context.Context, sessionToken uuid.UUID) (*models.Session, error) {

	session, err := s.GetSessionByID(ctx, sessionToken)
	if err != nil {
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, &custom_errors.UnauthorizedError{Message: "session expired"}
	}
	return session, nil

}

func (s *Service) Logout(ctx context.Context, sessionToken string) error {

	sessionId, err := uuid.Parse(sessionToken)
	if err != nil {
		return err
	}

	return s.DeleteSessionByID(ctx, sessionId)
}
