package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type LogPayload struct{
	Name string `json:"name"`
	Data string `json:"data"`
}
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payLoad := jsonResponse{
		Error:   false,
		Message: "Hit the broker!",
	}

	_ = app.writeJSON(w, http.StatusOK, payLoad)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)

	case "log":
		app.logItem(w,requestPayload.Log)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}


func (app *Config) logItem(w http.ResponseWriter , entry LogPayload){

	jsonData , _ := json.MarshalIndent(entry,"","/t")

	logServiceUrl := "http://host.docker.internal/log"

	request ,err :=http.NewRequest("POST",logServiceUrl,bytes.NewBuffer(jsonData))

	if err != nil{
		app.errorJSON(w,err)
		return
	}
	request.Header.Set("Content-type","application/json")

	client := &http.Client{}

	response ,err := client.Do(request)

	if err != nil{
		app.errorJSON(w,err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK{
		app.errorJSON(w,err)
		return
	}

	var payload jsonResponse
	payload.Error =false
	payload.Message ="logged"
	app.writeJSON(w,http.StatusAccepted,payload)

}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// first convert  response into auth request json format
	jsonData, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// call the auth service
	request, err := http.NewRequest("POST", "http://host.docker.internal:8081/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// create the client and get the response
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// log the response body for debugging
	body, _ := ioutil.ReadAll(response.Body)
	log.Println("Response Body:", string(body))

	// ensure we get the proper status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	} else if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("error in calling the auth service"), http.StatusInternalServerError)
		return
	}
log.Println("test123")
	 // Decode the response
	 var jsonFromAuthService jsonResponse
	 err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&jsonFromAuthService)
	 if err != nil {
		 log.Println("Error decoding JSON:", err)
		 app.errorJSON(w, errors.New("failed to decode JSON from auth service"))
		 return
	 }
 
	 // Log the decoded response
	 log.Println("Decoded JSON:", jsonFromAuthService)
 
	 // Check if there was an error in the response
	 if jsonFromAuthService.Error {
		 app.errorJSON(w, errors.New(jsonFromAuthService.Message), http.StatusUnauthorized)
		 return
	 }
 
	 // Send success response to client
	 var payLoad jsonResponse
	 payLoad.Error = false
	 payLoad.Message = "Authenticated!"
	 payLoad.Data = jsonFromAuthService.Data
 
	 app.writeJSON(w, http.StatusAccepted, payLoad)
}
