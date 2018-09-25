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
	"go.uber.org/zap"
	"io"

	"github.com/thapovan-inc/orion-server/authprovider"
	"github.com/thapovan-inc/orion-server/util"
	"github.com/thapovan-inc/orionproto"
)

func (grpcServer) UploadSpanStream(streamServer orionproto.Tracer_UploadSpanStreamServer) error {
	logger := util.GetLogger("server", "grpcServer::UploadSpanStream")
	namespace := ""
	for {
		request, err := streamServer.Recv()
		if err == io.EOF {
			logger.Info("Stream EOF reached. Closing stream for namespace ", zap.String("namespace", namespace))
			return nil
		}
		controlReq := request.GetControlRequest()
		if controlReq != nil {
			switch controlReq.RequestType {
			case orionproto.ControlRequest_END_STREAM:
				logger.Info("End stream message received. Closing stream for namespace ", zap.String("namespace", namespace))
				return nil
			case orionproto.ControlRequest_AUTH:
				namespace, err = authprovider.GetNameSpaceFromAuthToken(controlReq.GetJsonString())
				if err != nil {
					response := &orionproto.ServerResponse{
						Message: err.Error(),
					}
					streamServer.Send(response)
					return nil
				}
			}
		} else {
			spanData := request.GetSpanData()
			if spanData == nil {
				logger.Info("Empty span data received. Closing stream for namespace ", zap.String("namespace", namespace))
				response := &orionproto.ServerResponse{
					Message: "Empty span received. Closing stream now.",
				}
				streamServer.Send(response)
				return nil
			} else {
				err := ingestSpan(spanData, namespace)
				if err != nil {
					response := &orionproto.ServerResponse{
						Message: err.Error(),
					}
					streamServer.Send(response)
					return nil
				} else {
					streamServer.Send(&orionproto.ServerResponse{Success: true})
				}
			}
		}
	}
}
