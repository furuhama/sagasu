package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func Start() {
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/", notFoundHandler)

	log.Print("Starting server ...")

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func searchHandler(w http.ResponseWriter, req *http.Request) {
	loggingRequest(req)

	w.WriteHeader(http.StatusOK)

	queryRegexp := req.URL.Query().Get("regexp")
	// Here's sample query
	results, _ := Search([]string{queryRegexp, "option1", "option2"})
	formatted, _ := formatJson(results)
	io.WriteString(w, formatted)
}

func notFoundHandler(w http.ResponseWriter, req *http.Request) {
	loggingRequest(req)

	w.WriteHeader(http.StatusNotFound)

	message := map[string]string{"message": "not found"}
	formatted, _ := formatJson(message)
	io.WriteString(w, formatted)
}

func formatJson(v interface{}) (string, error) {
	encoded, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	formatted := fmt.Sprintf("%s\n", string(encoded))

	return formatted, nil
}

func loggingRequest(req *http.Request) {
	requestLog := struct {
		Host   string      `json:"host,string"`
		Header http.Header `json:"header,[]string"`
		Method string      `json:"method,string"`
		Path   string      `json:"path,string"`
		Query  url.Values  `json:"query,map[string][]string"`
	}{
		Host:   req.Host,
		Header: req.Header,
		Method: req.Method,
		Path:   req.URL.Path,
		Query:  req.URL.Query(),
	}

	encodedReqLog, _ := json.Marshal(requestLog)

	log.Printf("request: %s\n", string(encodedReqLog))
}
