package main

import (
	"log"
)

var cmdFlushDisable = &Command{
	Run:   runFlushDisable,
	Usage: "flush-disable [index]",
	Short: "disable flush",
	Long: `
Disables flush. Useful if you want to prevent accidential flushes
from happening, e.g. while performing a backup via rsync.

Use flush-enable to revert this action.

This is basically the same as:

  curl -XPUT 'localhost:9200/_settings' -d '{
    "index" : {
      "translog.disable_flush" : "true"
    }
  }'

Example:

  $ es flush-disable
  $ es flush-disable twitter
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/admin-indices-update-settings.html",
}

func runFlushDisable(cmd *Command, args []string) {
	index := ""
	if len(args) >= 1 {
		index = args[0]
	}

	translogReq := make(map[string]interface{})
	translogReq["translog.disable_flush"] = "true"
	data := make(map[string]interface{})
	data["index"] = translogReq

	var response struct {
		Ok     bool   `json:"ok,omitempty"`
		Error  string `json:"error,omitempty"`
		Status int    `json:"status,omitempty"`
	}

	path := ""
	if index != "" {
		path += "/"+index
	}
	path += "/_settings"

	req := ESReq("PUT", path)
	req.SetBodyJson(data)
	req.Do(&response)

	if len(response.Error) > 0 {
		log.Fatalf("Error: %v (%v)\n", response.Error, response.Status)
	}
}
