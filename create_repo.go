package main

import (
	"log"
	"os"
)

var cmdCreateRepo = &Command{
	Run:   runCreateRepo,
	Usage: "create-repo <repo>",
	Short: "create snapshot repo",
	Long: `
Creates a new snapshot repo.

Example:

  $ es create-repo nfs < nfs-repo.json
`,
	ApiUrl: "http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/modules-snapshots.html#_repositories",
}

func runCreateRepo(cmd *Command, args []string) {
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

	data := getJsonFromStdin()
	req := ESReq("PUT", "/_snapshot/"+repo)
	req.SetBodyJson(data)
	req.Do(&response)
	if len(response.Error) > 0 {
		log.Fatalf("Error: %v (%v)\n", response.Error, response.Status)
	}
}
