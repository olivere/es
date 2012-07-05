package main

import (
	"log"
	"os"
)

var cmdDeleteIndex = &Command{
	Run:   runDeleteIndex,
	Usage: "delete [-f] <index>",
	Short: "deletes index",
	Long: `
Deletes an index.

Example:

  $ es delete marvel
`,
}

func init() {
	cmdDeleteIndex.Flag.BoolVar(&force, "f", false, "force")
}

func runDeleteIndex(cmd *Command, args []string) {
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
	ESReq("DELETE", "/"+index).Do(&response)
	if !force && len(response.Error) > 0 {
		log.Fatalf("Error: %v (%v)\n", response.Error, response.Status)
	}
}
