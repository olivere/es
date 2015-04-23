package main

import (
	"fmt"
	"log"
	_ "net/http/httputil"
	"os"

	"github.com/olivere/elastic"
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

You can also copy indices from one cluster to another by using the -source
and -target options.

Use the -v flag to print progress.

Example:

  $ es reindex -v twitter twitter-snapshot
  $ es reindex -v -source=http://cluster1:9200 -target=http://cluster2:9200 twitter twitter-snapshot
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/bulk.html",
}

var (
	sourceURL, targetURL string
)

func init() {
	cmdReindex.Flag.BoolVar(&verbose, "v", false, "verbose")
	cmdReindex.Flag.StringVar(&sourceURL, "source", "", "URL of source cluster")
	cmdReindex.Flag.StringVar(&targetURL, "target", "", "URL of target cluster")
}

func runReindex(cmd *Command, args []string) {
	if len(args) < 2 {
		cmd.printUsage()
		os.Exit(1)
	}

	sourceIndex := args[0]
	targetIndex := args[1]

	if sourceURL == "" {
		sourceURL = esUrl
	}
	if targetURL == "" {
		targetURL = esUrl
	}

	// Get a client
	sourceClient, err := elastic.NewClient(elastic.SetURL(sourceURL))
	if err != nil {
		log.Fatalf("%v", err)
	}
	targetClient := sourceClient
	if sourceURL != targetURL {
		targetClient, err = elastic.NewClient(elastic.SetURL(targetURL))
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	// Check if source index exists. Stop if it doesn't.
	exists, err := sourceClient.IndexExists(sourceIndex).Do()
	if err != nil {
		log.Fatalf("%v", err)
	}
	if !exists {
		os.Exit(0)
	}

	// Create target index, if not already exists
	exists, err = targetClient.IndexExists(targetIndex).Do()
	if err != nil {
		log.Fatalf("%v", err)
	}
	if !exists {
		_, err := targetClient.CreateIndex(targetIndex).Do()
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	// Progress callback
	progress := func(current, total int64) {
		if verbose {
			var percent int64
			if total > 0 {
				percent = 100 * current / total
			}
			fmt.Fprintf(os.Stderr, "Reindexing %9d of %9d (%3d%%)\r", current, total, percent)
		}
	}

	// Use the Elastic Reindexer
	ix := elastic.NewReindexer(sourceClient, sourceIndex, targetIndex)
	ix = ix.TargetClient(targetClient)
	ix = ix.BulkSize(1000)
	ix = ix.Scroll("5m")
	ix = ix.Progress(progress)
	ix = ix.StatsOnly(true)
	res, err := ix.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Fprintf(os.Stderr, "%d successful and %d failed request(s)\n", res.Success, res.Failed)
}
