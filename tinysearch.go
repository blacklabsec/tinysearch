package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/blacklabsec/tinysearch/googleapi"
)

func main() {
	// Define command line flags
	urlFlag := flag.String("u", "", "Domain or URL to search")
	listFlag := flag.String("l", "", "File containing list of domains or URLs to search")
	queryFlag := flag.String("q", "", "Additional query to search with the domain")
	fileTypeFlag := flag.String("t", "", "File type to search (e.g., pdf)")
	maxResultsFlag := flag.Int("m", 100, "Maximum number of results to query (default 100)")
	configFile := flag.String("c", "", "Path to configuration file (default: ~/.config/tinysearch/config.json)")
	flag.Parse()

	// Set default config file path if not provided
	if *configFile == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatalf("Failed to get current user: %v", err)
		}
		*configFile = filepath.Join(usr.HomeDir, ".config", "tinysearch", "config.json")
	}

	// Load configuration
	configData, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Configuration file is required at %s with the following format:\n\n%s", *configFile, `{
  "cx": "YOUR_CX_HERE",
  "api_key": "YOUR_API_KEY_HERE"
}`)
	}

	var config googleapi.Config
	if err := json.Unmarshal(configData, &config); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	var queries []string

	// Determine the queries to run
	if *urlFlag != "" {
		baseQuery := fmt.Sprintf("site:%s", *urlFlag)
		if *fileTypeFlag != "" {
			baseQuery = fmt.Sprintf("%s filetype:%s", baseQuery, *fileTypeFlag)
		}
		if *queryFlag != "" {
			queries = append(queries, fmt.Sprintf("%s %s", baseQuery, *queryFlag))
		} else {
			queries = append(queries, baseQuery)
		}
	} else if *listFlag != "" {
		file, err := os.Open(*listFlag)
		if err != nil {
			log.Fatalf("Failed to open list file: %v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				baseQuery := fmt.Sprintf("site:%s", line)
				if *fileTypeFlag != "" {
					baseQuery = fmt.Sprintf("%s filetype:%s", baseQuery, *fileTypeFlag)
				}
				if *queryFlag != "" {
					queries = append(queries, fmt.Sprintf("%s %s", baseQuery, *queryFlag))
				} else {
					queries = append(queries, baseQuery)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading list file: %v", err)
		}
	} else {
		log.Fatal("You must provide either a -u or -l flag")
	}

	// Perform searches
	for _, query := range queries {
		urls, err := googleapi.GoogleSearch(config, query, *maxResultsFlag)
		if err != nil {
			log.Printf("Failed to search for '%s': %v", query, err)
			continue
		}

		// Print results
		for _, u := range urls {
			fmt.Println(u)
		}
	}
}
