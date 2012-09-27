package main

import (
	"fmt"
	"log"
	"regexp"
)

var cmdAliases = &Command{
	Run:   runAliases,
	Usage: "aliases [<pattern>]",
	Short: "list all aliases",
	Long: `
Prints a list of all aliases matching the specified pattern.

Example:

  $ es aliases
  master
  marvel
  $ es aliases 'mas.*'
  master
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/admin-indices-aliases.html",
}

var printIndex bool

func init() {
	cmdAliases.Flag.BoolVar(&printIndex, "i", false, "index")
}

func runAliases(cmd *Command, args []string) {
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
			index, ok := indices[indexName].(map[string]interface{})
			if ok {
				aliases, ok := index["aliases"].([]interface{})
				if ok {
					for _, alias := range aliases {
						var aliasName = alias.(string)
						if len(pattern) > 0 {
							matched, err := regexp.MatchString(pattern, aliasName)
							if err != nil {
								log.Fatal("unvalid pattern")
							}
							if matched {
								if printIndex {
									fmt.Printf("%s -> %s\n", aliasName, indexName)
								} else {
									fmt.Println(aliasName)
								}
							}
						} else {
							if printIndex {
								fmt.Printf("%s -> %s\n", aliasName, indexName)
							} else {
								fmt.Println(aliasName)
							}
						}
					}
				}
			}
		}
	} else {
		// No indices
	}
}
