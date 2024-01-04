package main

import (
	"fmt"
	"net"
	pb "github.com/BravoLu/grpc_idl"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	grpc_validate "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc_exp_info/internal/services"
	"github.com/grpc_plugins/config"
	"github.com/grpc_plugins/log"
	"google.golang.org/grpc"
)

func main() {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_validate.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
		)))
	pb.RegisterExpInfoServerServer(s, &services.ExpInfoServerImpl{})
	conn, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", config.Config.Server.Port))
	if err != nil {
		log.Fatalf("invalid port")
	}
	if err := s.Serve(conn); err != nil {
		log.Fatalf("start server err: %+v", err)
	}
}
