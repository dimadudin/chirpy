package database

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"sync"
	"time"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Users         map[int]User            `json:"users"`
	Chirps        map[int]Chirp           `json:"chirps"`
	RefreshTokens map[string]RefreshToken `json:"revocations"`
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Chirp struct {
	Id       int    `json:"id"`
	AuthorId int    `json:"author_id"`
	Body     string `json:"body"`
}

type RefreshToken struct {
	Id        string    `json:"id"`
	RevokedAt time.Time `json:"revoked_at"`
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
		newDBStructure := DBStructure{
			Users:         make(map[int]User),
			Chirps:        make(map[int]Chirp),
			RefreshTokens: make(map[string]RefreshToken),
		}
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

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email string, password string) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	newUser := User{
		Id:       len(dbs.Users) + 1,
		Email:    email,
		Password: password,
	}
	for _, v := range dbs.Users {
		if newUser.Email == v.Email {
			return User{}, errors.New("a user with this email already exists")
		}
	}
	dbs.Users[newUser.Id] = newUser
	err = db.writeDB(dbs)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}

// GetUsers returns all users in the database
func (db *DB) GetUsers() ([]User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	users := make([]User, 0, len(dbs.Chirps))
	for _, v := range dbs.Users {
		users = append(users, v)
	}
	sort.Slice(users, func(i, j int) bool { return users[i].Id < users[j].Id })
	return users, nil
}

// GetUserByEmail returns a user with the specified email
func (db *DB) GetUserByEmail(email string) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	for _, user := range dbs.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, errors.New("no user with such email")
}

// UpdateUser updates the user with the specified id to contain the new email and password
// returns the updated user
func (db *DB) UpdateUser(id int, email string, password string) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	updatedUser := User{
		Id:       id,
		Email:    email,
		Password: password,
	}
	dbs.Users[id] = updatedUser
	err = db.writeDB(dbs)
	if err != nil {
		return User{}, err
	}
	return updatedUser, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, author_id int) (Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	newChirp := Chirp{
		Id:       len(dbs.Chirps) + 1,
		AuthorId: author_id,
		Body:     body,
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

// CreateToken creates a new refresh token  and saves it to disk
func (db *DB) CreateToken(tokenStr string) (RefreshToken, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}
	newToken := RefreshToken{Id: tokenStr, RevokedAt: time.Time{}.UTC()}
	dbs.RefreshTokens[newToken.Id] = newToken
	err = db.writeDB(dbs)
	if err != nil {
		return RefreshToken{}, err
	}
	return newToken, nil
}

// GetToken returns a refresh token with the specified id
func (db *DB) GetToken(tokenStr string) (RefreshToken, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}
	token, ok := dbs.RefreshTokens[tokenStr]
	if !ok {
		return RefreshToken{}, errors.New("no such token")
	}
	return token, nil
}

// RevokeToken sets the revoked at time
func (db *DB) RevokeToken(tokenStr string) (RefreshToken, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}
	token, ok := dbs.RefreshTokens[tokenStr]
	if !ok {
		return RefreshToken{}, errors.New("no such token")
	}
	token.RevokedAt = time.Now().UTC()
	dbs.RefreshTokens[tokenStr] = token
	err = db.writeDB(dbs)
	if err != nil {
		return RefreshToken{}, err
	}
	return token, nil
}
