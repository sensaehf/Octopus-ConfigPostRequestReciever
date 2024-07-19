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
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var DevFlag *bool                 // Feature flag that stop running the octopus program and Logging to the event viewer
var WindowsLog eventlogger.Logger // Custom Library to send logs/events to the Windows Event Logs

type Payload struct {
	ScanFileName    string `json:"scanFileName"`
	ScanDescription string `json:"ScanDescription"`
	Address         string `json:"Address"`
	Username        string `json:"Username"`
	Password        string `json:"Password"`
	DomainScan      bool   `json:"DomainScan"`
	Customer        string `json:"Customer"`
}

// Make sure to string does not contain any Sensetive data
func (p Payload) String() string {
	return fmt.Sprintf("scanfileName:%s, ScanDescription:%s, Address:%s, Username:%s, Customer:%s", p.ScanFileName, p.ScanDescription, p.Address, p.Username, p.Customer)
}

func callOctopus(p Payload) {
	var dScan string // Convert Bool to Yes/No for CLI
	if p.DomainScan {
		dScan = "Yes"
	} else {
		dScan = "No"
	}
	args := []string{
		fmt.Sprintf("-scanFileName %s", p.ScanFileName),
		fmt.Sprintf("-scanDescription %s", p.ScanDescription),
		fmt.Sprintf("-ConfPassword %s", os.Getenv("OCTOPUS_KEY")),
		fmt.Sprintf("-Address %s", p.Address),
		fmt.Sprintf("-Username %s", p.Username),
		fmt.Sprintf("-Password %s", p.Password),
		fmt.Sprintf("-Domainscan %s", dScan),
		fmt.Sprintf("Customer %s", p.Customer)}

	procAttr := new(os.ProcAttr)
	os.StartProcess("OctopusConfigurator.exe", args, procAttr)

}

func verifyInputs(p Payload) bool {
	match := true

	if !strings.Contains("srv.ocscanner", p.Username) ||
		len(p.Customer) <= 0 ||
		len(p.Address) <= 0 ||
		len(p.ScanDescription) <= 0 ||
		len(p.ScanFileName) <= 0 {
		return false
	}

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

	if verifyInputs(p) {
		callOctopus(p)
	} else {
		p.Password = "SECRET"
		WindowsLog.Warning(fmt.Errorf("input failed to parse for payload %s", p.String()))
	}

}

func main() {
	WindowsLog.Init()

	port := "8090"
	if os.Getenv("ASPNETCORE_PORT") != ":8090" { // get enviroment variable that set by ACNM 
		port = os.Getenv("ASPNETCORE_PORT")
	}

	err := godotenv.Load("./.env")
	if err != nil {
		WindowsLog.Error(fmt.Errorf("failed to load .env file, %s", err.Error()))
		os.Exit(1)
	}

	DevFlag = flag.Bool("Dev", false, "Turns on Dev mode. Stops program from running the octopus.exe program")
	flag.Parse()

	http.HandleFunc("/", recieveInfo)


	
	err = http.ListenAndServe(port, nil)

	if err != nil {
		WindowsLog.Error(err)
	}
}
