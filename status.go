package main

import (
	"fmt"
	"log"
)

var cmdStatus = &Command{
	Run:   runStatus,
	Usage: "status [index]",
	Short: "prints index status",
	Long: `
Lists status details of all indices or the specified index.

Example:

  $ es status
  $ es status twitter
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/admin-indices-status.html",
}

func runStatus(cmd *Command, args []string) {
	index := ""
	if len(args) >= 1 {
		index = args[0]
	}

	var response map[string]interface{}

	var body string
	if len(index) > 0 {
		body = ESReq("GET", "/"+index+"/_status?pretty=1").Do(&response)
	} else {
		body = ESReq("GET", "/_status?pretty=1").Do(&response)
	}

	if error, ok := response["error"]; ok {
		status, _ := response["status"]
		log.Fatalf("Error: %v (%v)\n", error, status)
	}
	fmt.Print(body)
}
