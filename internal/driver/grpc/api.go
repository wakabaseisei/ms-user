package grpc

import "github.com/wakabaseisei/ms-user/internal/domain/repository"

type UserService struct {
	services *repository.Services
}

func NewUserService(services *repository.Services) *UserService {
	return &UserService{
		services: services,
	}
}
