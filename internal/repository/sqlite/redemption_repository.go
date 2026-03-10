package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// RedemptionRepository stores redemption data in SQLite
type RedemptionRepository struct {
	db *sql.DB
}

func NewRedemptionRepository(db *sql.DB) *RedemptionRepository {
	return &RedemptionRepository{db: db}
}

func (r *RedemptionRepository) Init(ctx context.Context) error {
	const schema = `
CREATE TABLE IF NOT EXISTS redemptions (
	team_name TEXT PRIMARY KEY,
	redeemed_at INTEGER NOT NULL
);`
	if _, err := r.db.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("create redemptions table: %w", err)
	}
	return nil
}

func (r *RedemptionRepository) CreateRedemption(ctx context.Context, teamName string, redeemedAt int64) (bool, error) {
	const stmt = `
INSERT INTO redemptions (team_name, redeemed_at)
VALUES (?, ?)
ON CONFLICT(team_name) DO NOTHING;`
	result, err := r.db.ExecContext(ctx, stmt, teamName, redeemedAt)
	if err != nil {
		return false, fmt.Errorf("insert redemption: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("redemption rows affected: %w", err)
	}
	return affected == 1, nil
}
