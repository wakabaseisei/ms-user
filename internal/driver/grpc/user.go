package grpc

import (
	"context"
	"fmt"
	"log"

	"connectrpc.com/connect"
	"github.com/wakabaseisei/ms-user/internal/domain"
	"github.com/wakabaseisei/ms-user/internal/driver/grpc/converter"
	"github.com/wakabaseisei/ms-user/internal/usecase"

	userv1 "buf.build/gen/go/wakabaseisei/ms-protobuf/protocolbuffers/go/ms/user/v1"
)

func (s *UserService) CreateUser(
	ctx context.Context,
	req *connect.Request[userv1.CreateUserRequest],
) (*connect.Response[userv1.User], error) {
	log.Println("Request headers: ", req.Header())

	newUser := req.Msg.GetUser()

	cmd := &domain.UserCommand{
		ID:        newUser.GetUserId(),
		Name:      newUser.GetName(),
		CreatedAt: newUser.GetCreatedAt().AsTime(),
	}
	user, uerr := usecase.NewCreateUserInteractor(s.services.UserRepository).Invoke(ctx, cmd)
	if uerr != nil {
		return nil, fmt.Errorf("usecase.CreateUserInteractor.Invoke(): %v", uerr)
	}

	res := connect.NewResponse(converter.ConvertUserToUserPb(user))
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}
