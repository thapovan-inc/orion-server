package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thapovan-inc/orion-server/util"
	"github.com/thapovan-inc/orionproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"time"
)

type grpcServer struct{}

type httpServer struct {
	r         *gin.Engine
	srv       *http.Server
	Port      uint16
	isRunning bool
}

var gServer *grpc.Server
var hServer *httpServer

func StartGRPC(port uint16) {
	logger := util.GetLogger("server", "StartGRPC")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatalln("gRPC server failed to listen: ", err)
	}
	gServer = grpc.NewServer()
	orionproto.RegisterTracerServer(gServer, grpcServer{})
	// Register reflection service on gRPC server.
	reflection.Register(gServer)
	go func() {
		if err := gServer.Serve(lis); err != nil {
			gServer = nil
			logger.Fatalf("gRPC server failed to start: %v", err)
		}
	}()
	logger.Infoln("gRPC server started and listening on port ", port)
}

func StopGRPC() {
	logger := util.GetLogger("server", "StopGRPC")
	if gServer != nil {
		logger.Infoln("Gracefully stopping gRPC server")
		gServer.GracefulStop()
	}
	logger.Info("Stopped gRPC server")
}

func StartHTTP(port uint16) {
	logger := util.GetLogger("server", "StartHTTP")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	hServer = &httpServer{r: r, Port: port, isRunning: false}
	r.POST("/uploadSpan", hServer.UploadSpan)
	r.POST("/bulkUpload", hServer.UploadSpanBulk)
	srv := &http.Server{
		Addr:    fmt.Sprint(":", port),
		Handler: hServer.r,
	}
	go func() {
		defer func() { hServer.isRunning = false }()
		err := srv.ListenAndServe()
		if err != nil {
			logger.Fatalln("HTTP server failed to start: ", err)
		}
	}()
	hServer.isRunning = true
}

func StopHTTP() {
	logger := util.GetLogger("server", "StopHTTP")
	if hServer != nil {
		logger.Infoln("Gracefully stopping HTTP server")
		ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
		hServer.srv.Shutdown(ctx)
		hServer.isRunning = false
	}
	logger.Infoln("Stopped HTTP Server")
}
