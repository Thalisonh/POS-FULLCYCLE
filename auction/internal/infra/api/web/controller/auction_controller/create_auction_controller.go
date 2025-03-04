package auction_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thalisonh/auction/configuration/rest_err"
	"github.com/thalisonh/auction/internal/infra/api/web/validation"
	"github.com/thalisonh/auction/internal/usecase/auction_usecase"
)

type AuctionController struct {
	auctionUseCase auction_usecase.AuctionUseCaseInterface
}

func NewAuctionController(auctionUseCase auction_usecase.AuctionUseCaseInterface) *AuctionController {
	return &AuctionController{
		auctionUseCase: auctionUseCase,
	}
}

func (u *AuctionController) CreateAuction(c *gin.Context) {
	var auctionInput auction_usecase.AuctionInputDTO
	if err := c.ShouldBindJSON(&auctionInput); err != nil {
		restErr := validation.ValidateErr(err)

		c.JSON(restErr.Code, restErr)

		return
	}

	err := u.auctionUseCase.CreateAuction(c, auctionInput)
	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)

		return
	}

	c.Status(http.StatusCreated)
}
