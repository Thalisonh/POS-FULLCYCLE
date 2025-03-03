package bid_usecase

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/thalisonh/auction/configuration/logger"
	"github.com/thalisonh/auction/internal/entity/bid_entity"
	"github.com/thalisonh/auction/internal/internal_error"
)

type BidInputtDTO struct {
	UserId    string  `json:"user_id"`
	AuctionId string  `json:"auction_id"`
	Amount    float64 `json:"amount"`
}

type BidOutputDTO struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	AuctionId string    `json:"auction_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type BidUsecase struct {
	BidRepository bid_entity.BidEntityRepository

	timer               *time.Timer
	maxBatchSize        int
	batchInsertInterval time.Duration
	bidChannel          chan bid_entity.Bid
}

var bidBatch []bid_entity.Bid

type BidUsecaseInterface interface {
	CreateBid(
		ctx context.Context,
		bidInputDTO BidInputtDTO,
	) *internal_error.InternalError
	FindWinningBidByAuctionId(
		ctx context.Context,
		auctionId string,
	) (*BidOutputDTO, *internal_error.InternalError)
	FindBidByAuctionId(
		ctx context.Context,
		auctionId string,
	) ([]BidOutputDTO, *internal_error.InternalError)
}

func NewBidUsecase(bidRepository bid_entity.BidEntityRepository) BidUsecaseInterface {
	maxBatchSizeInterval := getMaxBatchSizeInterval()
	maxBatchSize := getMaxBatchSize()

	bidUseCase := &BidUsecase{
		BidRepository:       bidRepository,
		maxBatchSize:        maxBatchSize,
		batchInsertInterval: maxBatchSizeInterval,
		timer:               time.NewTimer(maxBatchSizeInterval),
		bidChannel:          make(chan bid_entity.Bid, maxBatchSize),
	}

	bidUseCase.triggerCreateRoutine(context.Background())

	return bidUseCase
}

func (bu *BidUsecase) triggerCreateRoutine(ctx context.Context) {
	go func() {
		defer close(bu.bidChannel)

		for {
			select {
			case bidEntity, ok := <-bu.bidChannel:
				if !ok {
					if len(bu.bidChannel) > 0 {
						if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
							logger.Error("error trying to process bid batch list", err)
						}
					}
					return
				}
				bidBatch = append(bidBatch, bidEntity)

				if len(bidBatch) >= bu.maxBatchSize {
					if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
						logger.Error("error trying to process bid batch list", err)
					}

					bidBatch = nil
					bu.timer.Reset(bu.batchInsertInterval)
				}
			case <-bu.timer.C:
				if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
					logger.Error("error trying to process bid batch list", err)
				}

				bidBatch = nil
				bu.timer.Reset(bu.batchInsertInterval)
			}

		}

	}()
}

func (bu *BidUsecase) CreateBid(
	ctx context.Context,
	bidInputDTO BidInputtDTO,
) *internal_error.InternalError {

	bidEntity, err := bid_entity.CreateBid(bidInputDTO.UserId, bidInputDTO.AuctionId, bidInputDTO.Amount)
	if err != nil {
		return err
	}

	bu.bidChannel <- *bidEntity

	return nil
}

func getMaxBatchSizeInterval() time.Duration {
	batchInsertInterval := os.Getenv("BATCH_INSERT_INTERVAL")
	duration, err := time.ParseDuration(batchInsertInterval)
	if err != nil {
		return 3 * time.Minute
	}

	return duration
}

func getMaxBatchSize() int {
	value, err := strconv.Atoi(os.Getenv("MAX_BATCH_SIZE"))
	if err != nil {
		return 5
	}

	return value
}
