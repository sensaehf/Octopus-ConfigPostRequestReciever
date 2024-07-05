package main

import (
	"encoding/json"
	"log"
	"net/http"
	SecretDecoder "octopus/configReciever/src/decoder"
)

type Payload struct {
	ScanFileName    string `json:ScanFileName""`
	ScanDescription string `json:ScanDescription",omitempty"`
	Address         string `json:Address""`
	Username        string `json:Username""`
	Password        string `json:Password""`
	DomainScan      bool   `json:DomainScan""`
	Customer        string `json:Customer""`
}

func callOctopus(p Payload) {} // TODO call OctopusCLI to create Configuration File using the inputs

func verifieInputs(p Payload) bool { return true } //TODO validate Inputs

func recieveInfo(w http.ResponseWriter, req *http.Request) {
	var p Payload
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&p)

	if err != nil {
		log.Println("Failed to Parse JSON")
		return
	}

	Password, _ := SecretDecoder.DecodeSecret(p.Password)
	p.Password = Password
	if verifieInputs(p) {
		callOctopus(p)
	} else {
		log.Println("Failed to Validate Inputs")
	}

}

func main() {
	http.HandleFunc("/", recieveInfo)
	http.ListenAndServe(":8090", nil)
}
