package main

import (
	"GoStart/project/api"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"net"
	"time"
)

func main() {
	g := grpc.NewServer()
	api.RegisterPromotionServer(g, &Server{})

	lis, err := net.Listen("tcp", "0.0.0.0:8088")
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	err = g.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}

type Server struct {}

func (s *Server) GetLuList(ctx context.Context, request *api.GetLuListReq) (*api.GetLuListResponse, error) {
	time.Sleep(5*time.Second)
	return nil, status.Errorf(404, "not found")
	//return &api.GetLuListResponse{List: []*api.LuResponse{&api.LuResponse{Luid: 10, Luname: "tonna"}}}, nil
}
