package main

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/grpc"

	api "github.com/nikhovas/grpc_course/2024_autumn/week7/api"
)

func TestStress(t *testing.T) {
	var (
		serverAddr   = "localhost:50051"
		workersCount = 10000
		timeout      = 100 * time.Millisecond
	)
	conn, _ := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	client := api.NewKeyValueServiceClient(conn)

	// set value
	_, _ = client.SetValue(context.Background(), &api.SetValueRequest{Key: "f", Value: "v"})

	var totalRequests atomic.Uint64
	var successRequests atomic.Uint64
	var errorRequests atomic.Uint64
	// unix nanos
	var totalDurationTime atomic.Uint64

	wg := sync.WaitGroup{}
	wg.Add(workersCount)
	for i := 0; i < workersCount; i++ {
		go func() {
			defer wg.Done()
			ctx, _ := context.WithTimeout(context.Background(), timeout)

			start := time.Now()
			_, err := client.GetValue(ctx, &api.GetValueRequest{Key: "f"})
			elapsed := time.Since(start)

			totalRequests.Add(1)
			if err != nil {
				errorRequests.Add(1)
			} else {
				successRequests.Add(1)
			}

			totalDurationTime.Add(uint64(elapsed / time.Nanosecond))
		}()
	}

	wg.Wait()

	avgTime := time.Duration(totalDurationTime.Load()/totalRequests.Load()) * time.Nanosecond

	t.Log("Stress test completed")
	t.Logf("Average time: %v", avgTime)
	t.Logf("Total requests: %d", totalRequests.Load())
	t.Logf("Success requests: %d", successRequests.Load())
	t.Logf("Error requests: %d", errorRequests.Load())
}
