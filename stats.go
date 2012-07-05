package main

import (
	"fmt"
	"log"
)

var cmdStats = &Command{
	Run:   runStats,
	Usage: "stats [index]",
	Short: "prints index statistics",
	Long: `
Lists detailed information such as the number of documents
in an index etc.

Example:

  $ es stats
  $ es stats twitter
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/admin-indices-stats.html",
}

func runStats(cmd *Command, args []string) {
	index := ""
	if len(args) >= 1 {
		index = args[0]
	}

	var response map[string]interface{}

	var body string
	if len(index) > 0 {
		body = ESReq("GET", "/"+index+"/_stats?pretty=1").Do(&response)
	} else {
		body = ESReq("GET", "/_stats?pretty=1").Do(&response)
	}

	if error, ok := response["error"]; ok {
		status, _ := response["status"]
		log.Fatalf("Error: %v (%v)\n", error, status)
	}
	fmt.Print(body)
}
