package main

import (
	"log"
	"os"
)

var cmdRestoreSnapshot = &Command{
	Run:   runRestoreSnapshot,
	Usage: "restore-snapshot [-f] <repo> <snapshot>",
	Short: "restore snapshot",
	Long: `
Restores the specified snapshot.

Example:

  $ es restore-snapshot my_backup snapshot1
  $ es restore-snapshot -f my_backup snapshot1
`,
	ApiUrl: "http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/modules-snapshots.html#_restore",
}

func runRestoreSnapshot(cmd *Command, args []string) {
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
	data := getJsonFromStdin()
	req := ESReq("POST", "/_snapshot/"+repo+"/"+snapshot+"/_restore")
	if data != nil {
		req.SetBodyJson(data)
	}
	req.Do(&response)
	if len(response.Error) > 0 {
		log.Fatalf("Error: %v (%v)\n", response.Error, response.Status)
	}
}
