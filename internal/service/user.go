package service

import (
	"context"
	"debez/internal/models"
)

type UserRepository interface {
	Select(ctx context.Context, offset, limit int) ([]models.User, error)
	SelectByID(ctx context.Context, id int64) (models.User, error)
	Insert(ctx context.Context, user models.User) error
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id int64) error
}

type UserService struct {
	Repository UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		Repository: repo,
	}
}
func (s *UserService) GetUsers(ctx context.Context, offset, limit int) ([]models.User, error) {
	return s.Repository.Select(ctx, offset, limit)
}
func (s *UserService) SaveUser(ctx context.Context, user models.User) error {
	return s.Repository.Insert(ctx, user)
}
func (s *UserService) GetUserByID(ctx context.Context, id int64) (models.User, error) {
	return s.Repository.SelectByID(ctx, id)
}
func (s *UserService) UpdateUser(ctx context.Context, user models.User) error {
	return s.Repository.Update(ctx, user)
}
func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	return s.Repository.Delete(ctx, id)
}
