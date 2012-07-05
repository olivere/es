package main

import (
	"fmt"
	"log"
	"os"
)

var cmdTemplate = &Command{
	Run:   runTemplate,
	Usage: "template <name>",
	Short: "print template",
	Long: `
Prints details of the specified template.

Example:

  $ es template dummy
  {
  }
`,
}

func init() {
	// parse args here, if necessary
}

func runTemplate(cmd *Command, args []string) {
	if len(args) < 1 {
		cmd.printUsage()
		os.Exit(1)
	}

	template := args[0]

	var response map[string]interface{}

	body := ESReq("GET", "/_template/"+template+"?pretty=1").Do(&response)
	if error, ok := response["error"]; ok {
		status, _ := response["status"]
		log.Fatalf("Error: %v (%v)\n", error, status)
	}
	fmt.Print(body)
}
