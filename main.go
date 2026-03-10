package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"march-tht/internal/loader"
)

func main() {
	mappingCSVPath := flag.String("mapping-csv", "./mapping.sample.csv", "path to mapping CSV file")
	flag.Parse()

	mappingStore, err := loader.NewMappingStoreFromCSV(*mappingCSVPath)
	if err != nil {
		log.Fatalf("failed to load mapping file: %v", err)
	}

	runMenu(mappingStore)
}

func runMenu(store *loader.MappingStore) {
	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nCSV Staff Lookup")
		fmt.Println("1) Lookup team")
		fmt.Println("2) Exit")

		choice := prompt(in, "Choose [1-2]: ")
		switch choice {
		case "1":
			staffPassID := prompt(in, "Enter staff_pass_id: ")
			team, found := store.TeamByStaffPass(strings.TrimSpace(staffPassID))
			if !found {
				printResultBlock("LOOKUP RESULT", []string{"Status: staff_not_found"})
				pause(in)
				continue
			}
			printResultBlock("LOOKUP RESULT", []string{"Status: found", fmt.Sprintf("Team: %s", team)})
			pause(in)
		case "2":
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
