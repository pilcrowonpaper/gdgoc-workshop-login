package main

import (
	"errors"
	"fmt"
	"log"
)

func (server *serverStruct) signUpAction(username string, password string) (userStruct, string, string) {
	const (
		errorCodeUnexpectedError     = "unexpected_error"
		errorCodeInvalidUsername     = "invalid_username"
		errorCodeInvalidPassword     = "invalid_password"
		errorCodeUsernameAlreadyUsed = "username_already_used"
	)

	usernameValid := verifyUsernamePattern(username)
	if !usernameValid {
		return userStruct{}, "", errorCodeInvalidUsername
	}

	passwordValid := verifyUserPasswordPattern(password)
	if !passwordValid {
		return userStruct{}, "", errorCodeInvalidPassword
	}

	user, err := server.createUser(username, password)
	if errors.Is(err, errUsernameAlreadyUsed) {
		return userStruct{}, "", errorCodeUsernameAlreadyUsed
	}
	if err != nil {
		log.Printf("Failed to create user: %s", err.Error())
		return userStruct{}, "", errorCodeUnexpectedError
	}

	session, sessionSecret, err := server.createSession(user.id)
	if err != nil {
		log.Printf("Failed to create session: %s", err.Error())
		return userStruct{}, "", errorCodeUnexpectedError
	}

	sessionToken := createSessionToken(session.id, sessionSecret)

	return user, sessionToken, ""
}

func (server *serverStruct) signInAction(username string, password string) (userStruct, string, string) {
	const (
		errorCodeUnexpectedError   = "unexpected_error"
		errorCodeInvalidUsername   = "invalid_username"
		errorCodeInvalidPassword   = "invalid_password"
		errorCodeIncorrectPassword = "incorrect_password"
		errorCodeUserNotFound      = "user_not_found"
		errorCodeRateLimited       = "rate_limited"
	)

	usernameValid := verifyUsernamePattern(username)
	if !usernameValid {
		return userStruct{}, "", errorCodeInvalidUsername
	}

	passwordValid := verifyUserPasswordPattern(password)
	if !passwordValid {
		return userStruct{}, "", errorCodeInvalidPassword
	}

	user, err := server.getUserByUsernameFromStorage(username)
	if errors.Is(err, errUserNotFound) {
		return userStruct{}, "", errorCodeUserNotFound
	}
	if err != nil {
		log.Printf("Failed to get user by username from storage: %s", err.Error())
		return userStruct{}, "", errorCodeUnexpectedError
	}

	rateLimitAllowed := server.passwordVerificationRateLimit.check(user.id)
	if !rateLimitAllowed {
		return userStruct{}, "", errorCodeRateLimited
	}

	passwordCorrect := server.verifyUserPassword(password, user.passwordHash, user.passwordSalt)
	if !passwordCorrect {
		return userStruct{}, "", errorCodeIncorrectPassword
	}

	session, sessionSecret, err := server.createSession(user.id)
	if err != nil {
		log.Printf("Failed to create session: %s", err.Error())
		return userStruct{}, "", errorCodeUnexpectedError
	}

	sessionToken := createSessionToken(session.id, sessionSecret)

	return user, sessionToken, ""
}

func (server *serverStruct) getUserAction(sessionToken string) (userStruct, string) {
	const (
		errorCodeUnexpectedError     = "unexpected_error"
		errorCodeInvalidSessionToken = "invalid_session_token"
	)

	session, err := server.validateSessionToken(sessionToken)
	if errors.Is(err, errInvalidSessionToken) {
		return userStruct{}, errorCodeInvalidSessionToken
	}
	if err != nil {
		fmt.Printf("Failed to validate session token: %s", err.Error())
		return userStruct{}, errorCodeUnexpectedError
	}

	user, err := server.getUserFromStorage(session.userId)
	if err != nil {
		fmt.Printf("Failed to get user from storage: %s", err.Error())
		return userStruct{}, errorCodeUnexpectedError
	}

	return user, ""
}

func (server *serverStruct) signOutAction(sessionToken string) string {
	const (
		errorCodeUnexpectedError     = "unexpected_error"
		errorCodeInvalidSessionToken = "invalid_session_token"
	)

	session, err := server.validateSessionToken(sessionToken)
	if errors.Is(err, errInvalidSessionToken) {
		return errorCodeInvalidSessionToken
	}
	if err != nil {
		fmt.Printf("Failed to validate session token: %s", err.Error())
		return errorCodeUnexpectedError
	}

	err = server.deleteSessionFromStorage(session.id)
	if err != nil {
		fmt.Printf("Failed to delete session from storage: %s", err.Error())
		return errorCodeUnexpectedError
	}

	return ""
}
