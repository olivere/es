package main

import (
	"fmt"
	"log"
	"os"
)

var cmdSnapshot = &Command{
	Run:   runSnapshot,
	Usage: "snapshot <repo> <snapshot>",
	Short: "print snapshot",
	Long: `
Print a snapshot.

Example:

  $ es snapshot nfs logstash_1
{}
`,
	ApiUrl: "http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/modules-snapshots.html#_snapshots",
}

func runSnapshot(cmd *Command, args []string) {
	if len(args) < 2 {
		cmd.printUsage()
		os.Exit(1)
	}

	repo := args[0]
	snapshot := args[1]

	var response map[string]interface{}

	body := ESReq("GET", "/_snapshot/"+repo+"/"+snapshot+"?pretty=1").Do(&response)
	if error, ok := response["error"]; ok {
		status, _ := response["status"]
		log.Fatalf("Error: %v (%v)\n", error, status)
	}
	fmt.Print(body)
}
