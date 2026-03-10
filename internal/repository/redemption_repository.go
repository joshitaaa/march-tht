package repository

import "context"

type RedemptionRepository interface {
	Init(ctx context.Context) error
	CreateRedemption(ctx context.Context, teamName string, redeemedAt int64) (bool, error)
}
