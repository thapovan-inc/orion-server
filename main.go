package main

import (
	"fmt"
	"github.com/thapovan-inc/orion-server/publisher"
	"github.com/thapovan-inc/orion-server/server"
	"github.com/thapovan-inc/orion-server/util"
	"os"
	"sync"
)

func main() {
	fmt.Println("Loading config file from default.toml")
	util.LoadConfigFromFile("default.toml")

	util.SetupLoggerConfig()

	logger := util.GetLogger("main", "main")

	err := publisher.InitSpanPublisherFromConfig()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	var wg sync.WaitGroup
	logger.Infof("Starting gRPC server on port 20691")
	server.StartGRPC(20691)
	wg.Add(1)
	logger.Infof("Starting HTTP server on port 20691")
	server.StartHTTP(9017)
	wg.Add(1)
	defer func() {
		server.StopGRPC()
		wg.Done()
		server.StopHTTP()
		wg.Done()
	}()
	wg.Wait()
}
