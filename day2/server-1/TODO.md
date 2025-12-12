# TODO

## \*server.createSession()

```go
func (server *serverStruct) createSession(userId string) (sessionStruct, []byte, error) {
	id := generateSessionId()

	secret := make([]byte, 32)
	rand.Read(secret)
	secretHash := sha256.Sum256(secret)

	session := sessionStruct{
		id:         id,
		secretHash: secretHash[:],
		userId:     userId,
		createdAt:  time.Now(),
	}

	err := server.addSessionToStorage(session)
	if err != nil {
		return sessionStruct{}, nil, fmt.Errorf("failed to add session to storage: %s", err.Error())
	}

	return session, secret, nil
}
```

## createSessionToken()

```go
func createSessionToken(sessionId string, sessionSecret []byte) string {
	encodedSessionSecret := base64.StdEncoding.EncodeToString(sessionSecret)
	sessionToken := fmt.Sprintf("%s.%s", sessionId, encodedSessionSecret)
	return sessionToken
}
```

## \*server.signUpAction()


```go
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
		log.Println("Failed to create session: %s", err.Error())
		return userStruct{}, "", errorCodeUnexpectedError
	}

	sessionToken := createSessionToken(session.id, sessionSecret)

	return user, sessionToken, ""
}
```

## \*server.signInAction()

```go
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
```

## \*server.validateSessionToken()

```go
func (server *serverStruct) validateSessionToken(sessionToken string) (sessionStruct, error) {
	sessionTokenParts := strings.Split(sessionToken, ".")
	if len(sessionTokenParts) != 2 {
		return sessionStruct{}, errInvalidSessionToken
	}

	sessionId := sessionTokenParts[0]
	sessionSecret, err := base64.StdEncoding.DecodeString(sessionTokenParts[1])
	if err != nil {
		return sessionStruct{}, errInvalidSessionToken
	}

	session, err := server.getSessionFromStorage(sessionId)
	if errors.Is(err, errSessionNotFound) {
		return sessionStruct{}, errInvalidSessionToken
	}
	if err != nil {
		return sessionStruct{}, fmt.Errorf("failed to get session: %s", err.Error())
	}

	if time.Since(session.createdAt) >= 24*time.Hour {
		err = server.deleteSessionFromStorage(sessionId)
		if err != nil {
			return sessionStruct{}, fmt.Errorf("failed to delete session from storage: %s", err.Error())
		}
		return sessionStruct{}, errSessionNotFound
	}

	sessionSecretHash := sha256.Sum256(sessionSecret)
	sessionSecretCorrect := subtle.ConstantTimeCompare(session.secretHash, sessionSecretHash[:]) == 1
	if !sessionSecretCorrect {
		return sessionStruct{}, errInvalidSessionToken
	}

	return session, nil
}
```

## \*server.getUserAction()

## \*server.createSessionToken()
