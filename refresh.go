package main

import (
	"log"
)

var cmdRefresh = &Command{
	Run:   runRefresh,
	Usage: "refresh [index]",
	Short: "refreshes indices",
	Long: `
Refresh allows to explicitly refresh one or more index, 
making all operations performed since the last refresh 
available for search.

Example:

  $ es refresh
  $ es refresh twitter
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/admin-indices-refresh.html",
}

func runRefresh(cmd *Command, args []string) {
	index := ""
	if len(args) >= 1 {
		index = args[0]
	}

	var response struct {
		Ok     bool   `json:"ok,omitempty"`
		Error  string `json:"error,omitempty"`
		Status int    `json:"status,omitempty"`
	}

	if len(index) > 0 {
		ESReq("POST", "/"+index+"/_refresh").Do(&response)
	} else {
		ESReq("POST", "/_refresh").Do(&response)
	}

	if len(response.Error) > 0 {
		log.Fatalf("Error: %v (%v)\n", response.Error, response.Status)
	}
}
