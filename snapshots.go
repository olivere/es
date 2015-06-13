package main

import (
	"fmt"
	"log"
	"regexp"
	"os"
)

var cmdSnapshots = &Command{
	Run:   runSnapshots,
	Usage: "snapshots <repo> [<pattern>]",
	Short: "list all snapshots in a repo",
	Long: `
Prints a list of all snapshots in a repo matching the specified pattern.

Example:

  $ es snapshots nfs
  logstash_1
  fluentd_1
  $ es snapshots nfs 'logstash.*'
  logstash_1
`,
	ApiUrl: "http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/modules-snapshots.html#_snapshot",
}

func init() {
	// parse args here, if necessary
}

func runSnapshots(cmd *Command, args []string) {
	if len(args) < 1 {
		cmd.printUsage()
		os.Exit(1)
	}

	repo := args[0]

	var pattern = ""
	if len(args) > 1 {
		pattern = args[1]
	}

	type snapshot struct {
		Snapshot string `json:"snapshot"`
	}
	var response struct {
		Snapshots []snapshot `json:"snapshots,omitempty"`
		Error     string   `json:"error,omitempty"`
		Status    int      `json:"status,omitempty"`
	}
	ESReq("GET", "/_snapshot/"+repo+"/_all").Do(&response)

	if len(response.Error) > 0 {
		log.Fatalf("Error: %v (%v)\n", response.Error, response.Status)
	} else {
		for _, snapshot := range response.Snapshots {
			if len(pattern) > 0 {
				matched, err := regexp.MatchString(pattern, snapshot.Snapshot)
				if err != nil {
					log.Fatal("invalid pattern")
				}
				if matched {
					fmt.Println(snapshot.Snapshot)
				}
			} else {
				fmt.Println(snapshot.Snapshot)
			}
		}
	}
}
