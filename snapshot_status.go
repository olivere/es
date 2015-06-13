package main

import (
	"fmt"
	"log"
	"strings"
)

var cmdSnapshotStatus = &Command{
	Run:   runSnapshotStatus,
	Usage: "snapshot-status <repo> <snapshot>",
	Short: "print snapshot status",
	Long: `
View snapshot status.

Example:

  $ es snapshot-status
{"snapshots":[]}

  $ es snapshot-status nfs
{"snapshots":[]}

  $ es snapshot-status nfs twitter_1
{"snapshots":[ {"snapshot": "twitter_1", ...} ]}
`,
	ApiUrl: "http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/modules-snapshots.html#_snapshot_status",
}

func runSnapshotStatus(cmd *Command, args []string) {
	paths := []string{}
	if len(args) > 0 {
		paths = append(paths, args[0])
	}
	if len(args) > 1 {
		paths = append(paths, strings.Join(args[1:],","))
	}
	paths = append(paths, "_status")

	var response map[string]interface{}

	body := ESReq("GET", "/_snapshot/"+strings.Join(paths,"/")+"?pretty=1").Do(&response)
	if error, ok := response["error"]; ok {
		status, _ := response["status"]
		log.Fatalf("Error: %v (%v)\n", error, status)
	}
	fmt.Print(body)
}
