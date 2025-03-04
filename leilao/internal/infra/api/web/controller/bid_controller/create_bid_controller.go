package bid_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thalisonh/auction/configuration/rest_err"
	"github.com/thalisonh/auction/internal/infra/api/web/validation"
	"github.com/thalisonh/auction/internal/usecase/bid_usecase"
)

type BidController struct {
	bidUseCase bid_usecase.BidUsecaseInterface
}

func NewBidController(bidUseCase bid_usecase.BidUsecaseInterface) *BidController {
	return &BidController{
		bidUseCase: bidUseCase,
	}
}

func (u *BidController) CreateBid(c *gin.Context) {
	var bidInput bid_usecase.BidInputDTO
	if err := c.ShouldBindJSON(&bidInput); err != nil {
		restErr := validation.ValidateErr(err)

		c.JSON(restErr.Code, restErr)

		return
	}

	err := u.bidUseCase.CreateBid(c, bidInput)
	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)

		return
	}

	c.Status(http.StatusCreated)
}
