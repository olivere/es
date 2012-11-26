package main

import (
	"bufio"
	"fmt"
	"log"
	"io"
	_ "net/http"
	_ "net/http/httputil"
	"os"
	"strings"
)

var cmdBulk = &Command{
	Run:   runBulk,
	Usage: "bulk [-v] <index>",
	Short: "bulk import",
	Long: `
The bulk API enables many index/delete operations in a single API call.

Use the -v/--verbose flag to print progress.

Notice that if the bulk data file specifies the index name, 
the argument specified on the command line will have no effect.

Example:

  $ es bulk -v twitter < twitter-data.json
`,
	ApiUrl: "http://www.elasticsearch.org/guide/reference/api/bulk.html",
}

func init() {
	cmdBulk.Flag.BoolVar(&verbose, "v", false, "verbose")
}

func bulkReader(input chan string, done chan bool) {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString(byte('\n'))
		if err == io.EOF {
			done <- true
			return
		}
		if err != nil {
			log.Fatal("error reading from input\n")
		}
		input <- strings.TrimRight(line, "\n")
	}
}

func bulkCommitter(index string, input chan string, done chan bool) {
	// collect up to bulkSize lines before committing
	numLines := 0
	bulkSize := 1000
	lines := make([]string, 0, bulkSize)

	for {
		select {
		case line := <-input:
			numLines++
			lines = append(lines, line)
			if numLines%bulkSize == 0 {
				bulkCommit(index, lines, numLines)
				lines = make([]string, 0, bulkSize)
			}
		case <-done:
			if numLines > 0 {
				bulkCommit(index, lines, numLines)
				lines = make([]string, 0, bulkSize)
			}
			return
		}
	}

	if verbose {
		fmt.Fprintln(os.Stderr, "\n")
	}
}

func bulkCommit(index string, lines []string, numLines int) {
	if verbose {
		fmt.Fprintf(os.Stderr, "Bulk importing %9d\r", numLines)
	}

	type itemsData struct {
		Index   string `json:"_index,omitempty"`
		Type    string `json:"_type,omitempty"`
		Id      string `json:"_id,omitempty"`
		Version string `json:"_version,omitempty"`
		Ok      bool   `json:"ok,omitempty"`
	}

	type items struct {
		Create itemsData `json:"create,omitempty"`
		Delete itemsData `json:"delete,omitempty"`
	}

	var response struct {
		Took  int     `json:"took,omitempty"`
		Items []items `json:"items,omitempty"`
	}

	bodyOfLines := strings.Join(lines, "\n") + "\n"

	req := ESReq("POST", "/"+index+"/_bulk")
	req.SetBodyString(bodyOfLines)
	req.Do(&response)

	// dump, _ := httputil.DumpRequest((*http.Request)(req), true)
	// log.Print(string(dump))

	//body := req.Do(&response)
	//if len(response.Error) > 0 {
	//	log.Fatalf("Error: %v (%v)\n", response.Error, response.Status)
	//}
	//log.Println(body)
}

func runBulk(cmd *Command, args []string) {
	if len(args) < 1 {
		cmd.printUsage()
		os.Exit(1)
	}

	index := args[0]

	// create a channel that reads line by line from Stdin
	done := make(chan bool)
	input := make(chan string)
	go bulkReader(input, done)

	bulkCommitter(index, input, done)
}
