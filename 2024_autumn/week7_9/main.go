package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/go-redis/redis/v8"
	api "github.com/nikhovas/grpc_course/2024_autumn/week7/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Server struct {
	api.UnimplementedKeyValueServiceServer
	redisClient *redis.Client
}

func NewServer() *Server {
	redisURL, exists := os.LookupEnv("REDIS_URL")
	if !exists {
		redisURL = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{Addr: redisURL})
	return &Server{redisClient: rdb}
}

func (s *Server) SetValue(ctx context.Context, req *api.SetValueRequest) (*api.SetValueResponse, error) {
	err := s.redisClient.Set(ctx, req.Key, req.Value, 0).Err()
	return &api.SetValueResponse{}, err
}

func (s *Server) GetValue(ctx context.Context, req *api.GetValueRequest) (*api.GetValueResponse, error) {
	val, err := s.redisClient.Get(ctx, req.Key).Result()
	if err == redis.Nil {
		return nil, status.Error(codes.NotFound, "key not found")
	} else if err != nil {
		return nil, err
	}

	return &api.GetValueResponse{Value: val}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	api.RegisterKeyValueServiceServer(s, NewServer())
	reflection.Register(s)
	log.Println("gRPC server running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
