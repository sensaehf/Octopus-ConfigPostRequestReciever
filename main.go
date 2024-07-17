//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	eventlogger "octopus/configReciever/src/EventLogger"
	SecretDecoder "octopus/configReciever/src/decoder"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joho/godotenv"
)

var DevFlag *bool // Feature flag that stop running the octopus program and Logging to the event viewer
var WindowsLog eventlogger.Logger

type Payload struct {
	ScanFileName    string `json:"scanFileName"`
	ScanDescription string `json:"ScanDescription"`
	Address         string `json:"Address"`
	Username        string `json:"Username"`
	Password        string `json:"Password"`
	DomainScan      bool   `json:"DomainScan"`
	Customer        string `json:"Customer"`
}

func (p Payload) String() string {
	return fmt.Sprintf("scanfileName:%s, ScanDescription:%s, Address:%s, Username:%s, Customer:%s", p.ScanFileName, p.ScanDescription, p.Address, p.Username, p.Customer)
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
		fmt.Sprintf("-ConfPassword %s", os.Getenv("OCTOPUS_KEY")),
		fmt.Sprintf("-Address %s", p.Address),
		fmt.Sprintf("-Username %s", p.Username),
		fmt.Sprintf("-Password %s", p.Password),
		fmt.Sprintf("-Domainscan %s", dScan),
		fmt.Sprintf("-Customer %s", p.Customer),
	)
	if !*DevFlag {
		err := cmd.Run()
		if err != nil {
			WindowsLog.Warning(
				fmt.Errorf("failed to run command using %s", p.String()))
		}

	}

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
		WindowsLog.Error(errors.New("failed to Parse JSON"))
		return
	}

	Password, _ := SecretDecoder.DecodeSecret(p.Password)
	p.Password = Password
	if verifyInputs(p) {
		callOctopus(p)
	} else {
		p.Password = "SECRET"
		WindowsLog.Warning(fmt.Errorf("input failed to parse for payload %s", p.String()))
	}

}

func main() {
	WindowsLog.Init()

	err := godotenv.Load("./.env")
	if err != nil {
		WindowsLog.Error(fmt.Errorf("failed to load .env file, %s", err.Error()))
		os.Exit(1)
	}

	DevFlag = flag.Bool("Dev", false, "Turns on Dev mode. Stops program from running the octopus.exe program")

	http.HandleFunc("/", recieveInfo)
	err = http.ListenAndServe(":8090", nil)

	if err != nil {
		WindowsLog.Error(err)
	}
}
