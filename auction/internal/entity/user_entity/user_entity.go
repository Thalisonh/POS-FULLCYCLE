package user_entity

import (
	"context"

	"github.com/thalisonh/auction/internal/internal_error"
)

type User struct {
	Id   string
	Name string
}

type UserRepositoryInterface interface {
	FindUserById(ctx context.Context, userId string) (*User, *internal_error.InternalError)
	FindUsers(ctx context.Context) ([]User, *internal_error.InternalError)
	CreateUser(ctx context.Context, userEntity *User) *internal_error.InternalError
}
