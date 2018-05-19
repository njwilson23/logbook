package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Entry struct {
	filename string
	topic    string
	date     time.Time
	exists   bool
}

func (e *Entry) AddHeader() error {
	if !e.exists {
		f, err := os.Create(e.filename)
		if err != nil {
			fmt.Printf("failed to prepare %s", e.filename)
			return err
		}
		defer f.Close()
		f.Write([]byte(fmt.Sprintf("# %s\n", e.date.Format("2006 Jan 02"))))
	}
	return nil
}

func getEditor() string {
	editor := os.Getenv("LBEDITOR")
	if editor == "" {
		editor = "vim"
	}
	return editor
}

func run(topic string, offsetDays int) {
	entry := getEntry(topic, offsetDays)

	err := entry.AddHeader()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	cmd := exec.Command(getEditor(), entry.filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}

func getEntry(topic string, offsetDays int) *Entry {
	notesDir := os.Getenv("LBDIR")
	if notesDir == "" {
		notesDir = filepath.Join(os.Getenv("HOME"), "notes")
	}

	date := time.Now().Add(time.Hour * time.Duration(24*offsetDays))
	var topicStr string
	if topic != "" {
		topicStr = fmt.Sprintf("%s-", topic)
	}
	filename := filepath.Join(notesDir, fmt.Sprintf("%s%s.md", topicStr, date.Format("2006-01-02")))

	exists := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exists = false
	}

	return &Entry{
		filename,
		topic,
		date,
		exists,
	}
}

func main() {
	days := flag.Int("days", 0, "number of days relative to current date to open entry for")
	topic := flag.String("topic", "", "entry topic")

	flag.Parse()

	run(*topic, *days)
}
