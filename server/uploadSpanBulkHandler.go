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

package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thapovan-inc/orion-server/authprovider"
	"github.com/thapovan-inc/orion-server/util"
	"github.com/thapovan-inc/orionproto"
	"golang.org/x/net/context"
)

func (grpcServer) UploadSpanBulk(context context.Context, request *orionproto.BulkRequest) (*orionproto.ServerResponse, error) {
	logger := util.GetLogger("server", "grpcServer::UploadSpanBulk")
	isSuccess := true
	errorMessage := ""
	if len(request.SpanData) > 0 {
		namespace := ""
		var err error = nil
		namespace, err = authprovider.GetNameSpaceFromAuthToken(request.AuthToken)
		if err == nil {
			for _, spanData := range request.SpanData {
				err := ingestSpan(spanData, namespace)
				if err != nil {
					errorMessage = fmt.Sprintf("The following error occured when processing span with ID %s: %s",
						spanData.SpanId, err)
					isSuccess = false
					break
				}
			}
		} else {
			errorMessage = err.Error()
		}
	} else {
		logger.Debug("request.SpanData is empty")
	}
	return &orionproto.ServerResponse{Success: isSuccess, Code: "", Message: errorMessage}, nil
}

func (httpServer) UploadSpanBulk(c *gin.Context) {
	logger := util.GetLogger("server", "httpServer::UploadSpanBulk")
	bulkRequest := &orionproto.BulkRequest{}
	err := orionproto.JsonToProto(c.Request.Body, bulkRequest)
	if err == nil {
		isSuccess := true
		errorMessage := ""
		if len(bulkRequest.SpanData) > 0 {
			namespace := ""
			var err error = nil
			namespace, err = authprovider.GetNameSpaceFromAuthToken(bulkRequest.AuthToken)
			if err == nil {
				for _, spanData := range bulkRequest.SpanData {
					err := ingestSpan(spanData, namespace)
					if err != nil {
						errorMessage = fmt.Sprintf("The following error occured when processing span with ID %s: %s",
							spanData.SpanId, err)
						isSuccess = false
						break
					}
				}
			} else {
				errorMessage = err.Error()
			}
		} else {
			logger.Debug("request.SpanData is empty")
		}
		response, err := orionproto.ProtoToJson(&orionproto.ServerResponse{Success: isSuccess, Code: "", Message: errorMessage})
		if err == nil {
			c.Data(200, "application/json", []byte(response))
			return
		}
	}
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": err.Error(),
		})
	}
}
