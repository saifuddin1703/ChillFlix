package services

import (
	"chillfix/config"
	"chillfix/internal/database"
	"chillfix/internal/repository/user"
	"chillfix/models"
	"context"
)

type UserService struct {
	UserRepository user.UserRepository
}

func NewUserService() (*UserService, error) {
	config, err := config.GetConfig()
	if err != nil {

		return nil, err
	}
	client, err := database.GetMongoClient(config.DatabaseURL, "chillfix")
	if err != nil {
		return nil, err
	}
	return &UserService{
		UserRepository: user.NewMongoUserRepository(client),
	}, nil
}

// add context to all methods
func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	return s.UserRepository.Create(ctx, user)
}
func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.UserRepository.FindByID(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.UserRepository.FindByEmail(ctx, email)
}
func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.UserRepository.Update(ctx, user)
}
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.UserRepository.Delete(ctx, id)
}
