package user_usecase

import (
	"context"

	"github.com/thalisonh/auction/internal/entity/user_entity"
	"github.com/thalisonh/auction/internal/internal_error"
)

func (au *UserUseCase) CreateUser(
	ctx context.Context,
	userInput UserInputDTO,
) *internal_error.InternalError {
	if err := au.UserRepository.CreateUser(ctx, &user_entity.User{
		Id:   userInput.Id,
		Name: userInput.Name,
	}); err != nil {
		return err
	}

	return nil
}
