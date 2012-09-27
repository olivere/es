package main

import (
	"fmt"
	"log"
	"regexp"
)

var cmdIndexAliases = &Command{
	Run:   runIndexAliases,
	Usage: "index-aliases [<pattern>]",
	Short: "list all aliases of an index",
	Long: `
Prints a list of all aliases of an index matching the specified pattern.

Example:

  $ es index-aliases
  master -> alias1, alias2
  marvel
  $ es index-aliases 'mas.*'
  master -> alias1, alias2
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/admin-indices-aliases.html",
}

func runIndexAliases(cmd *Command, args []string) {
	var pattern = ""
	if len(args) > 0 {
		pattern = args[0]
	}

	var response map[string]interface{}

	ESReq("GET", "/_cluster/state?pretty=1").Do(&response)

	error, hasError := response["error"];
	if hasError {
		status, _ := response["status"]
		log.Fatalf("Error: %v (%v)\n", error, status)
	}

	metadata, ok := response["metadata"].(map[string]interface{})
	if !ok {
		log.Fatalf("Could not find metadata element in response\n")
	}


	indices, ok := metadata["indices"].(map[string]interface{})
	if ok {
		for indexName, _ := range indices {
			if len(pattern) > 0 {
				matched, err := regexp.MatchString(pattern, indexName)
				if err != nil {
					log.Fatal("unvalid pattern")
				}
				if matched {
					printAliasesOfIndex(indexName, indices)
				}
			} else {
				printAliasesOfIndex(indexName, indices)
			}
		}
	} else {
		// No indices
	}
}

func printAliasesOfIndex(indexName string, indices map[string]interface{}) {
	fmt.Printf(indexName)

	if index, ok := indices[indexName].(map[string]interface{}); ok {
		aliases, ok := index["aliases"].([]interface{})
		if ok {
			for i, alias := range aliases {
				var aliasName = alias.(string)
				if i == 0 {
					fmt.Printf(" -> %s", aliasName)
				} else {
					fmt.Printf(", %s", aliasName)
				}
			}
		}
	}

	fmt.Println()
}

