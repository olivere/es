package main

import (
	"log"
	"os"
)

var cmdOpenIndex = &Command{
	Run:   runOpenIndex,
	Usage: "open [-f] <index>",
	Short: "open index",
	Long: `
Opens an index.

Example:

  $ es open marvel
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/admin-indices-open-close.html",
}

func init() {
	cmdOpenIndex.Flag.BoolVar(&force, "f", false, "force")
}

func runOpenIndex(cmd *Command, args []string) {
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
	ESReq("POST", "/"+index+"/_open").Do(&response)
	if !force && len(response.Error) > 0 {
		log.Fatalf("Error: %v (%v)\n", response.Error, response.Status)
	}
}
