package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/thalisonh/auction/configuration/logger"
	"github.com/thalisonh/auction/internal/entity/user_entity"
	"github.com/thalisonh/auction/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserEntittMongo struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}

type UserRespository struct {
	Collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) *UserRespository {
	return &UserRespository{
		Collection: database.Collection("users"),
	}
}

func (ur *UserRespository) FindUserById(ctx context.Context, userId string) (*user_entity.User, *internal_error.InternalError) {
	filter := bson.M{"_id": userId}

	var userEntittMongo UserEntittMongo
	err := ur.Collection.FindOne(ctx, filter).Decode(&userEntittMongo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error(fmt.Sprintf("User not found with this id = %s", userId), err)
			return nil, internal_error.NewNotFoundError(
				fmt.Sprintf("User not found with this id = %s", userId))
		}

		logger.Error("Error trying to find user by userId", err)
		return nil, internal_error.NewInternalServerError("Error trying to find user by userId")
	}

	return &user_entity.User{
		Id:   userEntittMongo.Id,
		Name: userEntittMongo.Name,
	}, nil
}
