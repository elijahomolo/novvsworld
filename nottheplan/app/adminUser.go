package app

import (
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type AdminUser struct {
	ID       int
	Username string
	Password string
}

func (a *AdminUser) Create() error {
	db := DB{}
	err := db.init()
	if err != nil {
		return fmt.Errorf("Failed to initialize database: %v", err)
	}

	// create a new member
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), 8)
	if err != nil {
		return fmt.Errorf("Failed to generate hashed password: %v", err)
	}

	authToken := fmt.Sprintf("%s:%s", a.Username, hashedPassword)

	sEnc := base64.StdEncoding.EncodeToString([]byte(authToken))

	authForm, err := db.prepare(`INSERT INTO auth(auth_type, token) VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("Failed to prepare auth insert statement: %v", err)
	}

	_, err = authForm.Exec(0, sEnc)
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

	userForm, err := db.prepare(`INSERT INTO admin_users(id, username) VALUES (?,?)`)
	if err != nil {
		return fmt.Errorf("Failed to prepare member insert statement: %v", err)
	}

	_, err = userForm.Exec(id, a.Username)
	if err != nil {
		return fmt.Errorf("Failed to execute member insert statement: %v", err)
	}

	return nil
}

func (a *AdminUser) Verify() error {

	db := DB{}
	// check if the credentials are correct

	err := db.init()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	selDB, err := db.Database.Query("SELECT id FROM admin_users WHERE username = ?", a.Username)
	if err != nil {
		return fmt.Errorf("failed to query admin_users table: %v", err)
	}

	var id int
	for selDB.Next() {
		err = selDB.Scan(&id)
		if err != nil {
			return fmt.Errorf("failed to scan members table: %v", err)
		}
	}

	selDB, err = db.Database.Query("SELECT token FROM auth WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to query auth table: %v", err)
	}

	var token string
	for selDB.Next() {
		err = selDB.Scan(&token)
		if err != nil {
			return fmt.Errorf("failed to scan auth table: %v", err)
		}
	}

	sDnc, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return fmt.Errorf("failed to decode token: %v", err)
	}

	auth := string(sDnc)
	pw := auth[len(a.Username)+1:]

	// compare the password
	err = bcrypt.CompareHashAndPassword([]byte(pw), []byte(a.Password))
	if err != nil {
		return fmt.Errorf("failed to compare password: %v", err)
	}

	return nil
}
