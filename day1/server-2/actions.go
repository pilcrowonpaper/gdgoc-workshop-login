package main

import (
	"errors"
	"log"
)

func (server *serverStruct) signUpAction(username string, password string) (userStruct, string) {
	const (
		errorCodeUnexpectedError     = "unexpected_error"
		errorCodeInvalidUsername     = "invalid_username"
		errorCodeInvalidPassword     = "invalid_password"
		errorCodeUsernameAlreadyUsed = "username_already_used"
		errorCodeIncorrectPassword   = "incorrect_password"
	)

	usernameValid := verifyUsernamePattern(username)
	if !usernameValid {
		return userStruct{}, errorCodeInvalidUsername
	}

	passwordValid := verifyUserPasswordPattern(password)
	if !passwordValid {
		return userStruct{}, errorCodeInvalidPassword
	}

	user, err := server.createUser(username, password)
	if errors.Is(err, errUsernameAlreadyUsed) {
		return userStruct{}, errorCodeUsernameAlreadyUsed
	}
	if err != nil {
		log.Println(err.Error())
		return userStruct{}, errorCodeUnexpectedError
	}

	return user, ""
}

func (server *serverStruct) signInAction(username string, password string) (userStruct, string) {
	const (
		errorCodeUnexpectedError   = "unexpected_error"
		errorCodeInvalidUsername   = "invalid_username"
		errorCodeInvalidPassword   = "invalid_password"
		errorCodeIncorrectPassword = "incorrect_password"
		errorCodeUserNotFound      = "user_not_found"
		errorCodeRateLimited       = "rate_limited"
	)

	// TODO
}
