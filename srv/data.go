package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Match struct {
	Document           Document `json:"logEntry"`
	MatchingLine       string   `json:"line"`
	MatchingLineNumber int      `json:"lineNumber"`
}

type Document struct {
	Path string `json:"path"`
}

func (L *Document) Search(query string) ([]Match, error) {
	var results []Match

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
			results = append(results, Match{
				Document:           *L,
				MatchingLine:       line,
				MatchingLineNumber: i,
			})
		}
	}
	return results, nil
}
