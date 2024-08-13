package db_server

import (
	"fmt"
	"time"

	"github.com/Anirudh-S-Kumar/disec/types"
	"golang.org/x/crypto/bcrypt"
)

func (s *ServerDB) AddUser(user_info *types.RegisterReq) error {
	hash_passwd, err := bcrypt.GenerateFromPassword([]byte(user_info.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("cannot generate password: %v", err)
	}

	stmt, err := s.db.Preparex(`	INSERT INTO 
									users (username, email, password_hash, date_created) 
								VALUES (?, ?, ?, ?)`)

	if err != nil {
		return fmt.Errorf("unable to create prepared statement: %v", err)
	}
	local_time := time.Now().String()

	_, err = stmt.Exec(user_info.Username, user_info.Email, hash_passwd, local_time)
	return err
}
