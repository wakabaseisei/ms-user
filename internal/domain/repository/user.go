package repository

import (
	"context"

	"github.com/wakabaseisei/ms-user/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, cmd *domain.UserCommand) error
	FindByID(ctx context.Context, ID string) (*domain.User, error)
	Ping(ctx context.Context) error
}
