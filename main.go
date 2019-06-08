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
	"sync"

	"github.com/ant0ine/go-webfinger"
	"github.com/ant0ine/go-webfinger/jrd"
)

var (
	port = flag.Int("port", 8080, "TCP port to listen on")

	// shared webfinger.Client used to process requests
	client *webfinger.Client

	// mu protects access to the log package while processing lookup requests
	mu sync.Mutex
)

func init() {
	http.HandleFunc("/", lookup)
}

func main() {
	flag.Parse()

	client = webfinger.NewClient(nil)
	client.WebFistServer = ""

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Listening on %v", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func lookup(w http.ResponseWriter, r *http.Request) {
	input := r.FormValue("resource")
	if input == "" {
		http.Error(w, "empty resource", http.StatusBadRequest)
		return
	}

	mu.Lock()
	flags := log.Flags()
	logs := new(bytes.Buffer)
	log.SetFlags(log.Ltime)
	log.SetOutput(logs)

	j, err := client.Lookup(input, nil)

	// reset standard logger back to normal
	log.SetOutput(os.Stderr)
	log.SetFlags(flags)
	mu.Unlock()

	if err != nil {
		msg := fmt.Sprintf("Error getting JRD: %v", err)
		log.Print(msg)
		http.Error(w, msg, http.StatusInternalServerError)
	}

	var data = struct {
		Resource string   `json:"resource"`
		JRD      *jrd.JRD `json:"jrd"`
		Logs     string   `json:"logs"`
	}{input, j, logs.String()}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(data)
}
