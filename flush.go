package main

import (
	"log"
)

var cmdFlush = &Command{
	Run:   runFlush,
	Usage: "flush [index]",
	Short: "flushes indices",
	Long: `
Flushes data to the index storage, frees memory and clears the
internal transaction log. You often need to flush before updating
ElasticSearch to a new version.

Example:

  $ es flush
  $ es flush twitter
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/admin-indices-flush.html",
}

func runFlush(cmd *Command, args []string) {
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
		ESReq("POST", "/"+index+"/_flush").Do(&response)
	} else {
		ESReq("POST", "/_flush").Do(&response)
	}

	if len(response.Error) > 0 {
		log.Fatalf("Error: %v (%v)\n", response.Error, response.Status)
	}
}
