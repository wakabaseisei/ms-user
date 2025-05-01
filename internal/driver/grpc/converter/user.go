package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	userv1 "buf.build/gen/go/wakabaseisei/ms-protobuf/protocolbuffers/go/ms/user/v1"
	"github.com/wakabaseisei/ms-user/internal/domain"
)

func ConvertUserToUserPb(user *domain.User) *userv1.User {
	return &userv1.User{
		UserId:    user.ID,
		Name:      user.Name,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
