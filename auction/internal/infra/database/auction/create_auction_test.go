package auction_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/thalisonh/auction/configuration/database/mongodb"
	"github.com/thalisonh/auction/internal/entity/auction_entity"
	"github.com/thalisonh/auction/internal/infra/database/auction"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()

	t.Setenv("MONGODB_URL", "mongodb://localhost:27017")
	t.Setenv("MONGODB_DB", "auction")

	database, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())

		return
	}
	repository := auction.NewAuctionRepository(database)

	t.Run("Create", func(t *testing.T) {
		// force the auction interval to be 1 second
		t.Setenv("AUCTION_INTERVAL", "1s")

		auctionId := uuid.NewString()
		ctx := context.Background()

		err := repository.CreateAuction(ctx, &auction_entity.Auction{
			Id: auctionId,
		})

		assert.Nil(t, err)

		// wait for interval to complete
		time.Sleep(time.Second * 2)
		auction, err := repository.FindAuctionById(ctx, auctionId)

		assert.Nil(t, err)
		assert.Equal(t, auctionId, auction.Id)
		assert.Equal(t, auction_entity.Completed, auction.Status)
	})
}
