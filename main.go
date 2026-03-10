package main

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"march-tht/internal/loader"
	"march-tht/internal/repository/sqlite"
	"march-tht/internal/service"
)

func main() {
	mappingCSVPath := flag.String("mapping-csv", "./mapping.sample.csv", "path to mapping CSV file")
	dbPath := flag.String("db-path", "./redemptions.db", "path to sqlite database")
	flag.Parse()

	mappingStore, err := loader.NewMappingStoreFromCSV(*mappingCSVPath)
	if err != nil {
		log.Fatalf("failed to load mapping file: %v", err)
	}

	svc, closeFn, err := buildService(mappingStore, *dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer closeFn()

	runMenu(svc)
}

func buildService(mappingStore *loader.MappingStore, dbPath string) (*service.RedemptionService, func(), error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open sqlite db: %w", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	repo := sqlite.NewRedemptionRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := repo.Init(ctx); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to initialize sqlite schema: %w", err)
	}

	svc := service.NewRedemptionService(mappingStore, repo)
	return svc, func() { _ = db.Close() }, nil
}

func runMenu(svc *service.RedemptionService) {
	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nChristmas Gift Redemption")
		fmt.Println("1) Lookup team")
		fmt.Println("2) Redeem gift")
		fmt.Println("3) Exit")

		choice := prompt(in, "Choose [1-3]: ")
		switch choice {
		case "1":
			staffPassID := prompt(in, "Enter staff_pass_id: ")
			team, found, err := svc.LookupTeam(staffPassID)
			if err != nil {
				printResultBlock("ERROR", []string{err.Error()})
				pause(in)
				continue
			}
			if !found {
				printResultBlock("LOOKUP RESULT", []string{"Status: staff_not_found"})
				pause(in)
				continue
			}
			printResultBlock("LOOKUP RESULT", []string{"Status: found", fmt.Sprintf("Team: %s", team)})
			pause(in)
		case "2":
			staffPassID := prompt(in, "Enter staff_pass_id: ")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			result, err := svc.Redeem(ctx, staffPassID)
			cancel()
			if err != nil {
				printResultBlock("ERROR", []string{err.Error()})
				pause(in)
				continue
			}
			lines := []string{fmt.Sprintf("Status: %s", result.Status)}
			if result.TeamName != "" {
				lines = append(lines, fmt.Sprintf("Team: %s", result.TeamName))
			}
			if result.RedeemedAt != 0 {
				lines = append(lines, fmt.Sprintf("RedeemedAt: %d", result.RedeemedAt))
			}
			printResultBlock("REDEEM RESULT", lines)
			pause(in)
		case "3":
			fmt.Println("Bye")
			return
		default:
			fmt.Println("Invalid choice, try again")
			pause(in)
		}
	}
}

func prompt(in *bufio.Reader, text string) string {
	fmt.Print(text)
	line, _ := in.ReadString('\n')
	return strings.TrimSpace(line)
}

func pause(in *bufio.Reader) {
	fmt.Print("Press Enter to continue...")
	_, _ = in.ReadString('\n')
}

func printResultBlock(title string, lines []string) {
	fmt.Println("\n================================")
	fmt.Println(title)
	fmt.Println("--------------------------------")
	for _, line := range lines {
		fmt.Println(line)
	}
	fmt.Println("================================")
}
