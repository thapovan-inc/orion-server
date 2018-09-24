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
		logger.Debugln("request.SpanData is empty")
	}
	return &orionproto.ServerResponse{Success: isSuccess, Code: "", Message: errorMessage}, nil
}

func (httpServer) UploadSpanBulk(c *gin.Context) {
	logger := util.GetLogger("server", "httpServer::UploadSpanBulk")
	bulkRequest := &orionproto.BulkRequest{}
	err := orionproto.JsonToProto(c.Request.Body, bulkRequest)
	if err == nil {
		logger.Debugln(*bulkRequest)
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
			logger.Debugln("request.SpanData is empty")
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
