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

package util

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

var logLevel logrus.Level

func GetLogger(packageName, function string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"PackageName": packageName,
		"Function":    function,
	})
}

func GetZapLogger(packageName, function string) {
	zap.NewProduction()
}

func SetupLoggerConfig() {
	config := GetConfig().Logger
	var err error
	if config.Format != "json" {
		shouldForceColors := config.UseColors
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: shouldForceColors})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	if config.Level != "" {
		logLevel, err = logrus.ParseLevel(config.Level)
		if err == nil {
			logrus.SetLevel(logLevel)
		}
	}
	if err != nil || config.Level == "" {
		logger := GetLogger("util", "SetupLoggerConfig")
		if err != nil {
			logger = logger.WithError(err)
		}
		logger.Warn("Setting default log level to info")
		logrus.SetLevel(logrus.InfoLevel)
	}
}
