package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	baseUrl       = "localhost:8081"
	createPostfix = "/notes"
	getPostfix    = "/notes"
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

type SyncMap struct {
	elems map[int64]*Note
	m     sync.RWMutex
}

var notes = &SyncMap{
	elems: make(map[int64]*Note),
}

func parseNoteId(idStr string) (int64, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func createNodeHandler(w http.ResponseWriter, r *http.Request) {
	info := &NodeInfo{}
	if err := json.NewDecoder(r.Body).Decode(info); err != nil {
		http.Error(w, "Failed to decode data", http.StatusBadRequest)
		return
	}

	rand.Seed(time.Now().Unix())
	now := time.Now()

	note := &Note{
		ID:        rand.Int63(),
		Info:      *info,
		CreatedAt: now,
		UpdatedAt: now,
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, "Failed to encode note data", http.StatusInternalServerError)
		return
	}

	notes.m.Lock()
	defer notes.m.Unlock()

	notes.elems[note.ID] = note
}

func getNodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("pisya")
	noteID := chi.URLParam(r, "id")
	id, err := parseNoteId(noteID)
	if err != nil {
		http.Error(w, "Failed to parse id from query", http.StatusBadRequest)
		return
	}
	notes.m.RLock()
	defer notes.m.RUnlock()

	note, ok := notes.elems[id]
	if !ok {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, "Can't encode note", http.StatusInternalServerError)
		return
	}
}

func main() {
	r := chi.NewRouter()

	r.Post(createPostfix, createNodeHandler)
	r.Get(getPostfix, getNodeHandler)

	err := http.ListenAndServe(baseUrl, r)
	if err != nil {
		log.Fatal("[FATAL]: ", err.Error())
		return
	}
}
