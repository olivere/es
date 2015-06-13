package main

import (
	"log"
	"os"
)

var cmdCreateTemplate = &Command{
	Run:   runCreateTemplate,
	Usage: "create-template <template>",
	Short: "create template",
	Long: `
Creates a new template.

Example:

  $ es create-template marvel < marvel-template.json
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/admin-indices-templates.html",
}

func runCreateTemplate(cmd *Command, args []string) {
	if len(args) < 1 {
		cmd.printUsage()
		os.Exit(1)
	}

	template := args[0]

	var response struct {
		Ok     bool   `json:"ok,omitempty"`
		Ack    bool   `json:"acknowledged,omitempty"`
		Error  string `json:"error,omitempty"`
		Status int    `json:"status,omitempty"`
	}

	data := getJsonFromStdin()
	req := ESReq("PUT", "/_template/"+template)
	req.SetBodyJson(data)
	req.Do(&response)
	if len(response.Error) > 0 {
		log.Fatalf("Error: %v (%v)\n", response.Error, response.Status)
	}
}
