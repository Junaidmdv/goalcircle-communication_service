package server

import (
	"fmt"
	"net"

	"github.com/Junaidmdv/goalcircle-communication_service/internal/config"
	notification_repo "github.com/Junaidmdv/goalcircle-communication_service/internal/domain/repository/notification"
	notification_handler "github.com/Junaidmdv/goalcircle-communication_service/internal/handler/grpc/notification"
	"github.com/Junaidmdv/goalcircle-communication_service/internal/infrastructure/persistence/postgres"
	notification_usecase "github.com/Junaidmdv/goalcircle-communication_service/internal/usecase/notification"
	"github.com/Junaidmdv/goalcircle-communication_service/pkg/logger"
	compb "github.com/Junaidmdv/goalcircle-protos/communication/v1"
	"google.golang.org/grpc"
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

	logger, err := logger.NewLogger()
	if err != nil {
		return err
	}

	notificatioRepo := notification_repo.NewNotificationRepository(psqldb.DB, logger)
	notificationUsecase := notification_usecase.NewNotificationUsecase(notificatioRepo, logger)
	notificationHandler := notification_handler.NewNotificationHandler(notificationUsecase, gs.Config.Server.TimeOut)

	compb.RegisterNotificationServiceServer(gs.Server, notificationHandler)

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
