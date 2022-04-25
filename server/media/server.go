package media

import (
	"context"
	"fmt"
	"github.com/dukryung/media_backend/server/types"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"

)

type Server struct {
	GrpcServer *grpc.Server
	config     types.AppConfig
	grpcMux    *runtime.ServeMux

	close chan bool
}

func NewServer(config types.AppConfig) *Server {
	return &Server{
		config: config,
		GrpcServer: grpc.NewServer(),
	}
}

func (s *Server) Run() {
	go s.RunGateway()

	listen, err := net.Listen("tcp", s.config.Server.GRPCAddress)
	if err != nil {
		panic(err)
	}
	s.GrpcServer.Serve(listen)
}

func (s *Server) Close() {
	s.close <- true
}

func (s *Server) RunGateway() {
	s.registerHandler()

	gateway := http.Server{
		Addr:    s.config.Server.GatewayAddress,
		Handler: cors.Default().Handler(s.grpcMux),
	}

	gateway.ListenAndServe()
}

func (s *Server) registerHandler() {

	conn, err := grpc.DialContext(
		context.Background(),
		s.config.Server.GRPCAddress,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Println("err : ",err)
	}

	allowCors := func(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		return nil
	}
	s.grpcMux = runtime.NewServeMux(
		runtime.WithForwardResponseOption(allowCors),
	)

	err = RegisterMediaHandler(context.Background(),s.grpcMux,conn)
	if err != nil {
		log.Println("err : ",err)
	}

}

func (s *Server) RequestMedia(ctx context.Context, req *MediaRequest) (*MediaResponse, error) {
	log.Println("req")
	conn , err := grpc.DialContext(ctx, s.config.Server.GRPCAddress,grpc.WithInsecure())
	if err != nil {
		log.Println("err : ",err)
		return nil, err
	}

	defer conn.Close()

	data := req.Data

	fmt.Println("data : ",string(data))

	return &MediaResponse{ Code: "200" }, nil

}
