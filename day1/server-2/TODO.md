# TODO

## \*serverStruct.verifyUserPassword()

`user.go`.

```go
import (
	"crypto/subtle"
)
```

```go
func (server *serverStruct) verifyUserPassword(password string, passwordHash []byte, passwordSalt []byte) bool {
	hashed := server.hashUserPasswordWithSalt([]byte(password), passwordSalt)
	return subtle.ConstantTimeCompare(passwordHash, hashed) == 1
}
```

## \*serverStruct.signInAction()

`actions.go`.

```go
import (
	"errors"
	"log"
)
```

```go
func (server *serverStruct) signInAction(username string, password string) (userStruct, string) {
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
		return userStruct{}, errorCodeInvalidUsername
	}

	passwordValid := verifyUserPasswordPattern(password)
	if !passwordValid {
		return userStruct{}, errorCodeInvalidPassword
	}

	user, err := server.getUserByUsername(username)
	if errors.Is(err, errUserNotFound) {
		return userStruct{}, errorCodeUserNotFound
	}
	if err != nil {
		log.Println(err)
		return userStruct{}, errorCodeUnexpectedError
	}

	passwordCorrect := server.verifyUserPassword(password, user.passwordHash, user.passwordSalt)
	if !passwordCorrect {
		return userStruct{}, errorCodeIncorrectPassword
	}

	return user, ""
}

```
