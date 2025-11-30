package main

import (
	"context"
	"embed"

	// "fmt"
	"log"
	"math"
	"net"
	"net/http"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/nikhovas/grpc_course/api"
)

type CalcServer struct {
	api.UnimplementedCalcServiceServer
	Logger    *zap.Logger
	Counter   *prometheus.CounterVec
	Histogram *prometheus.HistogramVec
}

func (cs *CalcServer) CalcDistance(ctx context.Context, req *api.CalcDistanceReq) (*api.CalcDistanceRsp, error) {
	a := math.Pow(req.First.Latitude-req.Second.Latitude, 2.0)
	b := math.Pow(req.First.Longitude-req.Second.Longitude, 2.0)

	res := math.Sqrt(a + b)

	// md, _ := metadata.FromIncomingContext(ctx)

	// authData := md.Get("authorization")
	// if len(authData) == 0 {
	// 	return nil, status.Error(codes.Unauthenticated, "bad auth")
	// } else {
	// 	fmt.Println(authData[0])
	// }

	cs.Logger.Info("Distance", zap.Float64("result", res))
	cs.Counter.WithLabelValues("meta").Inc()
	cs.Histogram.WithLabelValues("meta").Observe(res)

	return &api.CalcDistanceRsp{
		Distance: res,
	}, nil
}

//go:embed api/api.swagger.json
var swaggerData []byte

//go:embed swagger-ui
var swaggerFiles embed.FS

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "simple_counter",
		},
		[]string{"label"},
	)
	prometheus.MustRegister(counter)

	histogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "histogram",
			Buckets: []float64{0, 10.0, 20.0, 50.0},
		},
		[]string{"label"},
	)
	prometheus.MustRegister(histogram)

	go func() {
		server := &http.Server{
			Addr:    ":9000",
			Handler: promhttp.Handler(),
		}

		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal("Error while running metrics server", zap.Error(err))
		}
	}()

	cs := &CalcServer{Logger: logger, Counter: counter, Histogram: histogram}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_zap.UnaryServerInterceptor(logger),
			// grpc_prometheus.UnaryServerInterceptor,
		),
	)
	api.RegisterCalcServiceServer(grpcServer, cs)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		log.Println("gRPC server started on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	gwmux := runtime.NewServeMux()
	api.RegisterCalcServiceHandler(context.Background(), gwmux, conn)

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	mux.HandleFunc("/swagger-ui/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(swaggerData)
	})

	// Serve Swagger UI
	swaggerUIHandler := http.FileServer(http.FS(swaggerFiles))
	mux.Handle("/swagger-ui/", swaggerUIHandler)

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: mux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Println("Swagger UI available at: http://localhost:8090/swagger-ui/")
	log.Fatalln(gwServer.ListenAndServe())
}
