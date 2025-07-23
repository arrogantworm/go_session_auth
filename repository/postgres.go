package repository

import (
	"context"
	"fmt"
	custom_errors "session-auth/errors"
	"session-auth/models"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

type Postgres struct {
	db     *pgxpool.Pool
	pgOnce sync.Once
}

func NewPostgres(ctx context.Context, cfg *Config) (*Postgres, error) {
	postgres := &Postgres{}
	var err error

	postgres.pgOnce.Do(func() {
		postgres.db, err = pgxpool.New(ctx, fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode))
	})

	if err != nil {
		return nil, err
	}

	if err := postgres.ping(ctx); err != nil {
		return nil, err
	}

	return postgres, nil
}

func (pg *Postgres) ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *Postgres) Close(ctx context.Context) {
	pg.db.Close()
}

// Auth

func (pg *Postgres) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, password) VALUES (@username, @password)`
	args := pgx.NamedArgs{
		"username": user.Username,
		"password": user.PasswordHash,
	}
	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "users_username_key" {
				return &custom_errors.BadRequestError{Message: "username is already taken"}
			} else {
				return err
			}
		}
		return err
	}
	return nil
}

func (pg *Postgres) GetUserByUsername(ctx context.Context, user *models.User) error {
	query := `SELECT id, password FROM users WHERE username=@username`
	args := pgx.NamedArgs{
		"username": user.Username,
	}

	row := pg.db.QueryRow(ctx, query, args)
	if err := row.Scan(&user.Id, &user.PasswordHash); err != nil {
		return err
	}
	return nil

}

func (pg *Postgres) SaveSession(ctx context.Context, session *models.Session) error {
	query := `INSERT INTO sessions (session_id, user_id, expires_at) VALUES (@sessionId, @userId, @expiresAt)`
	args := pgx.NamedArgs{
		"sessionId": session.Token,
		"userId":    session.UserId,
		"expiresAt": session.ExpiresAt,
	}

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

func (pg *Postgres) GetSessionByID(ctx context.Context, session_id uuid.UUID) (*models.Session, error) {
	query := `SELECT session_id, user_id, expires_at, users.username 
		FROM sessions 
		JOIN users ON sessions.user_id=users.id
		WHERE sessions.session_id=@sessionId`
	args := pgx.NamedArgs{
		"sessionId": session_id,
	}

	row := pg.db.QueryRow(ctx, query, args)

	var session models.Session
	if err := row.Scan(&session.Token, &session.UserId, &session.ExpiresAt, &session.Username); err != nil {
		return nil, err
	}
	return &session, nil
}

func (pg *Postgres) DeleteSessionByID(ctx context.Context, session_id uuid.UUID) error {
	query := `DELETE FROM sessions WHERE session_id=@sessionId`
	args := pgx.NamedArgs{
		"sessionId": session_id,
	}

	_, err := pg.db.Exec(ctx, query, args)
	return err
}
