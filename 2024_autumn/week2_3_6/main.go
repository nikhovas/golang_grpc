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

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	api "github.com/nikhovas/grpc_course/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type CalcServer struct {
	api.UnimplementedCalcServerServer
	Logger        *zap.Logger
	SimpleCounter *prometheus.CounterVec
	Histogram     *prometheus.HistogramVec
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

	res := math.Sqrt(a + b)
	s.Logger.Info("Result", zap.Float64("value", res))
	err := errors.New("some error")
	s.Logger.Error("Error while calculating", zap.Error(err))

	s.SimpleCounter.WithLabelValues(authData[0]).Inc()
	s.Histogram.WithLabelValues("demo").Observe(res)

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

	counter2 := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "histogram_counter",
			Buckets: []float64{0, 10.0, 20.0, 30.0},
		},
		[]string{"label"},
	)
	prometheus.MustRegister(counter2)

	go func() {
		server := &http.Server{
			Addr:    ":9000",
			Handler: promhttp.Handler(),
		}

		log.Println("Serving metrics on http://0.0.0.0:9000")
		log.Fatalln(server.ListenAndServe())
	}()

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := &CalcServer{
		Logger:        logger,
		SimpleCounter: counter,
		Histogram:     counter2,
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_prometheus.UnaryServerInterceptor,
		),
	)
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
