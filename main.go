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

func callOctopus(p Payload) {}

func recieveInfo(w http.ResponseWriter, req *http.Request) {
	var p Payload
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&p)

	Password, _ := SecretDecoder.DecodeSecret(p.Password)
	p.Password = Password

	callOctopus(p)
	log.Println(p)
}

func main() {
	http.HandleFunc("/", recieveInfo)
	http.ListenAndServe(":8090", nil)
}
