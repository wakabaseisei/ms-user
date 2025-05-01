package usecase

import (
	"context"
	"fmt"

	"github.com/wakabaseisei/ms-user/internal/domain"
	"github.com/wakabaseisei/ms-user/internal/domain/repository"
)

type CreateUserInteractor interface {
	Invoke(ctx context.Context, cmd *domain.UserCommand) (*domain.User, error)
}

type createUserInteractor struct {
	userRepository repository.UserRepository
}

func NewCreateUserInteractor(userRepo repository.UserRepository) CreateUserInteractor {
	return &createUserInteractor{
		userRepository: userRepo,
	}
}

func (i *createUserInteractor) Invoke(
	ctx context.Context,
	cmd *domain.UserCommand,
) (
	*domain.User,
	error,
) {
	if rerr := i.userRepository.Create(ctx, cmd); rerr != nil {
		return nil, fmt.Errorf("userRepository Create: %v", rerr)
	}

	user, ferr := i.userRepository.FindByID(ctx, cmd.ID)
	if ferr != nil {
		return nil, fmt.Errorf("userRepository FindByID: %v", ferr)
	}

	return user, nil
}
