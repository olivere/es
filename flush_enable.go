package main

import (
	"log"
)

var cmdFlushEnable = &Command{
	Run:   runFlushEnable,
	Usage: "flush-enable [index]",
	Short: "enable flush",
	Long: `
Enables flush. Useful to return to normal usage after flush has
been temporarily disabled (see flush-disable) in order to
perform a backup (e.g. via rsync).

This is basically the same as:

  curl -XPUT 'localhost:9200/_settings' -d '{
    "index" : {
      "translog.disable_flush" : "false"
    }
  }'

Example:

  $ es flush-enable
  $ es flush-enable twitter
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/admin-indices-update-settings.html",
}

func runFlushEnable(cmd *Command, args []string) {
	index := ""
	if len(args) >= 1 {
		index = args[0]
	}

	translogReq := make(map[string]interface{})
	translogReq["translog.disable_flush"] = "false"
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
