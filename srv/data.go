package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Document represents a textual document
type Document struct {
	Path string `json:"path"`
}

// Match represents a matching line within a Document
type Match struct {
	Document           Document `json:"document"`
	MatchingLine       string   `json:"line"`
	MatchingLineNumber int      `json:"lineNumber"`
}

// Contents returns the contents of a Document as a byte array
func (doc *Document) Contents() ([]byte, error) {
	contents, err := ioutil.ReadFile(doc.Path)
	return contents, err
}

func (doc *Document) uniqueWords() ([]string, error) {
	contents, err := doc.Contents()
	if err != nil {
		return []string{}, err
	}

	wordMap := make(map[string]bool)
	for _, word := range bytes.Split(contents, []byte{' '}) {
		wordMap[string(word)] = true
	}

	words := make([]string, len(wordMap))
	cnt := 0
	for word := range wordMap {
		words[cnt] = word
		cnt++
	}

	return words, nil
}

// Search returns an array of matches for a query
func (doc *Document) Search(query string) ([]Match, error) {
	var results []Match

	f, err := os.Open(doc.Path)
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
				Document:           *doc,
				MatchingLine:       line,
				MatchingLineNumber: i,
			})
		}
	}
	return results, nil
}
