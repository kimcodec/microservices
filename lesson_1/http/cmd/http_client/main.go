package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"log"
	"net/http"
	"time"
)

const (
	baseUrl       = "http://localhost:8081"
	createPostfix = "/notes"
	getPostfix    = "/notes/id=%d"
)

type NodeInfo struct {
	Title    string `json:"title"`
	Context  string `json:"context"`
	Author   string `json:"author"`
	IsPublic bool   `json:"is_public"`
}

type Note struct {
	ID        int64     `json:"id"`
	Info      NodeInfo  `json:"info"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func createNote() (Note, error) {
	note := NodeInfo{
		Title:    gofakeit.BeerName(),
		Context:  gofakeit.IPv4Address(),
		Author:   gofakeit.Name(),
		IsPublic: gofakeit.Bool(),
	}

	data, err := json.Marshal(note)
	if err != nil {
		return Note{}, err
	}

	resp, err := http.Post(baseUrl+createPostfix, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return Note{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return Note{}, errors.New(fmt.Sprintf("wrong status code: %d", resp.StatusCode))
	}

	var createdNote Note
	if err := json.NewDecoder(resp.Body).Decode(&createdNote); err != nil {
		return Note{}, err
	}
	return createdNote, nil
}

func getNote(id int64) (Note, error) {
	resp, err := http.Get(fmt.Sprintf(baseUrl+getPostfix, id))
	if err != nil {
		return Note{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Note{}, errors.New("note not found")
	}

	if resp.StatusCode != http.StatusOK {
		return Note{}, errors.New("failed to get note")
	}

	var gotNote Note
	if err := json.NewDecoder(resp.Body).Decode(&gotNote); err != nil {
		return Note{}, err
	}

	return gotNote, nil
}

func main() {
	note, err := createNote()
	if err != nil {
		log.Fatalf("Failed to create note: %s", err.Error())
	}
	log.Printf("Note created:\n %+v\n", note)

	note1, err := getNote(note.ID)
	if err != nil {
		log.Fatalf("Failed to get note: %s", err.Error())
	}

	log.Printf("Note got:\n %+v\n", note1)
}
