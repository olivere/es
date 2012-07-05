package main

import (
	"fmt"
	"log"
	"regexp"
)

var cmdTemplates = &Command{
	Run:   runTemplates,
	Usage: "templates [<pattern>]",
	Short: "list all templates",
	Long: `
Prints a list of all templates matching the specified pattern.

Example:

  $ es templates
  master-template
  dummy-template
  $ es templates 'mas.*'
  master-template
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/admin-indices-templates.html",
}

func init() {
	// parse args here, if necessary
}

func runTemplates(cmd *Command, args []string) {
	var pattern = ""
	if len(args) > 0 {
		pattern = args[0]
	}

	type metadata struct {
		Templates map[string]interface{} `json:"templates,omitempty"`
	}

	var response struct {
		Metadata metadata `json:"metadata,omitempty"`
	}

	ESReq("GET", "/_cluster/state").Do(&response)

	for k, _ := range response.Metadata.Templates {
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
