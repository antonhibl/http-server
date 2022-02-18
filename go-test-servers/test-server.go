package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/awoodbeck/gnp/ch09/handlers"
)

func TestSimpleHTTPServer(t *testing.T) {
	// Defining basic Server properties
	srv := &http.Server{
		// The Server's Address
		Addr: "127.0.0.1:8081",
		// Some default handlers
		// http.TimeoutHandler() is a middleware function before requests are sent to http.DefaultHandler()
		Handler: http.TimeoutHandler(
			handlers.DefaultHandler(), 2*time.Minute, ""),
		// The length of time clients can remain idle between requests
		IdleTimeout: 5 * time.Minute,
		// How long the server should wait to read a request header
		ReadHeaderTimeout: time.Minute,
	}

	//
	/*
		srv := &http.Server{
			Addr: "127.0.0.1:8081",
			Handler: mux,
			IdleTimeout: 5 * time.Minute,
			ReadHeaderTimeout: time.Minute,
		}
	*/

	// Create a new listener bound to the server's address
	l, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		t.Fatal(err)
	}

	// Define a serve method which returns http.ErrServerClosed when exiting normally
	// mainly it serves requests from the listener
	go func() {
		err := srv.Serve(l)
		if err != http.ErrServerClosed {
			t.Error(err)
		}
	}()

	// This would allow me to add TLS support with the proper certificates
	/*
		go func() {
			err := srv.ServeTLS(l, "cert.pem", "key.pem")
			if err != http.ErrServerClosed {
				t.Error(err)
			}
		}()
	*/

	// Setting a struct for HTTP Request Test Cases
	testCases := []struct {
		method   string
		body     io.Reader
		code     int
		response string
	}{
		// GET Method
		{http.MethodGet, nil, http.StatusOK, "Hello, friend."},
		// POST Method
		{http.MethodPost, bytes.NewBufferString("<world>"),
			http.StatusOK,
			"Hello, &lt;world&gt;!"},
		// http Headers
		{http.MethodHead, nil, http.StatusMethodNotAllowed, ""},
	}

	// establishing a client
	client := new(http.Client)
	// formatting the filepath in the server
	path := fmt.Sprintf("http://%s/", srv.Addr)

	// iterate over testcases
	for i, c := range testCases {
		// establish the new request
		r, err := http.NewRequest(c.method, path, c.body)

		// If there is an error, report it
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		// tell client to perform some request
		resp, err := client.Do(r)

		// if there is an error, report it
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		// If an unexpected status is returned
		if resp.StatusCode != c.code {
			t.Errorf("%d: unexpected status code: %q", i, resp.Status)
		}

		// read the body contents if not empty
		b, err := ioutil.ReadAll(resp.Body)

		// If there is an error report it
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		// close the body
		_ = resp.Body.Close()

		// test if response is a string, if not report error
		if c.response != string(b) {
			t.Errorf("%d: expected %q; actual %q", i, c.response, b)
		}
	}

	// close the server and check for an error in doing so
	if err := srv.Close(); err != nil {
		t.Fatal(err)
	}
}
