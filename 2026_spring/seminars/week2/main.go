package main

import (
	"context"
	"embed"
	"log"
	"net"
	"net/http"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/nikhovas/grpc_course/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type CalcService struct {
	api.UnimplementedCalcServiceServer
	Logger    *zap.Logger
	Counter   *prometheus.CounterVec
	Histogram *prometheus.HistogramVec
}

func (cs *CalcService) Add(ctx context.Context, req *api.AddReq) (*api.AddRsp, error) {
	// md, _ := metadata.FromIncomingContext(ctx)

	// authData := md.Get("authorization")
	// if len(authData) == 0 {
	// 	return nil, status.Error(codes.Unauthenticated, "bad auth")
	// } else {
	// 	fmt.Println(authData[0])
	// }

	sum := req.A + req.B

	cs.Logger.Info("Result", zap.Float64("result", float64(sum)))
	cs.Counter.WithLabelValues("meta").Inc()
	cs.Histogram.WithLabelValues("meta").Observe(float64(sum))

	return &api.AddRsp{
		Sum: sum,
	}, nil
}

// ПРОБЕЛ ПОСЛЕ // НЕ СТАВИТЬ!!!

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
			Name:    "histogram_metric",
			Buckets: []float64{0, 10.0, 15.0, 20.0, 200.0},
		},
		[]string{"label"},
	)
	prometheus.MustRegister(histogram)

	go func() {
		server := &http.Server{
			Addr:    ":9000",
			Handler: promhttp.Handler(),
		}
		log.Fatalln(server.ListenAndServe())
	}()

	cs := &CalcService{Logger: logger, Counter: counter, Histogram: histogram}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_prometheus.UnaryServerInterceptor,
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

	// TODO

	gwMux := runtime.NewServeMux()
	api.RegisterCalcServiceHandler(context.Background(), gwMux, conn)

	mux := http.NewServeMux()
	mux.Handle("/", gwMux)

	mux.HandleFunc("/swagger-ui/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(swaggerData)
	})

	// mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.FS(swaggerFiles))))
	mux.Handle("/swagger-ui/", http.FileServer(http.FS(swaggerFiles)))

	gwServer := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
