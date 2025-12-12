package main

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"time"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func (server *serverStruct) createSession(userId string) (sessionStruct, []byte, error) {
	// TODO

	// failed to add session to storage
}

func (server *serverStruct) addSessionToStorage(session sessionStruct) error {
	err := sqlitex.Execute(server.conn, "INSERT INTO session (id, secret_hash, user_id, created_at) VALUES (?, ?, ?, ?)", &sqlitex.ExecOptions{
		Args: []any{session.id, session.secretHash, session.userId, session.createdAt.Unix()},
	})
	if sqlite.ErrCode(err) == sqlite.ResultConstraintUnique {
		return errUsernameAlreadyUsed
	}
	if err != nil {
		return fmt.Errorf("insert query failed: %s", err.Error())
	}

	return nil
}

var errInvalidSessionToken = errors.New("invalid session token")

func (server *serverStruct) validateSessionToken(sessionToken string) (sessionStruct, error) {
	// TODO

	// failed to get session from storage
	// failed to delete session from storage
}

var errSessionNotFound = errors.New("session not found")

func (server *serverStruct) getSessionFromStorage(sessionId string) (sessionStruct, error) {
	sessions := []sessionStruct{}
	err := sqlitex.Execute(server.conn, "SELECT secret_hash, user_id, created_at FROM session WHERE id = ?", &sqlitex.ExecOptions{
		Args: []any{sessionId},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			secretHash := make([]byte, stmt.ColumnLen(0))
			stmt.ColumnBytes(0, secretHash)
			userId := stmt.ColumnText(1)
			createdAtUnix := stmt.ColumnInt64(2)
			session := sessionStruct{
				id:         sessionId,
				secretHash: secretHash,
				userId:     userId,
				createdAt:  time.Unix(createdAtUnix, 0),
			}
			sessions = append(sessions, session)
			return nil
		},
	})
	if err != nil {
		return sessionStruct{}, fmt.Errorf("select query failed: %s", err.Error())
	}
	if len(sessions) != 1 {
		return sessionStruct{}, errSessionNotFound
	}
	return sessions[0], nil
}

func (server *serverStruct) deleteSessionFromStorage(sessionId string) error {
	err := sqlitex.Execute(server.conn, "DELETE FROM session WHERE id = ?", &sqlitex.ExecOptions{
		Args: []any{sessionId},
	})
	if err != nil {
		return fmt.Errorf("delete query failed: %s", err.Error())
	}
	return nil
}

func generateSessionId() string {
	b := make([]byte, 10)
	rand.Read(b)
	alphabet := "abcdefghijkmnpqrstuvwyxz23456789"
	id := base32.NewEncoding(alphabet).EncodeToString(b)
	return id
}

type sessionStruct struct {
	id         string
	secretHash []byte
	userId     string
	createdAt  time.Time
}

func createSessionToken(sessionId string, sessionSecret []byte) string {
	// TODO
}
