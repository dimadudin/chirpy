package database

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"sync"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := DB{path: path, mu: &sync.RWMutex{}}
	err := db.ensureDB()
	if err != nil {
		return nil, err
	}
	return &db, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.Stat(db.path)
	if errors.Is(err, os.ErrNotExist) {
		newDBStructure := DBStructure{Chirps: make(map[int]Chirp)}
		err := db.writeDB(newDBStructure)
		return err
	} else {
		return err
	}
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	db.mu.Lock()
	err = os.WriteFile(db.path, dat, 0666)
	db.mu.Unlock()
	return err
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	dat, err := os.ReadFile(db.path)
	db.mu.RUnlock()
	if err != nil {
		return DBStructure{}, err
	}
	dbs := DBStructure{}
	err = json.Unmarshal(dat, &dbs)
	if err != nil {
		return DBStructure{}, err
	}
	return dbs, err
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	newChirp := Chirp{
		Id:   len(dbs.Chirps) + 1,
		Body: body,
	}
	dbs.Chirps[newChirp.Id] = newChirp
	err = db.writeDB(dbs)
	if err != nil {
		return Chirp{}, err
	}
	return newChirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := make([]Chirp, 0, len(dbs.Chirps))
	for _, v := range dbs.Chirps {
		chirps = append(chirps, v)
	}
	sort.Slice(chirps, func(i, j int) bool { return chirps[i].Id < chirps[j].Id })
	return chirps, nil
}

// GetChirpByID returns a chirp with the specified id
func (db *DB) GetChirpByID(id int) (Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	chirp, ok := dbs.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("no chirp with such ID")
	}
	return chirp, nil
}
