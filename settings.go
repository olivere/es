package main

import (
	"fmt"
	"log"
)

var cmdSettings = &Command{
	Run:   runSettings,
	Usage: "settings [index]",
	Short: "prints index settings",
	Long: `
Prints index settings.

Example:

  $ es settings
  $ es settings twitter
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/admin-indices-get-settings.html",
}

func runSettings(cmd *Command, args []string) {
	index := ""
	if len(args) >= 1 {
		index = args[0]
	}

	var response map[string]interface{}

	var body string
	if len(index) > 0 {
		body = ESReq("GET", "/"+index+"/_settings?pretty=1").Do(&response)
	} else {
		body = ESReq("GET", "/_settings?pretty=1").Do(&response)
	}

	if error, ok := response["error"]; ok {
		status, _ := response["status"]
		log.Fatalf("Error: %v (%v)\n", error, status)
	}
	fmt.Print(body)
}
