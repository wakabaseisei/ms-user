package repository

type Services struct {
	UserRepository UserRepository
}

func NewServices(userRepository UserRepository) *Services {
	return &Services{
		UserRepository: userRepository,
	}
}
