package main

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func verifyUsernamePattern(username string) bool {
	// TODO
}

func verifyUserPasswordPattern(password string) bool {
	// TODO
}

var errUsernameAlreadyUsed = errors.New("username already used")

func (server *serverStruct) createUser(username string, password string) (userStruct, error) {
	id := generateUserId()
	passwordHash, passwordSalt := server.hashPassword(password)

	err := sqlitex.Execute(server.conn, "INSERT INTO user (id, username, password_hash, password_salt) VALUES (?, ?, ?, ?)", &sqlitex.ExecOptions{
		Args: []any{id, username, passwordHash, passwordSalt},
	})
	if sqlite.ErrCode(err) == sqlite.ResultConstraintUnique {
		return userStruct{}, errUsernameAlreadyUsed
	}
	if err != nil {
		return userStruct{}, fmt.Errorf("insert query failed: %s", err.Error())
	}

	user := userStruct{
		id:           id,
		username:     username,
		passwordHash: passwordHash,
		passwordSalt: passwordSalt,
	}
	return user, nil
}

func generateUserId() string {
	b := make([]byte, 10)
	rand.Read(b)
	alphabet := "abcdefghijkmnpqrstuvwyxz23456789"
	id := base32.NewEncoding(alphabet).EncodeToString(b)
	return id
}

func (server *serverStruct) hashPassword(password string) ([]byte, []byte) {
	// TODO
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
