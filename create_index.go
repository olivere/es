package main

import (
	"log"
	"os"
)

var cmdCreateIndex = &Command{
	Run:   runCreateIndex,
	Usage: "create [-f] <index>",
	Short: "create index",
	Long: `
Creates an empty index.

Example:

  $ es create marvel
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/admin-indices-create-index.html",
}

func init() {
	cmdCreateIndex.Flag.BoolVar(&force, "f", false, "force")
}

func runCreateIndex(cmd *Command, args []string) {
	if len(args) < 1 {
		cmd.printUsage()
		os.Exit(1)
	}

	index := args[0]

	var response struct {
		Ok     bool   `json:"ok,omitempty"`
		Ack    bool   `json:"acknowledged,omitempty"`
		Error  string `json:"error,omitempty"`
		Status int    `json:"status,omitempty"`
	}
	ESReq("PUT", "/"+index).Do(&response)
	if !force && len(response.Error) > 0 {
		log.Fatalf("Error: %v (%v)\n", response.Error, response.Status)
	}
}
