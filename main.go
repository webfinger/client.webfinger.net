// Copyright 2013 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"webfinger.net/go/webfinger"
)

func main() {
	flag.Parse()
	http.HandleFunc("/", lookup)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Listening on %v", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func lookup(w http.ResponseWriter, r *http.Request) {
	resource := strings.TrimSpace(r.FormValue("resource"))
	if resource == "" {
		fmt.Fprint(w, "OK")
		return
	}

	// redirect old lookup URL to webfinger.net
	if strings.HasPrefix(r.URL.Path, "/lookup") {
		u := *r.URL
		u.Host = "webfinger.net"
		http.Redirect(w, r, u.String(), http.StatusFound)
		return
	}

	client := webfinger.NewClient(nil)
	logs := new(bytes.Buffer)
	client.Logger = log.New(logs, "", log.Ltime)

	jrd, err := client.Lookup(resource, nil)
	if err != nil {
		client.Logger.Printf("Error getting JRD: %v", err)
	}

	var data = struct {
		Resource string         `json:"resource"`
		JRD      *webfinger.JRD `json:"jrd"`
		Logs     string         `json:"logs"`
	}{resource, jrd, logs.String()}

	w.Header().Set("Access-Control-Allow-Origin", "https://webfinger.net")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(data)
}
