//go:build windows
// +build windows

package eventlogger

import (
	"golang.org/x/sys/windows/svc/eventlog"
)

type Logger struct {
	Name string
}

func (l Logger) Init() {
	l.Name = "OctopussPOSTReciever"
	const supports = eventlog.Error | eventlog.Warning | eventlog.Info
	eventlog.InstallAsEventCreate(l.Name, supports)
}

func (l Logger) Info(err error) {
	log, _ := eventlog.Open(l.Name)
	defer log.Close()
	err = log.Info(1, err.Error())
}

func (l Logger) Warning(err error) {
	log, _ := eventlog.Open(l.Name)
	defer log.Close()
	err = log.Info(2, err.Error())
}

func (l Logger) Error(err error) {
	log, _ := eventlog.Open(l.Name)
	defer log.Close()
	err = log.Info(3, err.Error())
}
