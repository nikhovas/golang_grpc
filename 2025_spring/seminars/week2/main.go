package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/nikhovas/grpc_course/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type ABA struct {
	api.UnimplementedCalcServerServer
}

func (c *ABA) Add(ctx context.Context, req *api.AddReq) (*api.AddRsp, error) {
	fmt.Println(req.Temp.Data)
	return &api.AddRsp{Result: req.A + req.B}, nil
}

func (c *ABA) Add2(ctx context.Context, req *api.AddReq2) (*api.AddRsp, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Fatal("why no?")
	}

	authHeaders := md.Get("Authorization")
	if len(authHeaders) != 1 {
		return nil, status.Error(codes.NotFound, "bad auth header")
		// return nil, errors.New("bad auth header")
	}

	fmt.Println("HERE")
	fmt.Println(authHeaders[0])

	var res int32 = 0
	for _, elem := range req.Values {
		res += elem
	}

	return &api.AddRsp{Result: res}, nil
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

	grpcServer := grpc.NewServer()

	c := ABA{}

	api.RegisterCalcServerServer(grpcServer, &c)
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	gwmux := runtime.NewServeMux()

	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	api.RegisterCalcServerHandler(context.Background(), gwmux, conn)

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
