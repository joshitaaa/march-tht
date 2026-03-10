package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"march-tht/internal/loader"
)

type fakeClock struct {
	now time.Time
}

func (f fakeClock) Now() time.Time { return f.now }

type fakeRepo struct {
	created map[string]int64
}

func (f *fakeRepo) Init(_ context.Context) error { return nil }

func (f *fakeRepo) IsTeamRedeemed(_ context.Context, teamName string) (bool, error) {
	_, ok := f.created[teamName]
	return ok, nil
}

func (f *fakeRepo) CreateRedemption(_ context.Context, teamName string, redeemedAt int64) (bool, error) {
	if _, ok := f.created[teamName]; ok {
		return false, nil
	}
	f.created[teamName] = redeemedAt
	return true, nil
}

func TestRedeem_Success(t *testing.T) {
	t.Parallel()
	mapping := createMappingStore(t, "staff_pass_id,team_name,created_at\nA001,Team Red,1700000000000\n")
	repo := &fakeRepo{created: map[string]int64{}}

	now := time.UnixMilli(1700100000000)
	svc := NewRedemptionService(mapping, repo).WithClock(fakeClock{now: now})

	result, err := svc.Redeem(context.Background(), "A001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Status != StatusRedeemed {
		t.Fatalf("expected redeemed status, got %s", result.Status)
	}
	if result.TeamName != "Team Red" {
		t.Fatalf("unexpected team: %s", result.TeamName)
	}
	if result.RedeemedAt != now.UnixMilli() {
		t.Fatalf("unexpected redeemed timestamp: %d", result.RedeemedAt)
	}
}

func TestRedeem_AlreadyRedeemed(t *testing.T) {
	t.Parallel()
	mapping := createMappingStore(t, "staff_pass_id,team_name,created_at\nA001,Team Red,1700000000000\n")
	repo := &fakeRepo{created: map[string]int64{"Team Red": 1700000000000}}
	svc := NewRedemptionService(mapping, repo)

	result, err := svc.Redeem(context.Background(), "A001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Status != StatusAlreadyRedeemed {
		t.Fatalf("expected already_redeemed status, got %s", result.Status)
	}
}

func TestRedeem_UnknownStaffPass(t *testing.T) {
	t.Parallel()
	mapping := createMappingStore(t, "staff_pass_id,team_name,created_at\nA001,Team Red,1700000000000\n")
	repo := &fakeRepo{created: map[string]int64{}}
	svc := NewRedemptionService(mapping, repo)

	result, err := svc.Redeem(context.Background(), "A999")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Status != StatusStaffNotFound {
		t.Fatalf("expected staff_not_found status, got %s", result.Status)
	}
}

func createMappingStore(t *testing.T, content string) *loader.MappingStore {
	t.Helper()
	path := writeTempCSV(t, content)
	store, err := loader.NewMappingStoreFromCSV(path)
	if err != nil {
		t.Fatalf("failed to create mapping store: %v", err)
	}
	return store
}

func writeTempCSV(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "mapping.csv")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp csv: %v", err)
	}
	return path
}
