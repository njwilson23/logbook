package main

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type LogEntry struct {
	Path string `json:"path"`
}

type LogMatch struct {
	LogEntry           *LogEntry `json:"logEntry"`
	MatchingLine       string    `json:"line"`
	MatchingLineNumber int       `json:"lineNumber"`
}

func (L *LogEntry) Search(query string) ([]LogMatch, error) {
	var results []LogMatch

	f, err := os.Open(L.Path)
	if err != nil {
		return results, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	var line string
	done := false
	for i := 0; !done; i++ {
		line, err = reader.ReadString('\n')
		if err == io.EOF {
			done = true
		} else if err != nil {
			return results, err
		}
		if strings.Contains(line, query) {
			// TODO: fuzzy match
			results = append(results, LogMatch{
				LogEntry:           L,
				MatchingLine:       line,
				MatchingLineNumber: i,
			})
		}
	}
	return results, nil
}

func getLogEntrys(logdir string) ([]LogEntry, error) {
	files, err := ioutil.ReadDir(logdir)
	if err != nil {
		return nil, err
	}

	var filenames []LogEntry
	for _, file := range files {
		if (filepath.Ext(file.Name()) == ".md") && !file.IsDir() {
			filenames = append(filenames, LogEntry{filepath.Join(logdir, file.Name())})
		}
	}
	return filenames, nil
}

func getLogDir() (logdir string) {
	logdir = os.Getenv("LBDIR")
	if len(logdir) == 0 {
		logdir = os.ExpandEnv("$HOME/notes")
	}
	return
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	logdir := getLogDir()
	logs, err := getLogEntrys(logdir)
	if err != nil {
		w.WriteHeader(500)
	}

	data, err := json.Marshal(logs)
	if err != nil {
		w.WriteHeader(500)
	}
	w.Write(data)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	logdir := getLogDir()
	logs, err := getLogEntrys(logdir)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("ERROR: failed to get log entries (%s)\n", err)
		return
	}

	// Extract query
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(400)
		log.Printf("INFO: bad query parameters ('%s')\n", r.URL.RawQuery)
		return
	}

	toFind, ok := queryValues["find"]
	if !ok {
		w.WriteHeader(400)
		log.Println("INFO: bad query parameters missing 'find'")
		return
	}

	var matches []LogMatch
	for _, logEntry := range logs {
		ms, err := logEntry.Search(toFind[0])
		if err != nil {
			w.WriteHeader(500)
			log.Printf("ERROR: search error (%s)\n", err)
			return
		}
		if len(ms) != 0 {
			matches = append(matches, ms...)
		}
	}

	data, err := json.Marshal(matches)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("ERROR: failure marshaling JSON (%s)\n", err)
		return
	}
	w.Write(data)
}

func main() {

	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/search", searchHandler)
	log.Println("starting up...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
