package main

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nikhovas/grpc_course/2024_autumn/week7/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestStress(t *testing.T) {
	var (
		workersCount = 20000
		timeout      = 100 * time.Millisecond
	)

	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := api.NewKeyValueServiceClient(conn)

	_, _ = client.SetValue(context.Background(), &api.SetValueRequest{Key: "f", Value: "v"})

	var totalRequests atomic.Uint64
	var successRequests atomic.Uint64
	var errorRequests atomic.Uint64
	var totalDurationTime atomic.Uint64

	wg := sync.WaitGroup{}
	wg.Add(workersCount)
	for i := 0; i < workersCount; i++ {
		go func() {
			defer wg.Done()

			start := time.Now()
			ctx, _ := context.WithTimeout(context.Background(), timeout)
			_, err := client.GetValue(ctx, &api.GetValueRequest{Key: "f"})
			elapsed := time.Since(start)

			totalRequests.Add(1)
			if err != nil {
				errorRequests.Add(1)
			} else {
				successRequests.Add(1)
			}

			totalDurationTime.Add(uint64(elapsed / time.Millisecond))
		}()
	}

	wg.Wait()

	avgTime := time.Duration(totalDurationTime.Load()/totalRequests.Load()) * time.Millisecond

	t.Log("Stress test completed")
	t.Logf("Average time: %v", avgTime)
	t.Logf("Total requests: %d", totalRequests.Load())
	t.Logf("Success requests: %d", successRequests.Load())
	t.Logf("Error requests: %d", errorRequests.Load())
}
