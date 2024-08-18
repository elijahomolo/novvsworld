package app

import (
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type Member struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Username  string
	Password  string
	Instagram string
	Twitter   string
	Threads   string
	Country   string
}

func (m *Member) VerifyUser(credentials Member) error {

	db := DB{}
	err := db.init()
	if err != nil {
		return fmt.Errorf("Failed to initialize database: %v", err)
	}

	// check if the credentials are correct
	m.Username = credentials.Username

	err = m.getUserID()
	if err != nil {
		return fmt.Errorf("Failed to get user id: %v", err)
	}

	selDB, err := db.Database.Query("SELECT token FROM auth WHERE id = ?", m.ID)
	if err != nil {
		return fmt.Errorf("Failed to query auth table: %v", err)
	}

	var token string
	for selDB.Next() {
		err = selDB.Scan(&token)
		if err != nil {
			return fmt.Errorf("Failed to scan auth table: %v", err)
		}
	}

	sDnc, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return fmt.Errorf("Failed to decode token: %v", err)
	}

	auth := string(sDnc)
	pw := strings.Trim(auth, fmt.Sprintf(credentials.Username+":"))

	// compare the password
	err = bcrypt.CompareHashAndPassword([]byte(pw), []byte(credentials.Password))
	if err != nil {
		return fmt.Errorf("Failed to compare password: %v", err)
	}

	return nil
}

func (m *Member) getUserID() error {
	db := DB{}
	err := db.init()
	if err != nil {
		return fmt.Errorf("Failed to initialize database: %v", err)
	}

	selDB, err := db.Database.Query("SELECT id FROM members WHERE username = ?", m.Username)
	if err != nil {
		return fmt.Errorf("Failed to query members table: %v", err)
	}

	var id int
	for selDB.Next() {
		err = selDB.Scan(&id)
		if err != nil {
			return fmt.Errorf("Failed to scan members table: %v", err)
		}
	}

	m.ID = id
	return nil
}

func (m *Member) Create() error {

	db := DB{}
	err := db.init()
	if err != nil {
		return fmt.Errorf("Failed to initialize database: %v", err)
	}

	// create a new member
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(m.Password), 8)
	if err != nil {
		return fmt.Errorf("Failed to generate hashed password: %v", err)
	}

	authToken := fmt.Sprintf("%s:%s", m.Username, hashedPassword)

	sEnc := base64.StdEncoding.EncodeToString([]byte(authToken))

	authForm, err := db.prepare(`INSERT INTO auth(auth_type, token) VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("Failed to prepare auth insert statement: %v", err)
	}

	_, err = authForm.Exec(1, sEnc)
	if err != nil {
		return fmt.Errorf("Failed to execute auth insert statement: %v", err)
	}

	selDB, err := db.Database.Query("SELECT id FROM auth WHERE token = ?", sEnc)
	if err != nil {
		return fmt.Errorf("Failed to query auth table: %v", err)
	}

	var id int
	for selDB.Next() {
		err = selDB.Scan(&id)
		if err != nil {
			return fmt.Errorf("Failed to scan auth table: %v", err)
		}
	}

	userForm, err := db.prepare(`INSERT INTO members(id, first_name, last_name, email, instagram_link, twitter_link, x_link, country, username) VALUES (?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		return fmt.Errorf("Failed to prepare member insert statement: %v", err)
	}

	_, err = userForm.Exec(id, m.FirstName, m.LastName, m.Email, m.Instagram, m.Twitter, m.Threads, m.Country, m.Username)
	if err != nil {
		return fmt.Errorf("Failed to execute member insert statement: %v", err)
	}

	return nil
}

func (m *Member) CreateUser() error {

	db := DB{}
	err := db.init()
	if err != nil {
		return fmt.Errorf("Failed to initialize database: %v", err)
	}

	// create a new member
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(m.Password), 8)
	if err != nil {
		return fmt.Errorf("Failed to generate hashed password: %v", err)
	}

	authToken := fmt.Sprintf("%s:%s", m.Username, hashedPassword)

	sEnc := base64.StdEncoding.EncodeToString([]byte(authToken))

	authForm, err := db.prepare(`INSERT INTO auth(auth_type, token) VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("Failed to prepare auth insert statement: %v", err)
	}

	_, err = authForm.Exec(1, sEnc)
	if err != nil {
		return fmt.Errorf("Failed to execute auth insert statement: %v", err)
	}

	selDB, err := db.Database.Query("SELECT id FROM auth WHERE token = ?", sEnc)
	if err != nil {
		return fmt.Errorf("Failed to query auth table: %v", err)
	}

	var id int
	for selDB.Next() {
		err = selDB.Scan(&id)
		if err != nil {
			return fmt.Errorf("Failed to scan auth table: %v", err)
		}
	}

	userForm, err := db.prepare(`INSERT INTO members(id, first_name, last_name, email, instagram_link, twitter_link, x_link, country, username) VALUES (?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		return fmt.Errorf("Failed to prepare member insert statement: %v", err)
	}

	_, err = userForm.Exec(id, m.FirstName, m.LastName, m.Email, m.Instagram, m.Twitter, m.Threads, m.Country, m.Username)
	if err != nil {
		return fmt.Errorf("Failed to execute member insert statement: %v", err)
	}

	return nil

}
