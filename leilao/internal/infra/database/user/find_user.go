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

type UserEntityMongo struct {
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

	var userEntityMongo UserEntityMongo
	err := ur.Collection.FindOne(ctx, filter).Decode(&userEntityMongo)
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
		Id:   userEntityMongo.Id,
		Name: userEntityMongo.Name,
	}, nil
}

func (ar *UserRespository) FindUsers(
	ctx context.Context,
) ([]user_entity.User, *internal_error.InternalError) {
	filter := bson.M{}

	cursor, err := ar.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error trying to find users", err)
		return nil, internal_error.NewInternalServerError("Error trying to find users")
	}
	defer cursor.Close(ctx)

	var entityMongo []UserEntityMongo
	if err := cursor.All(ctx, &entityMongo); err != nil {
		logger.Error("Error trying to find users", err)
		return nil, internal_error.NewInternalServerError("Error trying to find users")
	}

	var entity []user_entity.User
	for _, item := range entityMongo {
		entity = append(entity, user_entity.User{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	return entity, nil
}
