package server

import (
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
			logger.Info("Stream EOF reached. Closing stream for namespace ", namespace)
			return nil
		}
		controlReq := request.GetControlRequest()
		if controlReq != nil {
			switch controlReq.RequestType {
			case orionproto.ControlRequest_END_STREAM:
				logger.Info("End stream message received. Closing stream for namespace ", namespace)
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
				logger.Info("Empty span data received. Closing stream for namespace ", namespace)
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
