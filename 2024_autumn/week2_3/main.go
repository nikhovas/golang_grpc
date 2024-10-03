package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"math"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	api "github.com/nikhovas/grpc_course/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type CalcServer struct {
	api.UnimplementedCalcServerServer
}

func (s *CalcServer) CalcDistance(
	ctx context.Context, req *api.CalcDistanceReq,
) (*api.CalcDistanceRsp, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	authData := md.Get("authorization")
	if len(authData) == 0 {
		return nil, errors.New("bad auth")
	} else {
		fmt.Println(authData[0])
	}

	a := math.Pow(req.First.Latitude-req.Second.Latitude, 2.0)
	b := math.Pow(req.First.Longitude-req.Second.Longitude, 2.0)
	return &api.CalcDistanceRsp{
		Distance: math.Sqrt(a + b),
	}, nil
}

//go:embed api/api.swagger.json
var swaggerData []byte

//go:embed swagger-ui
var swaggerFiles embed.FS

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

	go func() {
		log.Println("gRPC server started on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()

	ctx := context.Background()
	if err := api.RegisterCalcServerHandler(ctx, gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	mux := http.NewServeMux()

	mux.Handle("/", gwmux)

	mux.HandleFunc("/swagger-ui/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(swaggerData)
	})

	fSys, err := fs.Sub(swaggerFiles, "swagger-ui")
	if err != nil {
		panic(err)
	}

	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.FS(fSys))))

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: mux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
