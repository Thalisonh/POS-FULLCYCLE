package auction_usecase

import (
	"context"
	"time"

	"github.com/thalisonh/auction/internal/entity/auction_entity"
	"github.com/thalisonh/auction/internal/entity/bid_entity"
	"github.com/thalisonh/auction/internal/internal_error"
	"github.com/thalisonh/auction/internal/usecase/bid_usecase"
)

func NewAuctionUseCase(
	auctionRepository auction_entity.AuctionRepositoryInterface,
	bidRepository bid_entity.BidEntityRepository,
) AuctionUseCaseInterface {
	return &AuctionUseCase{
		auctionRepositoryInterface: auctionRepository,
		bidRepositoryInterface:     bidRepository,
	}
}

type AuctionUseCaseInterface interface {
	CreateAuction(
		ctx context.Context,
		auctionInput AuctionInputDTO,
	) *internal_error.InternalError
	FindAuctionById(
		ctx context.Context, id string,
	) (*AuctionOutputDTO, *internal_error.InternalError)
	FindAuctions(
		ctx context.Context,
		status auction_entity.AuctionStatus,
		category, productName string,
	) ([]AuctionOutputDTO, *internal_error.InternalError)
	FindWinningBidByAuctionId(
		ctx context.Context, auctionId string,
	) (*WinningInfoOutputDTO, *internal_error.InternalError)
}

type AuctionInputDTO struct {
	ProductName string           `json:"product_name" binding:"required,min=1"`
	Category    string           `json:"category" binding:"required,min=2"`
	Description string           `json:"description" binding:"required,min=10,max=200"`
	Condition   ProductCondition `json:"condition"`
}

type AuctionOutputDTO struct {
	Id          string           `json:"id"`
	ProductName string           `json:"product_name"`
	Category    string           `json:"category"`
	Description string           `json:"description"`
	Condition   ProductCondition `json:"condition"`
	Status      AuctionStatus    `json:"status"`
	Timestamp   time.Time        `json:"timestamp" time_format:"2006-01-02 15:30:00"`
}

type WinningInfoOutputDTO struct {
	Auction AuctionOutputDTO          `json:"auction"`
	Bid     *bid_usecase.BidOutputDTO `json:"bid,omitempty"`
}

type ProductCondition int64
type AuctionStatus int64

type AuctionUseCase struct {
	auctionRepositoryInterface auction_entity.AuctionRepositoryInterface
	bidRepositoryInterface     bid_entity.BidEntityRepository
}

func (au *AuctionUseCase) CreateAuction(
	ctx context.Context,
	auctionInput AuctionInputDTO,
) *internal_error.InternalError {
	auction, err := auction_entity.CreateAuction(
		auctionInput.ProductName,
		auctionInput.Category,
		auctionInput.Description,
		auction_entity.ProductCondition(auctionInput.Condition))
	if err != nil {
		return err
	}

	if err := au.auctionRepositoryInterface.CreateAuction(ctx, auction); err != nil {
		return err
	}

	return nil
}
