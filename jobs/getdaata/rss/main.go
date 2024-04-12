package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/timotewb/cpu/jobs/data/rss/app"
)

// var getURL string = "https://rss.nytimes.com/services/xml/rss/nyt/World.xml"
// var getURL string = "https://www.rnz.co.nz/rss/business.xml"
var deBug bool = false

func main() {

	var configDir string
	var help bool
	// Define CLI flags in shrot and long form
	flag.StringVar(&configDir, "c", "", "Path where configuration file is stored (shorthand)")
	flag.StringVar(&configDir, "config", "", "Path where configuration file is stored")
	flag.BoolVar(&help, "h", false, "Show usage instructions (shorthand)")
	flag.BoolVar(&help, "help", false, "Show usage instructions")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "----------------------------------------------------------------------------------------")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Pass -c to specify where the configuration file is stored:")
		fmt.Fprintln(os.Stderr, "  -c\t\tstring\n  --config")
		fmt.Fprintln(os.Stderr, "  \tPath where configuration file is stored")
		fmt.Fprintln(os.Stderr, "\n  -h\t\tboolean\n  --help")
		fmt.Fprintln(os.Stderr, "  \tShow usage instructions")
		fmt.Fprintln(os.Stderr, "----------------------------------------------------------------------------------------")
	}
	flag.Parse()

	// Print the Help docuemntation to the terminal if user passes help flag
	if help {
		flag.Usage()
		return
	}

	// Read All Config
	allConfig, err := app.ReadAllConfig(configDir)
	if err != nil {
		log.Fatalf("function ReadAllConfig() failed: %v", err)
		return
	}

	// Read Job Config
	jobConfig, err := app.ReadJobConfig(configDir)
	if err != nil {
		log.Fatalf("function ReadJobConfig() failed: %v", err)
		return
	}

	// get sqlite db
	db, err := app.GetOrCreateSQLiteDB(allConfig, "rss")
	if err != nil {
		log.Fatalf("function GetOrCreateSQLiteDB() failed: %v", err)
	}
	defer db.Close()

	// create target table if not exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS rss (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		link TEXT,
		description TEXT,
		creator TEXT,
		pubDate TEXT
	)`)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	//----------------------------------------------------------------------------------------
	// loop over each url
	//----------------------------------------------------------------------------------------
	for i := 0; i < len(jobConfig.URLs); i++ {

		if xmlBytes, err := app.GetXML(jobConfig.URLs[i].URL); err != nil {
			log.Fatal("failed to get xml: ", err)
		} else {
			// remove any difficult strings
			xmlBytes = []byte(strings.ReplaceAll(string(xmlBytes), "atom:link", "atomlink"))

			// choose model
			var result interface{}
			if jobConfig.URLs[i].Format == 1 {
				result = &app.RssChannelFormat{}
				xml.Unmarshal(xmlBytes, &result)
				if err := xml.Unmarshal(xmlBytes, &result); err != nil {
					log.Fatal("format = 1 unmarshal error: ", err)
				}

			} else if jobConfig.URLs[i].Format == 2 {
				result = &app.FeedFormat{}
				xml.Unmarshal(xmlBytes, &result)
				if err := xml.Unmarshal(xmlBytes, &result); err != nil {
					log.Fatal("format = 2 unmarshal error: ", err)
				}
			} else {
				log.Fatal("unknown format type: ", err)
			}

			switch r := result.(type) {
			case *app.RssChannelFormat:
				for _, s := range r.Rss.Items {
					var desc string
					desc = strings.TrimSpace(s.Description)
					// Create a regular expression to match HTML tags
					re := regexp.MustCompile(`(?s)<[^>]*>`)

					// Replace all HTML tags in desc with an empty string
					desc = re.ReplaceAllString(desc, "")
					if deBug {
						fmt.Printf("Title:       %v\n", s.Title)
						fmt.Printf("Link:        %v\n", s.Link)
						fmt.Printf("Description: %v\n", desc)
						fmt.Printf("Creator:     %v\n", s.Creator)
						fmt.Printf("PubDate:     %v\n", s.PubDate)
						fmt.Println()
						fmt.Println("--------------------")
						fmt.Println()
					}
					_, err = db.Exec(`INSERT INTO rss (title, link, description, creator, pubDate) VALUES (?, ?, ?, ?, ?)`, s.Title, s.Link, desc, s.Creator, s.PubDate)
					if err != nil {
						log.Fatal("failed to insert data from rss channel data: ", err)
					}
				}
			case *app.FeedFormat:
				for _, s := range r.Rss {
					var desc string
					if strings.TrimSpace(s.Description) == "" && strings.TrimSpace(s.Content) != "" {
						desc = strings.TrimSpace(s.Content)
					} else {
						desc = strings.TrimSpace(s.Description)
					}
					// Create a regular expression to match HTML tags
					re := regexp.MustCompile(`(?s)<[^>]*>`)

					// Choose the date, leave blank if not found
					if s.PubDate == "" {
						if s.PublishedDate != "" {
							s.PubDate = s.PublishedDate
						} else if s.UpdateDate != "" {
							s.PubDate = s.UpdateDate
						}
					}

					// Replace all HTML tags in desc with an empty string
					desc = re.ReplaceAllString(desc, "")
					if deBug {
						fmt.Printf("Title:       %v\n", s.Title)
						fmt.Printf("Link:        %v\n", s.Link)
						fmt.Printf("Description: %v\n", desc)
						fmt.Printf("Creator:     %v\n", s.Creator.Name[0])
						fmt.Printf("PubDate:     %v\n", s.PubDate)
						fmt.Println()
						fmt.Println("--------------------")
						fmt.Println()
					}
					_, err = db.Exec(`INSERT INTO rss (title, link, description, creator, pubDate) VALUES (?, ?, ?, ?, ?)`, s.Title, s.Link, desc, s.Creator.Name[0], s.PubDate)
					if err != nil {
						log.Fatal("failed to insert data from feed data: ", err)
					}
				}
			default:
				log.Fatal("unexpected type")

			}
		}

		// remvoe duplicates from table
		_, err = db.Exec(`DELETE FROM rss WHERE id NOT IN (SELECT MIN(id) FROM rss GROUP BY link)`)
		if err != nil {
			log.Fatal("failed to remove duplicates from rss table: ", err)
			return
		}

	}
}
