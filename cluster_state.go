package main

import (
	"fmt"
	"log"
)

var cmdClusterState = &Command{
	Run:   runClusterState,
	Usage: "cluster-state",
	Short: "prints cluster state",
	Long: `
Cluster health prints a comprehensive state information of the whole cluster.

Example:

  $ es cluster-state
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/admin-cluster-state.html",
}

func runClusterState(cmd *Command, args []string) {
	var response map[string]interface{}

	body := ESReq("GET", "/_cluster/state?pretty=1").Do(&response)

	if error, ok := response["error"]; ok {
		status, _ := response["status"]
		log.Fatalf("Error: %v (%v)\n", error, status)
	}

	fmt.Print(body)
}
