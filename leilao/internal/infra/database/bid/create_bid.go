package bid

import (
	"context"
	"sync"

	"github.com/thalisonh/auction/configuration/logger"
	"github.com/thalisonh/auction/internal/entity/auction_entity"
	"github.com/thalisonh/auction/internal/entity/bid_entity"
	"github.com/thalisonh/auction/internal/infra/database/auction"
	"github.com/thalisonh/auction/internal/internal_error"
	"go.mongodb.org/mongo-driver/mongo"
)

type BidEntityMongo struct {
	Id        string  `bson:"_id"`
	UserId    string  `bson:"user_id"`
	AuctionId string  `bson:"auction_id"`
	Amount    float64 `bson:"amount"`
	Timestamp int64   `bson:"timestamp"`
}

type BidRepository struct {
	Collection        *mongo.Collection
	AuctionRepository *auction.AuctionRepository
}

func (bd *BidRepository) CreateBid(
	ctx context.Context,
	bidEntities []bid_entity.Bid,
) *internal_error.InternalError {
	var wg sync.WaitGroup

	for _, bid := range bidEntities {
		wg.Add(1)

		go func(bidValue bid_entity.Bid) {
			defer wg.Done()

			auctionEntity, err := bd.AuctionRepository.FindAuctionById(ctx, bidValue.AuctionId)
			if err != nil {
				logger.Error("Error trying to find auction by id", err)
				return
			}

			if auctionEntity.Status != auction_entity.Active {
				return
			}

			bidEntityMongo := &BidEntityMongo{
				Id:        bidValue.Id,
				UserId:    bidValue.UserId,
				AuctionId: bidValue.AuctionId,
				Amount:    bidValue.Amount,
				Timestamp: bidValue.Timestamp.Unix(),
			}

			if _, err := bd.Collection.InsertOne(ctx, bidEntityMongo); err != nil {
				logger.Error("Error trying to insert bid", err)

				return
			}
		}(bid)
	}

	wg.Wait()

	return nil
}
