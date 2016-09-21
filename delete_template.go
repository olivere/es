package main

import (
	"log"
	"os"
)

var cmdDeleteTemplate = &Command{
	Run:    runDeleteTemplate,
	Usage:  "delete-template [-f] <template>",
	Short:  "delete template",
	Long:   `Deletes the specified template.`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/admin-indices-templates.html",
}

func init() {
	cmdDeleteTemplate.Flag.BoolVar(&force, "f", false, "force")
}

func runDeleteTemplate(cmd *Command, args []string) {
	if len(args) < 1 {
		cmd.printUsage()
		os.Exit(1)
	}

	template := args[0]

	var response struct {
		Ack bool `json:"acknowledged,omitempty"`
	}
	ESReq("DELETE", "/_template/"+template).Do(&response)
	if !force && !response.Ack {
		log.Fatalf("Error: %v\n", response)
	}
}
