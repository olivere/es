package main

import (
	"fmt"
	"log"
	"regexp"
)

var cmdIndices = &Command{
	Run:   runIndices,
	Usage: "indices [<pattern>]",
	Short: "list all indices",
	Long: `
Prints a list of all indices matching the specified pattern.

Example:

  $ es indices
  master
  marvel
  $ es indices 'mas.*'
  master
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/index_.html",
}

func init() {
	// parse args here, if necessary
}

func runIndices(cmd *Command, args []string) {
	var pattern = ""
	if len(args) > 0 {
		pattern = args[0]
	}

	var mappings map[string]interface{}
	ESReq("GET", "/_mapping").Do(&mappings)
	for k, _ := range mappings {
		if len(pattern) > 0 {
			matched, err := regexp.MatchString(pattern, k)
			if err != nil {
				log.Fatal("invalid pattern")
			}
			if matched {
				fmt.Println(k)
			}
		} else {
			fmt.Println(k)
		}
	}
}
