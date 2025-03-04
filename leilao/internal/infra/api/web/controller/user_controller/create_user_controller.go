package user_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thalisonh/auction/configuration/rest_err"
	"github.com/thalisonh/auction/internal/infra/api/web/validation"
	"github.com/thalisonh/auction/internal/usecase/user_usecase"
)

func (u *UserController) CreateUser(c *gin.Context) {
	var input user_usecase.UserInputDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		restErr := validation.ValidateErr(err)

		c.JSON(restErr.Code, restErr)

		return
	}

	err := u.userUseCase.CreateUser(c, input)
	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)

		return
	}

	c.Status(http.StatusCreated)
}
