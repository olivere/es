package main

import (
	"log"
	"os"
)

var cmdCloseIndex = &Command{
	Run:   runCloseIndex,
	Usage: "close [-f] <index>",
	Short: "close index",
	Long: `
Closes an index.

Example:

  $ es close marvel
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/admin-indices-open-close.html",
}

func init() {
	cmdCloseIndex.Flag.BoolVar(&force, "f", false, "force")
}

func runCloseIndex(cmd *Command, args []string) {
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
	ESReq("POST", "/"+index+"/_close").Do(&response)
	if !force && len(response.Error) > 0 {
		log.Fatalf("Error: %v (%v)\n", response.Error, response.Status)
	}
}
