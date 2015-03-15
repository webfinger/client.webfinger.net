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
	"html/template"
	"log"
	"net/http"
	"os"

	webfinger "github.com/ant0ine/go-webfinger"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

var (
	lookupTemplate = template.Must(template.ParseFiles("lookup.html"))
)

const (
	webfingerHome = "https://webfinger.net/"
)

func lookup(w http.ResponseWriter, r *http.Request) {
	flags := log.Flags()
	defer func() {
		// reset standard logger back to normal
		log.SetOutput(os.Stderr)
		log.SetFlags(flags)
	}()

	logs := new(bytes.Buffer)
	log.SetFlags(log.Ltime)
	log.SetOutput(logs)

	ctx := appengine.NewContext(r)
	client := webfinger.NewClient(urlfetch.Client(ctx))
	client.AllowHTTP = true

	var jrd string

	input := r.FormValue("resource")
	if input == "" {
		http.Redirect(w, r, webfingerHome, http.StatusFound)
		return
	}

	j, err := client.Lookup(input, nil)
	if err != nil {
		log.Printf("Error getting JRD: %v", err.Error())
	} else {
		b, err := json.MarshalIndent(j, "", "  ")
		if err != nil {
			log.Printf("Error marshalling JRD: %v", err.Error())
		} else {
			jrd = string(b)
		}
	}

	var data = struct {
		Resource string
		JRD      string
		Logs     string
	}{input, jrd, logs.String()}
	lookupTemplate.Execute(w, data)
}

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// redirect "/" to webfinger.net
		http.Redirect(w, r, webfingerHome, http.StatusFound)
	})
	http.HandleFunc("/lookup", lookup)
}
