package main

import (
	"Shorty.Server.Go.Mangment/cmd/api/server"
	"Shorty.Server.Go.Mangment/pkg/logger"
	"github.com/sirupsen/logrus"
	"runtime"
)

func main() {
	rtCpu := runtime.NumCPU()

	if rtCpu > 2 {
		runtime.GOMAXPROCS(rtCpu / 2)
		rtCpu = rtCpu / 2
	}

	app, err := server.NewApp()
	if err != nil {
		logger.Panic(err.Error(), nil)
	}

	logger.InfoF("Starting server with %d CPUs", nil, rtCpu)
	if err := app.Start(); err != nil {
		logger.Panic("Server error", logrus.Fields{"error": err.Error()})
	}
}
