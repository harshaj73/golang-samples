// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// [START gae_flex_analytics_track_event]

// Sample analytics demonstrates Google Analytics calls from App Engine flexible environment.
package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	uuid "github.com/gofrs/uuid"

	"google.golang.org/appengine"
)

// gaPropertyEnvar is the name of the environment variable containing a
// tracking ID for Google Analytics.
const gaPropertyEnvar = "GA_TRACKING_ID"

func main() {
	http.HandleFunc("/", handle)
	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if err := trackEvent(r, "Example", "Test action", "label", ""); err != nil {
		fmt.Fprintf(w, "Event did not track: %v", err)
		return
	}
	fmt.Fprint(w, "Event tracked.")
}

func trackEvent(r *http.Request, category, action, label string, customValue string) error {
	gaPropertyID := os.Getenv(gaPropertyEnvar)
	if gaPropertyID == "" {
		return fmt.Errorf("analytics: %s environment variable is missing", gaPropertyEnvar)
	}
	if category == "" || action == "" {
		return fmt.Errorf("analytics: category and action are required")
	}

	v := url.Values{
		"v":   {"1"},
		"tid": {gaPropertyID},
		// Anonymously identifies a particular user. See the parameter guide for
		// details:
		// https://developers.google.com/analytics/devguides/collection/protocol/v1/parameters#cid
		//
		// Depending on your application, this might want to be associated with the
		// user in a cookie.
		"cid": {uuid.Must(uuid.NewV4()).String()},
		"t":   {"event"},
		"ec":  {category},
		"ea":  {action},
		"ua":  {r.UserAgent()},
	}

	if label != "" {
		v.Set("el", label)
	}

	if customValue != "" {
		v.Set("ev", fmt.Sprintf("%s", customValue))
	}

	// NOTE: Google Analytics returns a 200, even if the request is malformed.
	_, err := http.PostForm("https://www.google-analytics.com/collect", v)
	return err
}

// [END gae_flex_analytics_track_event]
