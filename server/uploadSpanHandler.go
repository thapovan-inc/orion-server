package server

import "golang.org/x/net/context"
import (
	"github.com/gin-gonic/gin"
	"github.com/thapovan-inc/orion-server/authprovider"
	"github.com/thapovan-inc/orion-server/util"
	"github.com/thapovan-inc/orionproto"
)

func (grpcServer) UploadSpan(context context.Context, request *orionproto.UnaryRequest) (*orionproto.ServerResponse, error) {
	logger := util.GetLogger("server", "grpcServer::UploadSpan")
	logger.Debugln(*request)
	isSuccess := true
	namespace := ""
	var err error = nil
	namespace, err = authprovider.GetNameSpaceFromAuthToken(request.AuthToken)
	if err == nil {
		err = ingestSpan(request.SpanData, namespace)
	}
	isSuccess = err == nil
	errorMessage := ""
	if err != nil {
		errorMessage = err.Error()
	}
	return &orionproto.ServerResponse{Success: isSuccess, Code: "", Message: errorMessage}, nil
}

func (httpServer) UploadSpan(c *gin.Context) {
	logger := util.GetLogger("server", "httpServer::UploadSpan")
	unaryRequest := &orionproto.UnaryRequest{}
	err := orionproto.JsonToProto(c.Request.Body, unaryRequest)
	if err == nil {
		logger.Debugln(*unaryRequest)
		isSuccess := true
		namespace := ""
		var err error = nil
		namespace, err = authprovider.GetNameSpaceFromAuthToken(unaryRequest.AuthToken)
		if err == nil {
			err = ingestSpan(unaryRequest.SpanData, namespace)
		}
		isSuccess = err == nil
		errorMessage := ""
		if err != nil {
			errorMessage = err.Error()
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
