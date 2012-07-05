package main

import (
	"fmt"
	"log"
)

var cmdClusterHealth = &Command{
	Run:   runClusterHealth,
	Usage: "cluster-health",
	Short: "prints cluster health",
	Long: `
Cluster health prints a very simple status on the health of the cluster.

Example:

  $ es cluster-health
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/admin-cluster-health.html",
}

func runClusterHealth(cmd *Command, args []string) {
	var response map[string]interface{}

	body := ESReq("GET", "/_cluster/health?pretty=1").Do(&response)

	if error, ok := response["error"]; ok {
		status, _ := response["status"]
		log.Fatalf("Error: %v (%v)\n", error, status)
	}

	fmt.Print(body)
}
