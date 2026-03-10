package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"march-tht/internal/loader"
	"march-tht/internal/repository"
)

const (
	StatusRedeemed        = "redeemed"
	StatusAlreadyRedeemed = "already_redeemed"
	StatusStaffNotFound   = "staff_not_found"
)

var ErrInvalidStaffPassID = errors.New("staff_pass_id is required")

type Clock interface {
	Now() time.Time
}

type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }

type RedeemResult struct {
	Status     string `json:"status"`
	TeamName   string `json:"team_name,omitempty"`
	RedeemedAt int64  `json:"redeemed_at,omitempty"`
}

type RedemptionService struct {
	mapping *loader.MappingStore
	repo    repository.RedemptionRepository
	clock   Clock
}

func NewRedemptionService(mapping *loader.MappingStore, repo repository.RedemptionRepository) *RedemptionService {
	return &RedemptionService{mapping: mapping, repo: repo, clock: systemClock{}}
}

func (s *RedemptionService) WithClock(clock Clock) *RedemptionService {
	s.clock = clock
	return s
}

func (s *RedemptionService) LookupTeam(staffPassID string) (string, bool, error) {
	id := strings.TrimSpace(staffPassID)
	if id == "" {
		return "", false, ErrInvalidStaffPassID
	}
	team, ok := s.mapping.TeamByStaffPass(id)
	return team, ok, nil
}

func (s *RedemptionService) Redeem(ctx context.Context, staffPassID string) (RedeemResult, error) {
	team, ok, err := s.LookupTeam(staffPassID)
	if err != nil {
		return RedeemResult{}, err
	}
	if !ok {
		return RedeemResult{Status: StatusStaffNotFound}, nil
	}

	nowMs := s.clock.Now().UnixMilli()
	created, err := s.repo.CreateRedemption(ctx, team, nowMs)
	if err != nil {
		return RedeemResult{}, err
	}
	if !created {
		return RedeemResult{Status: StatusAlreadyRedeemed, TeamName: team}, nil
	}

	return RedeemResult{Status: StatusRedeemed, TeamName: team, RedeemedAt: nowMs}, nil
}
