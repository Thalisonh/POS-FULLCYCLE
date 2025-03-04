package auction_controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thalisonh/auction/configuration/rest_err"
	"github.com/thalisonh/auction/internal/entity/auction_entity"
)

func (u *AuctionController) FindAuctionById(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		errRest := rest_err.NewBadRequestError("Invalid auction id", rest_err.Causes{
			Field:   "auction_id",
			Message: "Invalid UUID",
		})

		c.JSON(errRest.Code, errRest)

		return
	}

	auctionData, err := u.auctionUseCase.FindAuctionById(c, auctionId)
	if err != nil {
		errRest := rest_err.ConvertError(err)

		c.JSON(errRest.Code, errRest)

		return
	}

	c.JSON(http.StatusOK, auctionData)

}

func (u *AuctionController) FindAuctions(c *gin.Context) {
	status := c.Query("status")
	category := c.Query("category")
	productName := c.Query("product_name")

	statusNumber, errConv := strconv.Atoi(status)
	if errConv != nil {
		errRest := rest_err.NewBadRequestError("Invalid status", rest_err.Causes{
			Field:   "status",
			Message: "Invalid status",
		})

		c.JSON(errRest.Code, errRest)

		return
	}

	auctionData, err := u.auctionUseCase.FindAuctions(c, auction_entity.AuctionStatus(statusNumber), category, productName)
	if err != nil {
		errRest := rest_err.ConvertError(err)

		c.JSON(errRest.Code, errRest)

		return
	}

	c.JSON(http.StatusOK, auctionData)
}

func (u *AuctionController) FindWinningBidByAuctionId(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		errRest := rest_err.NewBadRequestError("Invalid auction id", rest_err.Causes{
			Field:   "auction_id",
			Message: "Invalid UUID",
		})

		c.JSON(errRest.Code, errRest)

		return
	}

	auctionData, err := u.auctionUseCase.FindWinningBidByAuctionId(c, auctionId)
	if err != nil {
		errRest := rest_err.ConvertError(err)

		c.JSON(errRest.Code, errRest)

		return
	}

	c.JSON(http.StatusOK, auctionData)
}
