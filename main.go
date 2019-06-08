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
)

var (
	allowHTTP = flag.Bool("allow_http", false, "allow falling back to non-secure HTTP connections")
	port      = flag.Int("port", 8080, "TCP port to listen on")

	// shared webfinger.Client used to process requests
	wfClient *webfinger.Client

	// mu protects access to the log package while processing lookup requests
	mu sync.Mutex
)

func init() {
	http.HandleFunc("/", lookup)
}

func webfingerClient(_ *http.Request) *webfinger.Client {
	return wfClient
}

func main() {
	flag.Parse()

	wfClient = webfinger.NewClient(nil)
	wfClient.AllowHTTP = *allowHTTP
	wfClient.WebFistServer = ""

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
	defer func() {
		// reset standard logger back to normal
		log.SetOutput(os.Stderr)
		log.SetFlags(flags)
		mu.Unlock()
	}()

	logs := new(bytes.Buffer)
	log.SetFlags(log.Ltime)
	log.SetOutput(logs)

	var jrd string

	j, err := webfingerClient(r).Lookup(input, nil)
	if err != nil {
		log.Printf("Error getting JRD: %v", err)
	} else {
		b, err := json.MarshalIndent(j, "", "  ")
		if err != nil {
			log.Printf("Error marshalling JRD: %v", err)
		} else {
			jrd = string(b)
		}
	}

	var data = struct {
		Resource string
		JRD      string
		Logs     string
	}{input, jrd, logs.String()}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(data)
}
