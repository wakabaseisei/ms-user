package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/wakabaseisei/ms-user/internal/domain"
)

type userRepository struct {
	conn *sql.DB
}

func NewUserRepository(conn *sql.DB) *userRepository {
	return &userRepository{
		conn: conn,
	}
}

func (r *userRepository) Create(ctx context.Context, cmd *domain.UserCommand) error {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `INSERT INTO Users (UserID, Name, CreatedAt) VALUES (?, ?, ?)`
	_, err = tx.ExecContext(ctx, query, cmd.ID, cmd.Name, cmd.CreatedAt)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("insert user: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *userRepository) FindByID(ctx context.Context, ID string) (*domain.User, error) {
	query := `SELECT UserID, Name, CreatedAt FROM Users WHERE UserID = ?`

	var user domain.User
	err := r.conn.QueryRowContext(ctx, query, ID).Scan(&user.ID, &user.Name, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("fetch user: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Ping(ctx context.Context) error {
	if err := r.conn.Ping(); err != nil {
		log.Printf("Ping: %v", err)
		return err
	}
	return nil
}
