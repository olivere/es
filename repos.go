package main

import (
	"fmt"
	"log"
	"regexp"
)

var cmdRepos = &Command{
	Run:   runRepos,
	Usage: "repos [<pattern>]",
	Short: "list all repos",
	Long: `
Prints a list of all repos matching the specified pattern.

Example:

  $ es repos
  nfs
  s3
  $ es repos 's3.*'
  s3
`,
	ApiUrl: "http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/modules-snapshots.html#_repositories",
}

func init() {
	// parse args here, if necessary
}

func runRepos(cmd *Command, args []string) {
	var pattern = ""
	if len(args) > 0 {
		pattern = args[0]
	}

	var repos map[string]interface{}
	ESReq("GET", "/_snapshot").Do(&repos)
	for k, _ := range repos {
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
