package app

import (
	"database/sql"
	"fmt"
)

type Submission struct {
	ID           int64
	UserID       int
	CommissionID int
	Demo         string
	SubmittedAt  []uint8
}

func (s *Submission) create() error {
	db := &DB{}

	err := db.init()
	if err != nil {
		return err
	}

	query := `INSERT INTO submissions(user_id, commission_id, demo_link) VALUES (?, ?, ?)`

	stmt, err := db.prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare submission insert statement: %v", err)
	}

	var res sql.Result

	res, err = db.execute(stmt, s.UserID, s.CommissionID, s.Demo)
	if err != nil {
		return fmt.Errorf("failed to execute submission insert statement: %v", err)
	}

	s.ID, err = res.LastInsertId()

	defer db.close()

	return nil
}

func (s *Submission) get() error {
	db := &DB{}

	err := db.init()
	if err != nil {
		return err
	}

	query := `SELECT * FROM submissions WHERE id = ?`

	stmt, err := db.prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare commission select statement: %v", err)
	}

	rows, err := stmt.Query(s.ID)
	if err != nil {
		return fmt.Errorf("failed to execute commission select statement: %v", err)
	}

	for rows.Next() {
		err = rows.Scan(&s.ID, &s.UserID, &s.CommissionID, &s.Demo, &s.SubmittedAt)
		if err != nil {
			return fmt.Errorf("failed to scan commission: %v", err)
		}
	}

	defer db.close()

	return nil
}

func (s *Submission) getAll() ([]Submission, error) {
	db := &DB{}

	err := db.init()
	if err != nil {
		return nil, err
	}

	query := `SELECT * FROM submissions`

	stmt, err := db.prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare submission select statement: %v", err)
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to execute submission select statement: %v", err)
	}

	submissions := []Submission{}

	for rows.Next() {
		submission := Submission{}
		err = rows.Scan(&submission.ID, &submission.UserID, &submission.CommissionID, &submission.Demo, &submission.SubmittedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan submission: %v", err)
		}

		submissions = append(submissions, submission)
	}

	defer db.close()

	return submissions, nil
}

func (s *Submission) update() error {
	db := &DB{}

	err := db.init()
	if err != nil {
		return err
	}

	query := `UPDATE submissions SET user_id = ?, commission_id = ?, demo_link = ? WHERE id = ?`

	stmt, err := db.prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare submission update statement: %v", err)
	}

	_, err = db.execute(stmt, s.UserID, s.CommissionID, s.Demo, s.ID)
	if err != nil {
		return fmt.Errorf("failed to execute submission update statement: %v", err)
	}

	defer db.close()

	return nil
}

func (s *Submission) delete() error {
	db := &DB{}

	err := db.init()
	if err != nil {
		return err
	}

	query := `DELETE FROM submissions WHERE id = ?`

	stmt, err := db.prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare submission delete statement: %v", err)
	}

	_, err = db.execute(stmt, s.ID)
	if err != nil {
		return fmt.Errorf("failed to execute submission delete statement: %v", err)
	}

	defer db.close()

	return nil
}

func (s *Submission) getByUserID(userID int) ([]Submission, error) {
	db := &DB{}

	err := db.init()
	if err != nil {
		return nil, err
	}

	query := `SELECT * FROM submissions WHERE user_id = ?`

	stmt, err := db.prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare submission select statement: %v", err)
	}

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute submission select statement: %v", err)
	}

	submissions := []Submission{}

	for rows.Next() {
		submission := Submission{}
		err = rows.Scan(&submission.ID, &submission.UserID, &submission.CommissionID, &submission.Demo, &submission.SubmittedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan submission: %v", err)
		}

		submissions = append(submissions, submission)
	}

	defer db.close()

	return submissions, nil
}

func (s *Submission) getCommission() (*Commission, error) {
	commission := &Commission{ID: s.CommissionID}

	err := commission.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get commission: %v", err)
	}

	return commission, nil
}

func (s *Submission) listCommissionsByUser(userID int) ([]Commission, error) {
	commissions := []Commission{}

	submissions, err := s.getByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get submissions by user: %v", err)
	}

	for _, submission := range submissions {
		commission := &Commission{ID: submission.CommissionID}
		err = commission.Get()
		if err != nil {
			return nil, fmt.Errorf("failed to get commission: %v", err)
		}

		commissions = append(commissions, *commission)
	}

	return commissions, nil
}

func (s *Submission) getByCommissionID(commissionID int) ([]Submission, error) {
	db := &DB{}

	err := db.init()
	if err != nil {
		return nil, err
	}

	query := `SELECT * FROM submissions WHERE commission_id = ?`

	stmt, err := db.prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare submission select statement: %v", err)
	}

	rows, err := stmt.Query(commissionID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute submission select statement: %v", err)
	}

	submissions := []Submission{}

	for rows.Next() {
		submission := Submission{}
		err = rows.Scan(&submission.ID, &submission.UserID, &submission.CommissionID, &submission.Demo, &submission.SubmittedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan submission: %v", err)
		}

		submissions = append(submissions, submission)
	}

	defer db.close()

	return submissions, nil
}

func (s *Submission) getByCommissionIDAndUserID(commissionID, userID int) ([]Submission, error) {
	db := &DB{}

	err := db.init()
	if err != nil {
		return nil, err
	}

	query := `SELECT * FROM submissions WHERE commission_id = ? AND user_id = ?`

	stmt, err := db.prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare submission select statement: %v", err)
	}

	rows, err := stmt.Query(commissionID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute submission select statement: %v", err)
	}

	submissions := []Submission{}

	for rows.Next() {
		submission := Submission{}
		err = rows.Scan(&submission.ID, &submission.UserID, &submission.CommissionID, &submission.Demo, &submission.SubmittedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan submission: %v", err)
		}

		submissions = append(submissions, submission)
	}

	defer db.close()

	return submissions, nil
}

func (s *Submission) checkDemoExists() ([]Submission, error) {
	db := &DB{}

	err := db.init()
	if err != nil {
		return nil, err
	}

	query := `SELECT * FROM submissions WHERE demo_link = ?`

	stmt, err := db.prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare submission select statement: %v", err)
	}

	rows, err := stmt.Query(s.Demo)
	if err != nil {
		return nil, fmt.Errorf("failed to execute submission select statement: %v", err)
	}

	submissions := []Submission{}

	for rows.Next() {
		submission := Submission{}
		err = rows.Scan(&submission.ID, &submission.UserID, &submission.CommissionID, &submission.Demo, &submission.SubmittedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan submission: %v", err)
		}

		submissions = append(submissions, submission)
	}

	defer db.close()

	return submissions, nil
}

func (s *Submission) deleteByUserID(userID int) error {
	db := &DB{}

	err := db.init()
	if err != nil {
		return err
	}

	query := `DELETE FROM submissions WHERE user_id = ?`

	stmt, err := db.prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare submission delete statement: %v", err)
	}

	_, err = db.execute(stmt, userID)
	if err != nil {
		return fmt.Errorf("failed to execute submission delete statement: %v", err)
	}

	defer db.close()

	return nil
}

func (s *Submission) deleteByCommissionID(commissionID int) error {
	db := &DB{}

	err := db.init()
	if err != nil {
		return err
	}

	query := `DELETE FROM submissions WHERE commission_id = ?`

	stmt, err := db.prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare submission delete statement: %v", err)
	}

	_, err = db.execute(stmt, commissionID)
	if err != nil {
		return fmt.Errorf("failed to execute submission delete statement: %v", err)
	}

	defer db.close()

	return nil
}
