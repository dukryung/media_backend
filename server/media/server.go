package media

import (
	"context"
	"flag"
	"fmt"
	"github.com/dukryung/media_backend/server/types"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"io"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
)

type Server struct {
	MediaServer

	GrpcServer *grpc.Server
	config     types.AppConfig
	grpcMux    *runtime.ServeMux

	close chan bool
}

func NewServer(config types.AppConfig) *Server {
	return &Server{
		config:     config,
		GrpcServer: grpc.NewServer(),
	}
}

func (s *Server) Run() {

	go s.RunGateway()

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", s.config.Server.GRPCAddress))
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

	err := http.ListenAndServe(fmt.Sprintf(":%s", s.config.Server.GatewayAddress), s.grpcMux)
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func (s *Server) registerHandler() {

	allowCors := func(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		return nil
	}

	s.grpcMux = runtime.NewServeMux(runtime.WithForwardResponseOption(allowCors))
	err := s.grpcMux.HandlePath("POST", "/request/file", handlerBinaryFileUpload)
	if err != nil {
		print(err.Error())
	}
	grpcServerEndpoint := flag.String("grpc-server-endpoint", fmt.Sprintf("localhost:%s", s.config.Server.GRPCAddress), "gRPC server endpoint")

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = RegisterMediaHandlerFromEndpoint(context.Background(), s.grpcMux, *grpcServerEndpoint, opts)
	if err != nil {
		log.Println("err : ", err)
	}

}

func (s *Server) RequestMedia(ctx context.Context, req *MediaRequest) (*MediaResponse, error) {

	imgFile, err := os.Create("./test.jpg")
	if err != nil {
		log.Println("err : ", err)
		return nil, err
	}
	defer imgFile.Close()

	_, err = imgFile.Write(req.Data)
	if err != nil {
		log.Println("err : ", err)
		return nil, err
	}

	return &MediaResponse{Code: "200"}, nil

}

func handlerBinaryFileUpload(w http.ResponseWriter, r *http.Request, params map[string]string) {
	log.Println("upload file")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse form: %s", err.Error()), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("attachment")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get file 'attachment': %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer file.Close()

	f, err := os.OpenFile("./downloaded.mp4", os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()

	io.Copy(f, file)

}
