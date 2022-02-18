package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
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

// Function to test a multi-part POST request to our server
func TestMultipartPost(t *testing.T) {
	// create a buffer to hold our request body
	reqBody := new(bytes.Buffer)
	// create a writer which will let us write in multiple parts at once in our request body
	w := multipart.NewWriter(reqBody)

	// JSON encoding in multiples
	for k, v := range map[string]string{
		// format the current time
		"date": time.Now().Format(time.RFC3339),
		// a description of our post request content
		"description": "Form values with attached files",
	} {
		// write our post request
		err := w.WriteField(k, v)
		// if error occurs
		if err != nil {
			// log and safely exit
			t.Fatal(err)
		}
	}

	// loop through byte string and assign two filepaths as values
	for i, file := range []string{
		"./files/hello.txt",
		"./files/goodbye.txt",
	} {
		// format filepath
		filePart, err := w.CreateFormFile(fmt.Sprintf("file%d", i+1),
			filepath.Base(file))
		// if error occurs
		if err != nil {
			// log and safely exit
			t.Fatal(err)
		}

		// open the filestream
		f, err := os.Open(file)
		// if error occurs
		if err != nil {
			// log and safely exit program
			t.Fatal(err)
		}

		// copy the filestream
		_, err = io.Copy(filePart, f)
		// close the filestream
		_ = f.Close()
		// check if there is an error
		if err != nil {
			// Log and exit safely
			t.Fatal(err)
		}
	}

	// close the post request
	err := w.Close()
	// check if error occured
	if err != nil {
		// log and exit safely if error occurs
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://httpbin.org/post", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d; actual status %d", http.StatusOK, resp.StatusCode)
	}

	t.Logf("\n%s", b)
}
