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
	"go.uber.org/zap"
)

var logger *zap.Logger

func GetLogger(packageName, function string) *zap.Logger {
	if logger == nil {
		SetupLoggerConfig()
	}
	l := logger.With(zap.String("package", packageName), zap.String("function", function))
	return l
}

func GetZapLogger(packageName, function string) {
	zap.NewProduction()
}

func SetupLoggerConfig() {
	//config := zap.NewDevelopmentConfig()
	//fmt.Printf("%+v", config)
	//l, err := config.Build()
	//if err != nil {
	//	fmt.Println(err)
	//	panic(err)
	//}
	logger, _ = zap.NewDevelopment()
}
