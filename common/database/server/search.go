package db_server

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func (s *ServerDB) GetUserInfo(username string) (*sqlx.Row, error) {
	stmt, err := s.db.Preparex("SELECT * FROM users WHERE username = ?")
	if err != nil {
		return nil, fmt.Errorf("error in preparing statement: %v", err)
	}
	defer stmt.Close()

	return stmt.QueryRowx(username), nil
}

func (s *ServerDB) GetUserPasswdHash(username string) (string, error) {
	stmt, err := s.db.Preparex("SELECT password_hash FROM users WHERE username = ?")
	if err != nil {
		return "", fmt.Errorf("error in preparing statement: %v", err)
	}
	defer stmt.Close()

	var hash string
	err = stmt.QueryRowx(username).Scan(&hash)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("error in getting password hash: %v", err)
	}
	return hash, nil
}

func (s *ServerDB) UserExists(username, email string) error {
	var username_exists, email_exists bool

	stmt_user, err_user := s.db.Preparex(`SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)`)
	stmt_email, err_email := s.db.Preparex(`SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)`)

	if err_user != nil {
		return fmt.Errorf("error in preparing user exists query: %v", err_user)
	}

	if err_email != nil {
		return fmt.Errorf("error in preparing email exists query: %v", err_email)
	}

	defer stmt_email.Close()
	defer stmt_user.Close()

	log.Printf("email: %v", email)

	err_user = stmt_user.QueryRowx(username).Scan(&username_exists)
	err_email = stmt_email.QueryRowx(email).Scan(&email_exists)

	if err_user != nil && err_user != sql.ErrNoRows || err_email != nil && err_email != sql.ErrNoRows {
		return fmt.Errorf("error in executing userExists queries")
	}

	if username_exists {
		return fmt.Errorf("username already exists")
	}

	if email_exists {
		return fmt.Errorf("email already exists")
	}

	return nil

}
