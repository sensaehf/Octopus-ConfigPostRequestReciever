package main

import (
	"encoding/json"
	"log"
	"net/http"
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

func recieveInfo(w http.ResponseWriter, req *http.Request) {
	var p Payload
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&p)
	log.Println(p)

}

func main() {
	http.HandleFunc("/", recieveInfo)
	http.ListenAndServe(":8090", nil)
}
