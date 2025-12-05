package main

import (
	"errors"

	"github.com/faroedev/go-json"
)

const (
	actionSignUp = "sign_up"
)

func (server *serverStruct) invokeAction(action string, argumentsJSONObject json.ObjectStruct) (string, error) {
	if action == actionSignUp {
		return server.invokeSignupAction(argumentsJSONObject)
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

	user, errorCode := server.signUpAction(username, password)
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
