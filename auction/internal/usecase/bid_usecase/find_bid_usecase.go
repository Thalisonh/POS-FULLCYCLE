package bid_usecase

import (
	"context"

	"github.com/thalisonh/auction/internal/internal_error"
)

func (bu *BidUsecase) FindBidByAuctionId(
	ctx context.Context,
	auctionId string,
) ([]BidOutputDTO, *internal_error.InternalError) {
	bidList, err := bu.BidRepository.FindBidByAuctionId(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	outputList := []BidOutputDTO{}
	for _, item := range bidList {
		outputList = append(outputList, BidOutputDTO{
			Id:        item.Id,
			UserId:    item.UserId,
			AuctionId: item.AuctionId,
			Amount:    item.Amount,
			Timestamp: item.Timestamp,
		})
	}

	return outputList, nil
}

func (bu *BidUsecase) FindWinningBidByAuctionId(
	ctx context.Context,
	auctionId string,
) (*BidOutputDTO, *internal_error.InternalError) {
	bidEntity, err := bu.BidRepository.FindWinningBidByAuctionId(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	return &BidOutputDTO{
		Id:        bidEntity.Id,
		UserId:    bidEntity.UserId,
		AuctionId: bidEntity.AuctionId,
		Amount:    bidEntity.Amount,
		Timestamp: bidEntity.Timestamp,
	}, nil
}
