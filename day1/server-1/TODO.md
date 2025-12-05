# TODO

## verifyUserPasswordPattern()

`user.go`.

```go
func verifyUserPasswordPattern(password string) bool {
	return len(password) >= 10
}
```

## verifyUsernamePattern()

`user.go`.

```go
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
```

## \*serverStruct.hashPassword()

`user.go`.

```go
import (
	"crypto/rand"
)
```

```go
func (server *serverStruct) hashPassword(password string) ([]byte, []byte) {
	salt := make([]byte, 32)
	rand.Read(salt)
	hash := server.hashUserPasswordWithSalt([]byte(password), salt)
	return hash, salt
}
```

## \*serverStruct.signUpAction()

`actions.go`.

```go
import (
	"errors"
	"log"
)
```

```go
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
```
