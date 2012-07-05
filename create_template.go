package main

import (
	"bufio"
	"encoding/json"
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

	// parse Stdin into JSON
	var data interface{}
	reader := bufio.NewReader(os.Stdin)
	if err := json.NewDecoder(reader).Decode(&data); err != nil {
		log.Fatal("invalid json\n")
	}

	req := ESReq("PUT", "/_template/"+template)
	req.SetBodyJson(data)
	req.Do(&response)
	if len(response.Error) > 0 {
		log.Fatalf("Error: %v (%v)\n", response.Error, response.Status)
	}
}
