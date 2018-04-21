package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Index is a mapping from words to document pointers
type Index struct {
	store map[string][]*Document
}

// NewIndex returns a new, empty document Index
func NewIndex() *Index {
	return &Index{}
}

// AddDocument indexes the tokens in a document
func (index *Index) AddDocument(doc Document) error {
	uniqueWords, err := doc.uniqueWords()
	if err != nil {
		return err
	}

	for _, word := range uniqueWords {

		arr, ok := index.store[word]
		if ok {
			arr = append(arr, &doc)
		} else {
			arr = []*Document{&doc}
		}
		index.store[word] = arr
	}

	return nil
}

// BuildFromPath adds all documents visitable from a path to an Index
func (index *Index) BuildFromPath(path string) error {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("Path %s does not exist", path)
	}

	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !(info.IsDir()) {
			doc := Document{path}
			index.AddDocument(doc)
		} else if strings.HasPrefix(info.Name(), ".") {
			// ignore hidden directories
			return filepath.SkipDir
		}
		return nil
	})
	return nil
}
