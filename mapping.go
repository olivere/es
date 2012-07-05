package main

import (
	"fmt"
	"log"
	"os"
)

var cmdMapping = &Command{
	Run:   runMapping,
	Usage: "mapping <index>",
	Short: "list mapping",
	Long: `
Lists the mapping of an index.

Example:

  $ es mapping twitter
`,
}

func runMapping(cmd *Command, args []string) {
	if len(args) < 1 {
		cmd.printUsage()
		os.Exit(1)
	}

	index := args[0]

	var mapping map[string]interface{}
	body := ESReq("GET", "/"+index+"/_mapping?pretty=1").Do(&mapping)
	if error, ok := mapping["error"]; ok {
		status, _ := mapping["status"]
		log.Fatalf("Error: %v (%v)\n", error, status)
	}
	fmt.Print(body)
}
