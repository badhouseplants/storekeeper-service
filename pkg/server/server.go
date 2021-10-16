package server

import (
	"fmt"
	"net"

	"github.com/droplez/droplez-uploader/pkg/queue"
	"github.com/droplez/droplez-uploader/pkg/rpc/uploader"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Serve starts grpc server
func Serve() (err error) {
	unlocked := make(chan bool)
	// go queue.FreeSpaceWatcher(writeLock)
	// Read Config
	// Starting grpc server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", "9090"))
	if err != nil {
		return err
	}
	go queue.FreeSpaceWatcher(unlocked)
	// opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(1024 * 1024)}
	opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(2147483648)}
	grpcServer := grpc.NewServer(opts...)
	uploader.Register(grpcServer, unlocked)
	reflection.Register(grpcServer)

	if err = grpcServer.Serve(listener); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
