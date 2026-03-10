package sqlite

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
)

func TestRedemptionRepository_CreateRedemption(t *testing.T) {
	t.Parallel()
	db := openTempDB(t)
	repo := NewRedemptionRepository(db)

	if err := repo.Init(context.Background()); err != nil {
		t.Fatalf("init repo: %v", err)
	}

	created, err := repo.CreateRedemption(context.Background(), "Team Red", 1700000000000)
	if err != nil {
		t.Fatalf("create redemption: %v", err)
	}
	if !created {
		t.Fatalf("expected first insert to succeed")
	}

	created, err = repo.CreateRedemption(context.Background(), "Team Red", 1700000000100)
	if err != nil {
		t.Fatalf("create duplicate redemption: %v", err)
	}
	if created {
		t.Fatalf("expected duplicate insert to be ignored")
	}
}

func openTempDB(t *testing.T) *sql.DB {
	t.Helper()
	path := filepath.Join(t.TempDir(), "redemptions.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	t.Cleanup(func() {
		_ = db.Close()
	})
	return db
}
