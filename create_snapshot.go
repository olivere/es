package main

import (
	"log"
	"net/url"
	"os"
)

var cmdCreateSnapshot = &Command{
	Run:   runCreateSnapshot,
	Usage: "create-snapshot <repo> <snapshot>",
	Short: "create snapshot",
	Long: `
Creates a new snapshot.

Example:

  $ es create-snapshot nfs logstash_1 < nfs-snapshot.json
`,
	ApiUrl: "http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/modules-snapshots.html#_snapshots",
}

var waitForCompletion bool

func init() {
	cmdCreateSnapshot.Flag.BoolVar(&waitForCompletion, "w", false, "wait")
}

func runCreateSnapshot(cmd *Command, args []string) {
	if len(args) < 2 {
		cmd.printUsage()
		os.Exit(1)
	}

	repo := args[0]
	snapshot := args[1]
	values := url.Values{}
	if waitForCompletion {
		values.Set("wait_for_completion", "true")
	}

	var response struct {
		Ok     bool   `json:"ok,omitempty"`
		Ack    bool   `json:"acknowledged,omitempty"`
		Error  string `json:"error,omitempty"`
		Status int    `json:"status,omitempty"`
	}

	data := getJsonFromStdin()
	req := ESReq("PUT", "/_snapshot/"+repo+"/"+snapshot+"?"+values.Encode())
	if data != nil {
		req.SetBodyJson(data)
	}
	req.Do(&response)
	if len(response.Error) > 0 {
		log.Fatalf("Error: %v (%v)\n", response.Error, response.Status)
	}
}
