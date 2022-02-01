package server

import (
	"fmt"
	"net"

	"github.com/badhouseplants/storekeeper-service/pkg/constants"
	"github.com/badhouseplants/storekeeper-service/pkg/internal/service"
	"github.com/badhouseplants/storekeeper-service/pkg/tools/logger"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Serve starts grpc server
func Serve() (err error) {
	// Init listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", viper.GetString(constants.ConstDroplezUploaderPort)))
	if err != nil {
		return
	}

	// Prepare gRPC server
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(2147483648),
	setupGrpcStreamOpts(),
		setupGrpcUnaryOpts(),
	)
	// Register services
	service.RegisterUploader(grpcServer)

	// Add reflection only if it's not production environment
	if viper.GetString(constants.ConstDroplezUploaderMode) != constants.ConstProdMode {
		reflection.Register(grpcServer)
	}

	// Start server
	if err = grpcServer.Serve(listener); err != nil {
		return
	}

	return
}

// Add interceptors for unary requests
func setupGrpcUnaryOpts() grpc.ServerOption {
	return grpc_middleware.WithUnaryServerChain(
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_logrus.UnaryServerInterceptor(logger.GrpcLogrusEntry, logger.GrpcLogrusOpts...),
	)
}

// Add interceptors for streams
func setupGrpcStreamOpts() grpc.ServerOption {
	return grpc_middleware.WithStreamServerChain(
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_logrus.StreamServerInterceptor(logger.GrpcLogrusEntry, logger.GrpcLogrusOpts...),
	)
}
