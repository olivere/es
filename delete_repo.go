package main

import (
	"log"
	"os"
)

var cmdDeleteRepo = &Command{
	Run:    runDeleteRepo,
	Usage:  "delete-repo [-f] <repo>",
	Short:  "delete repo",
	Long:   `Deletes the specified snapshot repo.`,
	ApiUrl: "http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/modules-snapshots.html#_repositories",
}

func runDeleteRepo(cmd *Command, args []string) {
	if len(args) < 1 {
		cmd.printUsage()
		os.Exit(1)
	}

	repo := args[0]

	var response struct {
		Ok     bool   `json:"ok,omitempty"`
		Ack    bool   `json:"acknowledged,omitempty"`
		Error  string `json:"error,omitempty"`
		Status int    `json:"status,omitempty"`
	}
	ESReq("DELETE", "/_snapshot/"+repo).Do(&response)
	if !force && len(response.Error) > 0 {
		log.Fatalf("Error: %v (%v)\n", response.Error, response.Status)
	}
}
