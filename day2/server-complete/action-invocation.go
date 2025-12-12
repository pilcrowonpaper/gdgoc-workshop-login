package main

import (
	"errors"

	"github.com/faroedev/go-json"
)

const (
	actionSignUp  = "sign_up"
	actionSignIn  = "sign_in"
	actionGetUser = "get_user"
	actionSignOut = "sign_out"
)

func (server *serverStruct) invokeAction(action string, argumentsJSONObject json.ObjectStruct) (string, error) {
	if action == actionSignUp {
		return server.invokeSignupAction(argumentsJSONObject)
	}
	if action == actionSignIn {
		return server.invokeSignInAction(argumentsJSONObject)
	}
	if action == actionGetUser {
		return server.invokeGetUserAction(argumentsJSONObject)
	}
	if action == actionSignOut {
		return server.invokeSignOutAction(argumentsJSONObject)
	}
	return "", errors.New("unknown action")
}

func (server *serverStruct) invokeSignupAction(argumentsJSONObject json.ObjectStruct) (string, error) {
	username, err := argumentsJSONObject.GetString("username")
	if err != nil {
		return "", errors.New("invalid or missing 'username' argument")
	}
	password, err := argumentsJSONObject.GetString("password")
	if err != nil {
		return "", errors.New("invalid or missing 'password' argument")
	}

	user, sessionToken, errorCode := server.signUpAction(username, password)
	if errorCode != "" {
		resultJSON := createActionErrorResultJSON(errorCode)
		return resultJSON, nil
	}

	userJSON := createUserJSON(user)

	valuesJSONBuilder := json.NewObjectBuilder()
	valuesJSONBuilder.AddJSON("user", userJSON)
	valuesJSONBuilder.AddString("session_token", sessionToken)
	valuesJSON := valuesJSONBuilder.Done()

	resultJSON := createActionSuccessResultJSON(valuesJSON)
	return resultJSON, nil
}

func (server *serverStruct) invokeSignInAction(argumentsJSONObject json.ObjectStruct) (string, error) {
	username, err := argumentsJSONObject.GetString("username")
	if err != nil {
		return "", errors.New("invalid or missing 'username' argument")
	}
	password, err := argumentsJSONObject.GetString("password")
	if err != nil {
		return "", errors.New("invalid or missing 'password' argument")
	}

	user, sessionToken, errorCode := server.signInAction(username, password)
	if errorCode != "" {
		resultJSON := createActionErrorResultJSON(errorCode)
		return resultJSON, nil
	}

	userJSON := createUserJSON(user)

	valuesJSONBuilder := json.NewObjectBuilder()
	valuesJSONBuilder.AddJSON("user", userJSON)
	valuesJSONBuilder.AddString("session_token", sessionToken)
	valuesJSON := valuesJSONBuilder.Done()

	resultJSON := createActionSuccessResultJSON(valuesJSON)
	return resultJSON, nil
}

func (server *serverStruct) invokeGetUserAction(argumentsJSONObject json.ObjectStruct) (string, error) {
	sessionToken, err := argumentsJSONObject.GetString("session_token")
	if err != nil {
		return "", errors.New("invalid or missing 'session_token' argument")
	}

	user, errorCode := server.getUserAction(sessionToken)
	if errorCode != "" {
		resultJSON := createActionErrorResultJSON(errorCode)
		return resultJSON, nil
	}

	userJSON := createUserJSON(user)

	valuesJSONBuilder := json.NewObjectBuilder()
	valuesJSONBuilder.AddJSON("user", userJSON)
	valuesJSON := valuesJSONBuilder.Done()

	resultJSON := createActionSuccessResultJSON(valuesJSON)
	return resultJSON, nil
}

func (server *serverStruct) invokeSignOutAction(argumentsJSONObject json.ObjectStruct) (string, error) {
	sessionToken, err := argumentsJSONObject.GetString("session_token")
	if err != nil {
		return "", errors.New("invalid or missing 'session_token' argument")
	}

	errorCode := server.signOutAction(sessionToken)
	if errorCode != "" {
		resultJSON := createActionErrorResultJSON(errorCode)
		return resultJSON, nil
	}

	valuesJSONBuilder := json.NewObjectBuilder()
	valuesJSON := valuesJSONBuilder.Done()

	resultJSON := createActionSuccessResultJSON(valuesJSON)
	return resultJSON, nil
}

func createUserJSON(user userStruct) string {
	jsonBuilder := json.NewObjectBuilder()
	jsonBuilder.AddString("id", user.id)
	jsonBuilder.AddString("username", user.username)
	userJSON := jsonBuilder.Done()
	return userJSON
}

func createActionSuccessResultJSON(valuesJSON string) string {
	objectBuilder := json.NewObjectBuilder()
	objectBuilder.AddBool("ok", true)
	objectBuilder.AddJSON("values", valuesJSON)
	resultJSON := objectBuilder.Done()
	return resultJSON
}

func createActionErrorResultJSON(errorCode string) string {
	jsonBuilder := json.NewObjectBuilder()
	jsonBuilder.AddBool("ok", false)
	jsonBuilder.AddString("error_code", errorCode)
	resultJSON := jsonBuilder.Done()
	return resultJSON
}
