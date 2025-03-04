package bid_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thalisonh/auction/configuration/rest_err"
)

func (u *BidController) FindBidByAuctionById(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		errRest := rest_err.NewBadRequestError("Invalid auction id", rest_err.Causes{
			Field:   "auction_id",
			Message: "Invalid UUID",
		})

		c.JSON(errRest.Code, errRest)

		return
	}

	BidOutput, err := u.bidUseCase.FindBidByAuctionId(c, auctionId)
	if err != nil {
		errRest := rest_err.ConvertError(err)

		c.JSON(errRest.Code, errRest)

		return
	}

	c.JSON(http.StatusOK, BidOutput)
}
