package main

import (
	"bytes"
	"encoding/json"
	"log"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
)

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

func (r *Request) Do(v interface{}) {
	res := checkResponse(http.DefaultClient.Do((*http.Request)(r)))
	defer res.Body.Close()

	if v != nil {
		err := json.NewDecoder(res.Body).Decode(v)
		if err != nil {
			log.Fatal(err)
		}
	}
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
