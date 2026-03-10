package loader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewMappingStoreFromCSV_Success(t *testing.T) {
	t.Parallel()
	path := writeTempCSV(t, "staff_pass_id,team_name,created_at\nA001,Team Red,1700000000000\n")

	store, err := NewMappingStoreFromCSV(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	team, ok := store.TeamByStaffPass("A001")
	if !ok {
		t.Fatalf("expected staff pass to be found")
	}
	if team != "Team Red" {
		t.Fatalf("unexpected team, got %s", team)
	}
}

func writeTempCSV(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "mapping.csv")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp csv: %v", err)
	}
	return path
}
