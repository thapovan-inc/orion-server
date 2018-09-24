// Copyright 2018-Present Thapovan Info Systems Pvt. Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http:// www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
