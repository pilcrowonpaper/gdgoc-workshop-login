package main

import (
	"errors"
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

	// TODO

	// Failed to create session
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

	// TODO

	// Failed to create session
}

func (server *serverStruct) getUserAction(sessionToken string) (userStruct, string) {
	const (
		errorCodeUnexpectedError     = "unexpected_error"
		errorCodeInvalidSessionToken = "invalid_session_token"
	)

	// TODO

	// Failed to validate session token
	// Failed to get user from storage
}

func (server *serverStruct) signOutAction(sessionToken string) string {
	const (
		errorCodeUnexpectedError     = "unexpected_error"
		errorCodeInvalidSessionToken = "invalid_session_token"
	)

	// TODO

	// Failed to validate session token
	// Failed to delete session from storage
}
