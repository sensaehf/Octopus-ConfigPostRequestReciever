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
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var DevFlag *bool // Feature flag that stop running the octopus program and Logging to the event viewer
var PortFlag *int64
var WindowsLog eventlogger.Logger // Custom Library to send logs/events to the Windows Event Logs

type Payload struct {
	ScanFileName    string `json:"ScanFileName"`
	ScanDescription string `json:"ScanDescription"`
	Address         string `json:"Address"`
	Username        string `json:"Username"`
	Password        string `json:"Password"`
	DomainScan      bool   `json:"DomainScan"`
	Customer        string `json:"Customer"`
}

// Make sure to string does not contain any Sensetive data
func (p Payload) String() string {
	return fmt.Sprintf("ScanfileName:%s, ScanDescription:%s, Address:%s, Username:%s, Customer:%s", p.ScanFileName, p.ScanDescription, p.Address, p.Username, p.Customer)
}

func callOctopus(p Payload) {
	var dScan string // Convert Bool to Yes/No for CLI
	if p.DomainScan {
		dScan = "Yes"
	} else {
		dScan = "No"
	}
	args := []string{
		fmt.Sprintf("-ScanFileName %s", p.ScanFileName),
		fmt.Sprintf("-ScanDescription %s", p.ScanDescription),
		fmt.Sprintf("-ConfPassword %s", os.Getenv("OCTOPUS_KEY")),
		fmt.Sprintf("-Address %s", p.Address),
		fmt.Sprintf("-Username %s", p.Username),
		fmt.Sprintf("-Password %s", p.Password),
		fmt.Sprintf("-Domainscan %s", dScan),
		fmt.Sprintf("-Customer %s", p.Customer)}

	attr := &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	os.StartProcess("OctopusConfigurator.exe", args, attr)

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
	DevFlag = flag.Bool("Dev", false, "Turns on Dev mode. Stops program from running the octopus.exe program")
	PortFlag = flag.Int64("port", 8090, "Selects what port is listned to")

	err := godotenv.Load("./.env")
	if err != nil {
		WindowsLog.Error(fmt.Errorf("failed to load .env file, %s", err.Error()))
		os.Exit(1)
	}

	flag.Parse()

	http.HandleFunc("/", recieveInfo)

	port := *PortFlag
	if os.Getenv("ASPNETCORE_PORT") != "" { // get enviroment variable that set by ACNM
		p := os.Getenv("ASPNETCORE_PORT")
		port, _ = strconv.ParseInt(p, 10, 32)

	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		WindowsLog.Error(err)
	}
}
