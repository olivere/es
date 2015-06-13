package main

import (
	"log"
	"os"
)

var cmdDeleteSnapshot = &Command{
	Run:    runDeleteSnapshot,
	Usage:  "delete-snapshot [-f] <snapshot>",
	Short:  "delete snapshot",
	Long:   `Deletes the specified snapshot.`,
	ApiUrl: "http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/modules-snapshots.html#_snapshot",
}

func runDeleteSnapshot(cmd *Command, args []string) {
	if len(args) < 2 {
		cmd.printUsage()
		os.Exit(1)
	}

	repo := args[0]
	snapshot := args[1]

	var response struct {
		Ok     bool   `json:"ok,omitempty"`
		Ack    bool   `json:"acknowledged,omitempty"`
		Error  string `json:"error,omitempty"`
		Status int    `json:"status,omitempty"`
	}
	ESReq("DELETE", "/_snapshot/"+repo+"/"+snapshot).Do(&response)
	if len(response.Error) > 0 {
		log.Fatalf("Error: %v (%v)\n", response.Error, response.Status)
	}
}
