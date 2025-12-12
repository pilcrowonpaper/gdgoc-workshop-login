package main

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base32"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func verifyUsernamePattern(username string) bool {
	chars := []rune(username)
	if len(chars) < 3 || len(chars) > 16 {
		return false
	}
	for _, char := range chars {
		if char >= 'a' && char <= 'z' {
			continue
		}
		if char >= '0' && char <= '9' {
			continue
		}
		return false
	}
	return true
}

func verifyUserPasswordPattern(password string) bool {
	return len(password) >= 10 && len(password) <= 100
}

var errUsernameAlreadyUsed = errors.New("username already used")

func (server *serverStruct) createUser(username string, password string) (userStruct, error) {
	id := generateUserId()
	passwordHash, passwordSalt := server.hashPassword(password)

	user := userStruct{
		id:           id,
		username:     username,
		passwordHash: passwordHash,
		passwordSalt: passwordSalt,
	}

	err := server.addUserToStorage(user)
	if errors.Is(err, errUsernameAlreadyUsed) {
		return userStruct{}, errUsernameAlreadyUsed
	}
	if err != nil {
		return userStruct{}, fmt.Errorf("failed to add user to storage: %s", err.Error())
	}

	return user, nil
}

func (server *serverStruct) addUserToStorage(user userStruct) error {
	err := sqlitex.Execute(server.conn, "INSERT INTO user (id, username, password_hash, password_salt) VALUES (?, ?, ?, ?)", &sqlitex.ExecOptions{
		Args: []any{user.id, user.username, user.passwordHash, user.passwordSalt},
	})
	if sqlite.ErrCode(err) == sqlite.ResultConstraintUnique {
		return errUsernameAlreadyUsed
	}
	if err != nil {
		return fmt.Errorf("insert query failed: %s", err.Error())
	}
	return nil
}

var errUserNotFound = errors.New("user not found")

func (server *serverStruct) getUserByUsernameFromStorage(username string) (userStruct, error) {
	users := []userStruct{}
	err := sqlitex.Execute(server.conn, "SELECT id, password_hash, password_salt FROM user WHERE username = ?", &sqlitex.ExecOptions{
		Args: []any{username},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			id := stmt.ColumnText(0)
			passwordHash := make([]byte, stmt.ColumnLen(1))
			stmt.ColumnBytes(1, passwordHash)
			passwordSalt := make([]byte, stmt.ColumnLen(2))
			stmt.ColumnBytes(2, passwordSalt)
			user := userStruct{
				id:           id,
				username:     username,
				passwordHash: passwordHash,
				passwordSalt: passwordSalt,
			}
			users = append(users, user)
			return nil
		},
	})
	if err != nil {
		return userStruct{}, fmt.Errorf("select query failed: %s", err.Error())
	}
	if len(users) != 1 {
		return userStruct{}, errUserNotFound
	}
	return users[0], nil
}

func (server *serverStruct) getUserFromStorage(userId string) (userStruct, error) {
	users := []userStruct{}
	err := sqlitex.Execute(server.conn, "SELECT username, password_hash, password_salt FROM user WHERE id = ?", &sqlitex.ExecOptions{
		Args: []any{userId},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			username := stmt.ColumnText(0)
			passwordHash := make([]byte, stmt.ColumnLen(1))
			stmt.ColumnBytes(1, passwordHash)
			passwordSalt := make([]byte, stmt.ColumnLen(2))
			stmt.ColumnBytes(2, passwordSalt)
			user := userStruct{
				id:           userId,
				username:     username,
				passwordHash: passwordHash,
				passwordSalt: passwordSalt,
			}
			users = append(users, user)
			return nil
		},
	})
	if err != nil {
		return userStruct{}, fmt.Errorf("select query failed: %s", err.Error())
	}
	if len(users) != 1 {
		return userStruct{}, errUserNotFound
	}
	return users[0], nil
}

func generateUserId() string {
	b := make([]byte, 10)
	rand.Read(b)
	alphabet := "abcdefghijkmnpqrstuvwyxz23456789"
	id := base32.NewEncoding(alphabet).EncodeToString(b)
	return id
}

func (server *serverStruct) hashPassword(password string) ([]byte, []byte) {
	salt := make([]byte, 32)
	rand.Read(salt)
	hash := server.hashUserPasswordWithSalt([]byte(password), salt)
	return hash, salt
}

func (server *serverStruct) verifyUserPassword(password string, passwordHash []byte, passwordSalt []byte) bool {
	hashed := server.hashUserPasswordWithSalt([]byte(password), passwordSalt)
	return subtle.ConstantTimeCompare(passwordHash, hashed) == 1
}

func (server *serverStruct) hashUserPasswordWithSalt(password []byte, salt []byte) []byte {
	server.passwordHashingSemaphore.Acquire(context.Background(), 1)
	hash := argon2.IDKey(password, salt, 3, 64*1024, 1, 32)
	server.passwordHashingSemaphore.Release(1)
	return hash
}

type userStruct struct {
	id           string
	username     string
	passwordHash []byte
	passwordSalt []byte
}
