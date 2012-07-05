package main

import (
	"fmt"
	"os"
	"net/url"
)

var cmdSearch = &Command{
	Run:   runSearch,
	Usage: "search <index> <query>",
	Short: "searches an index",
	Long: `
Performs a very basic search on an index via the request URI API.

Example:

  $ es search twitter "user:kimchy"
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/search/uri-request.html",
}

func runSearch(cmd *Command, args []string) {
	if len(args) < 2 {
		cmd.printUsage()
		os.Exit(1)
	}

	index := args[0]
	query := args[1]

	values := url.Values{}
	values.Set("pretty", "true")
	values.Set("q", query)

	var response map[string]interface{}

	body := ESReq("GET", "/"+index+"/_search?"+values.Encode()).Do(&response)
	fmt.Print(body)
}
