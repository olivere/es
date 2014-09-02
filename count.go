package main

import (
	"fmt"
	"net/url"
	"os"
)

var cmdCount = &Command{
	Run:   runCount,
	Usage: "count <index>",
	Short: "count documents in indices",
	Long: `
Returns the number of documents in one or more indices.

Example:

  $ es count catalog-1
  $ es count "catalog-*"
`,
	ApiUrl: "http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/search-count.html",
}

func runCount(cmd *Command, args []string) {
	if len(args) < 1 {
		cmd.printUsage()
		os.Exit(1)
	}

	indices := args[0]

	values := url.Values{}
	values.Set("pretty", "true")

	var response map[string]interface{}

	body := ESReq("GET", "/"+indices+"/_count?"+values.Encode()).Do(&response)
	fmt.Print(body)
}
