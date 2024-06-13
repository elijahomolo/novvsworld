package app

import (
	"fmt"
	"time"
)

type Commission struct {
	ID          int       `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Budget      float64   `json:"budget,omitempty"`
	Currency    string    `json:"currency,omitempty"`
	Location    string    `json:"location,omitempty"`
	Deadline    time.Time `json:"deadline,omitempty"`
	Status      string    `json:"status,omitempty"`
	Winner      string    `json:"winner,omitempty"`
	Entries     []uint8   `json:"entries,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

func (c *Commission) Create() error {

	db := &DB{}

	err := db.init()
	if err != nil {
		return err
	}

	query := `INSERT INTO commissions(name, description, budget, currency, location, deadline, status, winner, entries) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := db.prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare commission insert statement: %v", err)

	}

	err = db.execute(stmt, c.Name, c.Description, c.Budget, c.Currency, c.Location, c.Deadline, c.Status, c.Winner, fmt.Sprintf("%v", c.Entries))
	if err != nil {
		return fmt.Errorf("failed to execute commission insert statement: %v", err)
	}

	defer db.close()

	return nil
}

func (c *Commission) Update() error {
	db := &DB{}

	err := db.init()
	if err != nil {
		return err
	}

	query := `UPDATE commissions SET name = ?, description = ?, budget = ?, currency = ?, location = ?, deadline = ?, status = ?, winner = ?, entries = ? WHERE id = ?`

	stmt, err := db.prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare commission update statement: %v", err)
	}

	err = db.execute(stmt, c.Name, c.Description, c.Budget, c.Currency, c.Location, c.Deadline, c.Status, c.Winner, fmt.Sprintf("%v", c.Entries), c.ID)
	if err != nil {
		return fmt.Errorf("failed to execute commission update statement: %v", err)
	}

	defer db.close()

	return nil
}

func (c *Commission) Delete() error {
	db := &DB{}

	err := db.init()
	if err != nil {
		return err
	}

	query := `DELETE FROM commissions WHERE id = ?`

	stmt, err := db.prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare commission delete statement: %v", err)
	}

	err = db.execute(stmt, c.ID)
	if err != nil {
		return fmt.Errorf("failed to execute commission delete statement: %v", err)
	}

	defer db.close()

	return nil
}

func (c *Commission) Get() error {
	db := &DB{}

	err := db.init()
	if err != nil {
		return err
	}

	query := `SELECT * FROM commissions WHERE id = ?`

	stmt, err := db.prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare commission select statement: %v", err)
	}

	rows, err := stmt.Query(c.ID)
	if err != nil {
		return fmt.Errorf("failed to execute commission select statement: %v", err)
	}

	for rows.Next() {
		err = rows.Scan(&c.ID, &c.Name, &c.Description, &c.Budget, &c.Currency, &c.Location, &c.Deadline, &c.Status, &c.Winner, &c.Entries, &c.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to scan commission: %v", err)
		}
	}

	defer db.close()

	return nil
}

func (c *Commission) GetAll() ([]Commission, error) {
	db := &DB{}

	err := db.init()
	if err != nil {
		return nil, err
	}

	query := `SELECT * FROM commissions`

	stmt, err := db.prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare commission select statement: %v", err)
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to execute commission select statement: %v", err)
	}

	commissions := []Commission{}
	for rows.Next() {
		var commission Commission
		err = rows.Scan(&commission.ID, &commission.Name, &commission.Description, &commission.Budget, &commission.Currency, &commission.Location, &commission.Deadline, &commission.Status, &commission.Winner, &commission.Entries, &commission.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan commission: %v", err)
		}
		commissions = append(commissions, commission)
	}

	defer db.close()

	return commissions, nil
}
