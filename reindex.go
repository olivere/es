package main

import (
	"fmt"
	"github.com/olivere/elastic"
	"log"
	"net/http"
	_ "net/http/httputil"
	"os"
	"time"
)

var cmdReindex = &Command{
	Run:   runReindex,
	Usage: "reindex [-v] <source> <target>",
	Short: "reindex one index to another index",
	Long: `
The reindex command takes the documents from the source index
and bulk imports them to the specified target index.

This is quite handy if you want to change the settings of an
index and won't lose any data.

Use the -v/--verbose flag to print progress.

Example:

  $ es reindex -v twitter twitter-snapshot
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/bulk.html",
}

func init() {
	cmdReindex.Flag.BoolVar(&verbose, "v", false, "verbose")
}

func runReindex(cmd *Command, args []string) {
	if len(args) < 2 {
		cmd.printUsage()
		os.Exit(1)
	}

	source := args[0]
	target := args[1]

	// Get a client
	client, err := elastic.NewClient(http.DefaultClient, esUrl)
	if err != nil {
		log.Fatal("Error: %v", err)
	}

	// Create target index, if not already exists
	exists, err := client.IndexExists(target).Do()
	if err != nil {
		log.Fatal("Error: %v", err)
	}
	if !exists {
		_, err := client.CreateIndex(target).Do()
		if err != nil {
			log.Fatal("Error: %v", err)
		}

		// Flush
		client.Flush().Index(target).Do()
	}

	// Scan through source
	bulkSize := 1000
	bulk := client.Bulk().Index(target).DebugOnError(true)
	inserted := 0

	cursor, err := client.Scan(source).Do()
	if err != nil {
		log.Fatal("Error: %v", err)
	}

	// Iterate
	for {
		sr, err := cursor.Next()
		if err == elastic.EOS {
			break
		}
		if err != nil {
			log.Fatal("Error: %v", err)
		}

		if sr.Hits != nil && sr.Hits.Hits != nil {
			for _, hit := range sr.Hits.Hits {
				indexReq := &elastic.BulkIndexRequest{
					Id:   hit.Id,
					Type: hit.Type,
					Data: string(*hit.Source),
				}
				bulk.Add(indexReq)

				inserted += 1

				if verbose {
					fmt.Fprintf(os.Stderr, "Reindexing %9d\r", inserted)
				}

				if bulk.NumberOfActions() >= bulkSize {
					bulk, err = commitBulk(client, bulk, target)
					if err != nil {
						log.Fatal("Error: %v", err)
					}
				}
			}
		}
	}

	if bulk.NumberOfActions() > 0 {
		bulk, err = commitBulk(client, bulk, target)
		if err != nil {
			log.Fatal("Error: %v", err)
		}
	}

	// Flush
	_, err = client.Flush().Index(target).Do()
	if err != nil {
		log.Fatal("Error: %v", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Reindexed %9d    \n", inserted)
	}
}

// commitBulk commits the bulk actions to Elasticsearch. When an error
// occurs, it retries the bulk operation several times before failing.
func commitBulk(client *elastic.Client, bulk *elastic.BulkService, indexName string) (*elastic.BulkService, error) {
	sleepSeconds := int64(1)
	numRetries := 5

	for {
		_, err := bulk.Do()
		if err != nil {
			numRetries -= 1
			if numRetries < 0 {
				return bulk, err
			}
			time.Sleep(time.Duration(sleepSeconds) * time.Second)
			sleepSeconds *= 2
		} else {
			break
		}
	}

	// Return a new bulk request
	bulk = client.Bulk().Index(indexName).DebugOnError(true)

	return bulk, nil
}
