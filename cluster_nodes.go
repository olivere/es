package main

import (
	"fmt"
	"log"
)

var cmdClusterNodes = &Command{
	Run:   runClusterNodes,
	Usage: "cluster-nodes",
	Short: "prints cluster nodes information",
	Long: `
Cluster health prints information about all nodes in the cluster.

Example:

  $ es cluster-nodes
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/admin-cluster-nodes-info.html",
}

func runClusterNodes(cmd *Command, args []string) {
	var response map[string]interface{}

	params := "pretty=1" +
		"&settings=true" +
		"&os=true" +
		"&process=true" +
		"&jvm=true" +
		"&thread_pool=true" +
		"&network=true" +
		"&transport=true" +
		"&http=true"

	body := ESReq("GET", "/_cluster/nodes?"+params).Do(&response)

	if error, ok := response["error"]; ok {
		status, _ := response["status"]
		log.Fatalf("Error: %v (%v)\n", error, status)
	}

	fmt.Print(body)
}
