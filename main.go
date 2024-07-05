package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	SecretDecoder "octopus/configReciever/src/decoder"
	"os/exec"
	"path/filepath"
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

func callOctopus(p Payload) {
	var dScan string
	if p.DomainScan {
		dScan = "Yes"
	} else {
		dScan = "No"
	}
	path := filepath.Join("C:", "inetpub", "oc_configurator", "configs", "OctopusConfigurator.exe")

	cmd := exec.Command(path,
		fmt.Sprintf("-scanFileName %s", p.ScanFileName),
		fmt.Sprintf("-scanDescription %s", p.ScanDescription),
		fmt.Sprintf("-ConfPassword %s", "TODO ENV VARIABLE"),
		fmt.Sprintf("-Address %s", p.Address),
		fmt.Sprintf("-Username %s", p.Username),
		fmt.Sprintf("-Password %s", p.Password),
		fmt.Sprintf("-Domainscan %s", dScan),
		fmt.Sprintf("-Customer %s", p.Customer),
	)
	fmt.Println(cmd.Args)
	/*rr := cmd.Run()
	if err != nil {
		log.Println(err)
	}*/

}

func verifyInputs(p Payload) bool { //TODO verify input
	match := true
	return match

}

// Recieves Data, validates and Adds it to octopus
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
	if verifyInputs(p) {
		callOctopus(p)
	} else {
		log.Println("Failed to Validate Inputs")
	}

}

func main() {
	http.HandleFunc("/", recieveInfo)
	http.ListenAndServe(":8090", nil)
}
