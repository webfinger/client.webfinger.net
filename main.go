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
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ant0ine/go-webfinger"
)

var (
	allowHTTP = flag.Bool("allow_http", false, "allow falling back to non-secure HTTP connections")
	port      = flag.Int("port", 8080, "TCP port to listen on")

	// shared webfinger.Client used to process requests
	wfClient *webfinger.Client
)

func init() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
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
