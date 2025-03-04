package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/thalisonh/auction/configuration/logger"
	"github.com/thalisonh/auction/internal/entity/user_entity"
	"github.com/thalisonh/auction/internal/internal_error"
)

func (ar *UserRespository) CreateUser(
	ctx context.Context,
	userEntity *user_entity.User,
) *internal_error.InternalError {
	userEntityMongo := &UserEntityMongo{
		Id:   uuid.NewString(),
		Name: userEntity.Name,
	}

	_, err := ar.Collection.InsertOne(ctx, userEntityMongo)
	if err != nil {
		logger.Error("Error trying to inser auction", err)
		return internal_error.NewInternalServerError("Error trying to inser auction")
	}

	return nil
}
