package main

import (
	"context"
	"log"
	"math"
	"net"

	api "github.com/nikhovas/grpc_course/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type CalcServer struct {
	api.UnimplementedCalcServerServer
}

func (s *CalcServer) CalcDistance(
	ctx context.Context, req *api.CalcDistanceReq,
) (*api.CalcDistanceRsp, error) {
	a := math.Pow(req.First.Latitude-req.Second.Latitude, 2.0)
	b := math.Pow(req.First.Longitude-req.Second.Longitude, 2.0)
	return &api.CalcDistanceRsp{
		Distance: math.Sqrt(a + b),
	}, nil
}

func main() {
	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := &CalcServer{}

	grpcServer := grpc.NewServer()
	api.RegisterCalcServerServer(grpcServer, s)
	reflection.Register(grpcServer)

	log.Println("gRPC server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
