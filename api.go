package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var cmdApi = &Command{
	Usage:  "api <command>",
	Short:  "show API for command",
	Long:   `Opens browser with API details for the specified command.`,
	ApiUrl: "",
}

func init() {
	cmdApi.Run = runApi
}

func runApi(cmd *Command, args []string) {
	if len(args) == 0 {
		printUsage()
		return
	}

	if len(args) != 1 {
		log.Fatal("too many arguments")
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && len(cmd.ApiUrl) > 0 {
			// open browser
			openBrowser := "open"
			if _, err := exec.LookPath("xdg-open"); err == nil {
				openBrowser = "xdg-open"
			}
			exec.Command(openBrowser, cmd.ApiUrl).Start()
			return
		}
	}

	fmt.Fprintf(os.Stderr, "No API information for command: %s\n", args[0])
	os.Exit(2)
}

// Generic API for accessing ElasticSearch server starts here.

type Request http.Request

func ESReq(method, path string) *Request {
	req, err := http.NewRequest(method, esUrl+path, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("User-Agent", "es/"+Version+" ("+runtime.GOOS+"-"+runtime.GOARCH+")")
	req.Header.Add("Accept", "application/json")
	return (*Request)(req)
}

func (r *Request) SetBodyJson(data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	r.SetBody(bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
}

func (r *Request) SetBody(body io.Reader) {
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	r.Body = rc
	if body != nil {
		switch v := body.(type) {
		case *strings.Reader:
			r.ContentLength = int64(v.Len())
		case *bytes.Buffer:
			r.ContentLength = int64(v.Len())
		}
	}
}

func (r *Request) Do(v interface{}) string {
	res := checkResponse(http.DefaultClient.Do((*http.Request)(r)))
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if v != nil {
		jsonErr := json.Unmarshal(bodyBytes, v)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
	}

	return string(bodyBytes)
}

func checkResponse(res *http.Response, err error) *http.Response {
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode == 401 {
		log.Fatal("Unauthorized")
	}
	if res.StatusCode == 403 {
		log.Fatal("Unauthorized")
	}
	if res.StatusCode < 200 && res.StatusCode > 299 {
		log.Fatal("Unexpected error: ", res.Status)
	}
	return res
}
