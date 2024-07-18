//go:build windows
// +build windows

package eventlogger

import (
	"log"

	"golang.org/x/sys/windows/svc/eventlog"
)

type Logger struct {
	Name string
}

func (l Logger) Init() Logger {
	l.Name = "OctopussPOSTReciever"
	const supports = eventlog.Error | eventlog.Warning | eventlog.Info
	eventlog.InstallAsEventCreate(l.Name, supports)
	return l
}

func (l Logger) Info(err error) {
	Log, _ := eventlog.Open(l.Name)
	log.Printf(err.Error())
	defer Log.Close()
	err = Log.Info(1, err.Error())
}

func (l Logger) Warning(err error) {
	Log, _ := eventlog.Open(l.Name)
	log.Printf(err.Error())
	defer Log.Close()
	err = Log.Info(2, err.Error())
}

func (l Logger) Error(err error) {
	Log, _ := eventlog.Open(l.Name)
	log.Printf(err.Error())
	defer Log.Close()
	err = Log.Info(3, err.Error())
}
