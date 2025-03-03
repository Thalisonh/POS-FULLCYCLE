package user_usecase

import (
	"context"

	"github.com/thalisonh/auction/internal/entity/user_entity"
	"github.com/thalisonh/auction/internal/internal_error"
)

type UserUseCase struct {
	UserRepository user_entity.UserRepositoryInterface
}

type UserOutputDTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserUseCaseInterface interface {
	FindUserBydId(ctx context.Context, id string) (*UserOutputDTO, *internal_error.InternalError)
}

func (u *UserUseCase) FindUserBydId(
	ctx context.Context, id string,
) (*UserOutputDTO, *internal_error.InternalError) {
	userEntiry, err := u.UserRepository.FindUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &UserOutputDTO{
		Id:   userEntiry.Id,
		Name: userEntiry.Name,
	}, nil
}
