package main

import (
	"fmt"
	"log"
	"regexp"
	"sort"

	"github.com/olivere/elastic"
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

	// Get a client
	client, err := elastic.NewClient(elastic.SetURL(esUrl))
	if err != nil {
		log.Fatal(err)
	}
	indices, err := client.IndexNames()
	if err != nil {
		log.Fatal(err)
	}

	// Sort by default
	sort.Strings(indices)

	for _, index := range indices {
		if len(pattern) > 0 {
			matched, err := regexp.MatchString(pattern, index)
			if err != nil {
				log.Fatal("invalid pattern")
			}
			if matched {
				fmt.Println(index)
			}
		} else {
			fmt.Println(index)
		}
	}
}
