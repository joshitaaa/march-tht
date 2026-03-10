package loader

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	headerStaffPassID = "staff_pass_id"
	headerTeamName    = "team_name"
	headerCreatedAt   = "created_at"
)

type MappingStore struct {
	byStaffPass map[string]string
}

// NewMappingStoreFromCSV loads the CSV mapping file from the path indicated.
func NewMappingStoreFromCSV(path string) (*MappingStore, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open mapping csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, errors.New("mapping csv is empty")
		}
		return nil, fmt.Errorf("read csv header: %w", err)
	}

	index, err := headerIndex(headers)
	if err != nil {
		return nil, err
	}

	store := &MappingStore{byStaffPass: make(map[string]string)}
	rowNum := 1
	for {
		rowNum++
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read row %d: %w", rowNum, err)
		}

		staffPassID := strings.TrimSpace(record[index[headerStaffPassID]])
		teamName := strings.TrimSpace(record[index[headerTeamName]])
		createdAtRaw := strings.TrimSpace(record[index[headerCreatedAt]])

		if staffPassID == "" || teamName == "" || createdAtRaw == "" {
			return nil, fmt.Errorf("row %d contains empty required field", rowNum)
		}

		if _, err := strconv.ParseInt(createdAtRaw, 10, 64); err != nil {
			return nil, fmt.Errorf("row %d has invalid created_at: %w", rowNum, err)
		}

		store.byStaffPass[staffPassID] = teamName
	}

	return store, nil
}

func headerIndex(headers []string) (map[string]int, error) {
	index := map[string]int{}
	for i, header := range headers {
		normalized := strings.TrimSpace(header)
		index[normalized] = i
	}

	for _, required := range []string{headerStaffPassID, headerTeamName, headerCreatedAt} {
		if _, ok := index[required]; !ok {
			return nil, fmt.Errorf("missing required header: %s", required)
		}
	}
	return index, nil
}

// Returns the team name for a staff pass ID
func (m *MappingStore) TeamByStaffPass(staffPassID string) (string, bool) {
	team, ok := m.byStaffPass[staffPassID]
	return team, ok
}
