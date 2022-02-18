package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Struct to hold basic user data
type User struct {
	// First Name
	First string
	// Last Name
	Last string
}

// Handler function
func handlePostUser(t *testing.T) func(http.ResponseWriter, *http.Request) {
	// Returns a response writer handler
	return func(w http.ResponseWriter, r *http.Request) {
		// call a deferred reader function
		defer func(r io.ReadCloser) {
			// copy user input into our io reader
			_, _ = io.Copy(ioutil.Discard, r)
			// close reader
			_ = r.Close()
		}(r.Body)

		// test if given request matches the POST method
		if r.Method != http.MethodPost {
			// if not return error
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		// decode reader using json decoder
		var u User
		err := json.NewDecoder(r.Body).Decode(&u)
		// test for error
		if err != nil {
			// print the error
			t.Error(err)
			// print failure to decode
			http.Error(w, "Decode Failed", http.StatusBadRequest)
			return
		}

		// write http status header
		w.WriteHeader(http.StatusAccepted)
	}
}

// function to test a POST request to the server
func TestPostUser(t *testing.T) {
	// establish a server for our handler function
	ts := httptest.NewServer(http.HandlerFunc(handlePostUser(t)))
	// defer the closing of the server just opened
	defer ts.Close()

	// get the URL with a GET request
	resp, err := http.Get(ts.URL)
	// If there is an error
	if err != nil {
		// log error and safely exit
		t.Fatal(err)
	}
	// if response status coe is not allowed
	if resp.StatusCode != http.StatusMethodNotAllowed {
		// print expected vs actual status codes
		t.Fatalf("expected status %d; actual status %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}

	// Create a buffer of bytes to store the json data in when encoding
	buf := new(bytes.Buffer)
	// establish a JSON hash map
	u := User{First: "Anton", Last: "Hibl"}

	// Encode using a JSON encoder
	err = json.NewEncoder(buf).Encode(&u)
	// test if there is an error
	if err != nil {
		// log error and safely exit
		t.Fatal(err)
	}

	// post the response
	resp, err = http.Post(ts.URL, "application/json", buf)
	// test if there is an error
	if err != nil {
		// call to log error and exit safely
		t.Fatal(err)
	}

	// if the response status code is not accepted
	if resp.StatusCode != http.StatusAccepted {
		// print expected vs actual status codes
		t.Fatalf("expected status %d; actual status %d", http.StatusAccepted, resp.StatusCode)
	}
	// Close the response body
	_ = resp.Body.Close()
}
