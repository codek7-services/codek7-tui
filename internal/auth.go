package internal

import (
	"context"

	proto "github.com/codek7-services/codek7-tui/pkg/pb"
)

type AuthService struct {
	client proto.RepoServiceClient
}

func (a *AuthService) Register(username string) error {
	_, err := a.client.CreateUser(context.TODO(), &proto.CreateUserRequest{
		Username: username,
	})
	return err
}

func (a *AuthService) GetUser(id string) (*proto.UserResponse, error) {
	return a.client.GetUser(context.TODO(), &proto.GetUserRequest{
		Username: id,
	})
}
