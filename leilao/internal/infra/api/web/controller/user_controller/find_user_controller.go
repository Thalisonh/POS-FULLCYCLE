package user_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thalisonh/auction/configuration/rest_err"
	"github.com/thalisonh/auction/internal/usecase/user_usecase"
)

type UserController struct {
	userUseCase user_usecase.UserUseCaseInterface
}

func NewUserController(userUseCase user_usecase.UserUseCaseInterface) *UserController {
	return &UserController{
		userUseCase: userUseCase,
	}
}

func (u *UserController) FindUserById(c *gin.Context) {
	userId := c.Param("userId")

	if err := uuid.Validate(userId); err != nil {
		errRest := rest_err.NewBadRequestError("Invalid user id", rest_err.Causes{
			Field:   "user_id",
			Message: "Invalid UUID",
		})

		c.JSON(errRest.Code, errRest)

		return
	}

	userData, err := u.userUseCase.FindUserById(c, userId)
	if err != nil {
		errRest := rest_err.ConvertError(err)

		c.JSON(errRest.Code, errRest)

		return
	}

	c.JSON(http.StatusOK, userData)
}
