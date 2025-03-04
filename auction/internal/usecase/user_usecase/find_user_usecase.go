package user_usecase

import (
	"context"

	"github.com/thalisonh/auction/internal/entity/user_entity"
	"github.com/thalisonh/auction/internal/internal_error"
)

func NewUserUseCase(userRepository user_entity.UserRepositoryInterface) UserUseCaseInterface {
	return &UserUseCase{
		UserRepository: userRepository,
	}
}

type UserUseCase struct {
	UserRepository user_entity.UserRepositoryInterface
}

type UserOutputDTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserInputDTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserUseCaseInterface interface {
	FindUserById(ctx context.Context, id string) (*UserOutputDTO, *internal_error.InternalError)
	FindUsers(ctx context.Context) ([]UserOutputDTO, *internal_error.InternalError)
	CreateUser(ctx context.Context, userInput UserInputDTO) *internal_error.InternalError
}

func (u *UserUseCase) FindUserById(
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

func (u *UserUseCase) FindUsers(
	ctx context.Context,
) ([]UserOutputDTO, *internal_error.InternalError) {
	userEntity, err := u.UserRepository.FindUsers(ctx)
	if err != nil {
		return nil, err
	}

	outputList := []UserOutputDTO{}
	for _, item := range userEntity {
		outputList = append(outputList, UserOutputDTO{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	return outputList, nil
}
