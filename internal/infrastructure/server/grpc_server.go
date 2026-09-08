package server

import (
	"fmt"
	"net"

	"github.com/Junaidmdv/goalcircle-communication_service/internal/config"
	"github.com/Junaidmdv/goalcircle-communication_service/pkg/logger"
	"google.golang.org/grpc" 
	"github.com/Junaidmdv/goalcircle-communication_service/internal/infrastructure/persistence/postgres"
)

type GRPCServer struct {
	Config *config.Config
	Server *grpc.Server
}

func NewGRPCServer(cnfg *config.Config, logger logger.Logger) *GRPCServer {

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			RecoveryInterceptor(logger),
		),
	)

	return &GRPCServer{
		Server: server,
	}
}

func (gs *GRPCServer) BootstrapSetup() error {
	psqldb, err := postgres.NewPostgresDB(gs.Config.Postgres)
	if err != nil {
		return err
	}
	if err := psqldb.Migration(); err != nil {
		return err
	}

	return nil
}

func (gs *GRPCServer) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", gs.Config.Server.Port))

	if err != nil {

		return err
	}
	return gs.Server.Serve(lis)
}

func (gs *GRPCServer) GracefulShutdown() {
	gs.Server.GracefulStop()
}
