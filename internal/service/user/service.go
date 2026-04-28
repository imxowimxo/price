package user

import (
	us "Price/internal/domain/user"
	"context"
	"errors"
)

type Service interface {
	Create(ctx context.Context, user us.User) (us.User, error)
	GetByID(ctx context.Context, userID int64) (us.User, error)
}

type service struct {
	repo us.Repository
}

func NewService(repo us.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, user us.User) (us.User, error) {
	if user.Username == "" {
		return us.User{}, errors.New("username не может быть пустым")
	}

	if user.TgID == 0 {
		return us.User{}, errors.New("пустое тг айди пользователя")
	}

	res, err := s.repo.Create(ctx, user)
	if err != nil {
		return us.User{}, err
	}
	return res, nil
}

func (s *service) GetByID(ctx context.Context, userID int64) (us.User, error) {
	return s.repo.GetByID(ctx, userID)
}
