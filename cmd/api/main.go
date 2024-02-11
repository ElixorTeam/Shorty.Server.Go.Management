package main

import (
	"Shorty.Server.Go.Mangment/cmd/api/server"
	"Shorty.Server.Go.Mangment/pkg/logger"
	"github.com/sirupsen/logrus"
	"runtime"
)

func main() {
	numCpu := runtime.NumCPU()
	logger.InfoF("Starting server with %d CPUs", nil, numCpu)

	if runtime.NumCPU() > 2 {
		runtime.GOMAXPROCS(numCpu / 2)
	}

	app, err := server.NewApp()
	if err != nil {
		logger.Panic(err.Error(), nil)
	}

	if err := app.Start(); err != nil {
		logger.Panic("Server error", logrus.Fields{"error": err.Error()})
	}
}
