package user

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, user User) (User, error)
	GetByID(ctx context.Context, userId int64) (User, error)
}
