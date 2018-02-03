// curl -X POST -H "Content-Type: application/octet-stream" --data-binary '@filename' http://127.0.0.1:5050/upload

package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/golang/glog"
)

const (
	sessionIDHTTPHeader = "X-Session-ID"
	defaultAddress      = "0.0.0.0"
	defaultPort         = 5050
)

func main() {

	address := flag.String("address", defaultAddress, "address to serve on")
	port := flag.Int("port", defaultPort, "port to serve on")
	flag.Parse()

	http.HandleFunc("/upload", uploadHandler)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *address, *port), nil)
	if err != nil {
		glog.Errorf("Error starting server: %s", err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	sessionID := r.Header.Get(sessionIDHTTPHeader)
	if sessionID == "" {
		glog.Errorf("Request did not contain session ID HTTP header: %s", sessionIDHTTPHeader)
		return
	}

	filename := fmt.Sprintf("session-%s", sessionID)
	//fmt.Printf("%s", sessionID)
	//filename := fmt.Sprintf("session-test")

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	defer file.Close()
	if err != nil {
		glog.Errorf("Failed to open file %s: %s", filename, err)
		return
	}
	n, err := io.Copy(file, r.Body)
	if err != nil {
		glog.Errorf("Failed to copy body to file: %s", filename)
		return
	}

	glog.Info("%d bytes copied to %s", n, filename)
}
