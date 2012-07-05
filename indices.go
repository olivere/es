package main

import (
	"fmt"
)

var cmdIndices = &Command{
	Run:   runIndices,
	Usage: "indices",
	Short: "list all indices",
	Long:  `Prints a list of all indices.`,
}

func init() {
	// parse args here, if necessary
}

func runIndices(cmd *Command, args []string) {
	var mappings map[string]interface{}
	ESReq("GET", "/_mapping").Do(&mappings)
	for k, _ := range mappings {
		fmt.Println(k)
	}
}
