package main

import (
	"github.com/op/go-logging"
	"os"
)

func NewLogger(moduleName string, level logging.Level) (logger *logging.Logger) {
	logger = logging.MustGetLogger(moduleName)
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	format := `%{color}%{time:15:04:05.000} %{shortpkg}/%{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`
	formatter := logging.MustStringFormatter(format)
	backendFormatter := logging.NewBackendFormatter(backend, formatter)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(level, moduleName)
	logging.SetBackend(backendLeveled)
	return
}
