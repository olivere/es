package main

import (
	"log"
	"net/url"
)

var cmdOptimize = &Command{
	Run:   runOptimize,
	Usage: "optimize [-r -f] [index]",
	Short: "optimizes indices",
	Long: `
Optimize basically optimizes the index for faster search operations.

Use -r, -r=true, -r=false to enable/disable refresh (enabled by default).
Use -f, -f=true, -f=false to enable/disable flush (enabled by default).

Example:

  $ es optimize
  $ es optimize -f=false twitter
  $ es optimize twitter
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/admin-indices-optimize.html",
}

func init() {
	cmdOptimize.Flag.BoolVar(&flush, "f", true, "flush")
	cmdOptimize.Flag.BoolVar(&refresh, "r", true, "refresh")
}

func runOptimize(cmd *Command, args []string) {
	index := ""
	if len(args) >= 1 {
		index = args[len(args)-1]
	}

	values := url.Values{}
	if flush {
		values.Set("flush", "1")
	} else {
		values.Set("flush", "0")
	}
	if refresh {
		values.Set("refresh", "1")
	} else {
		values.Set("refresh", "0")
	}

	var response struct {
		Ok     bool   `json:"ok,omitempty"`
		Error  string `json:"error,omitempty"`
		Status int    `json:"status,omitempty"`
	}

	if len(index) > 0 {
		ESReq("POST", "/"+index+"/_optimize?"+values.Encode()).Do(&response)
	} else {
		ESReq("POST", "/_optimize?"+values.Encode()).Do(&response)
	}

	if len(response.Error) > 0 {
		log.Fatalf("Error: %v (%v)\n", response.Error, response.Status)
	}
}
