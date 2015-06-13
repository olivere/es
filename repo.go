package main

import (
	"fmt"
	"log"
	"os"
)

var cmdRepo = &Command{
	Run:   runRepo,
	Usage: "repo <name>",
	Short: "print snapshot repo",
	Long: `
Prints details of the specified snapshot repo.

Example:

  $ es repo dummy
  {
  }
`,
	ApiUrl: "http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/modules-snapshots.html#_repositories",
}

func init() {
	// parse args here, if necessary
}

func runRepo(cmd *Command, args []string) {
	if len(args) < 1 {
		cmd.printUsage()
		os.Exit(1)
	}

	repo := args[0]

	var response map[string]interface{}

	body := ESReq("GET", "/_snapshot/"+repo+"?pretty=1").Do(&response)
	if error, ok := response["error"]; ok {
		status, _ := response["status"]
		log.Fatalf("Error: %v (%v)\n", error, status)
	}
	fmt.Print(body)
}
