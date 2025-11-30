package main

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nikhovas/grpc_course/2024_autumn/week7/api"
	"google.golang.org/grpc"
)

func TestStress(t *testing.T) {
	var (
		serverAddr       = "localhost:50051"
		workerCount      = 5000
		timeout          = 100 * time.Millisecond
		workerIterations = 10000
		// workerIterations = 1
	)

	conn, _ := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	client := api.NewKeyValueServiceClient(conn)
	_, _ = client.SetValue(context.Background(), &api.SetValueRequest{Key: "f", Value: "v"})

	var errorCount atomic.Uint64
	var successCount atomic.Uint64
	var totalDurationTime atomic.Uint64

	wg := sync.WaitGroup{}
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < workerIterations; j++ {
				ctx, _ := context.WithTimeout(context.Background(), timeout)

				start := time.Now()
				_, err := client.GetValue(ctx, &api.GetValueRequest{Key: "f"})
				elapsed := time.Since(start)

				if err != nil {
					errorCount.Add(1)
				} else {
					successCount.Add(1)
				}

				totalDurationTime.Add(uint64(elapsed) / uint64(time.Nanosecond))
			}
		}()
	}

	wg.Wait()

	avgTime := time.Duration(totalDurationTime.Load()/(errorCount.Load()+successCount.Load())) * time.Nanosecond

	t.Log("Stress test completed")
	t.Logf("Average time: %v", avgTime)
	t.Logf("Success requests: %d", successCount.Load())
	t.Logf("Error requests: %d", errorCount.Load())
}
