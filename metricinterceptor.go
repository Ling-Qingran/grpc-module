package metricinterceptor

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"time"
)

const (
	serverURL      = "http://34.86.236.100/"
	influxDBToken  = "I_UycfPULIG3VFr6eT-b0EzSIESMVb6rxZlS3n49zwHAcmpjPXQPS4u0eaZNY69hsWIVErE--T3lodcHQyx5rA=="
	influxDBOrg    = "API-Observability"
	influxDBBucket = "combined_metrics"
)

func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)
	// Measure request and response size (assuming they can be converted to string)
	reqSize := proto.Size(req.(proto.Message))
	respSize := 0
	if resp != nil {
		respSize = proto.Size(resp.(proto.Message))
	}

	// Get method name
	methodName := info.FullMethod
	statusCode := status.Code(err).String()

	// Extract peer information
	p, ok := peer.FromContext(ctx)
	ipAddress := ""
	if ok && p.Addr != net.Addr(nil) {
		host, _, err := net.SplitHostPort(p.Addr.String())
		if err == nil {
			ipAddress = host
		} else {
			log.Fatalf("Error while parsing peer address: %v", err)
		}
	}

	// Increment request count
	requestCount := 1

	// Error rate
	errorRate := 0
	if err != nil {
		errorRate = 1
	}

	// Extract metadata from context
	md, ok := metadata.FromIncomingContext(ctx)
	userAgent := ""
	if ok {
		// Metadata keys are normalized to lowercase
		if ua, exists := md["user-agent"]; exists && len(ua) > 0 {
			userAgent = ua[0]
		}
	}

	// Record metrics to InfluxDB (or print to console, log, etc.)
	writeMetrics(duration, err, reqSize, respSize, methodName, statusCode, requestCount, errorRate, ipAddress, userAgent)
	return resp, err
}

func writeMetrics(duration time.Duration, err error, reqSize, respSize int, methodName, statusCode string, requestCount, errorRate int, ipAddress, userAgent string) {
	// Create a new InfluxDB client
	client := influxdb2.NewClient(serverURL, influxDBToken)
	defer client.Close()

	// Create a write API (this can be reused)
	writeAPI := client.WriteAPI(influxDBOrg, influxDBBucket)

	// Create a point to write (measurement name is "gRPCMetrics")
	point := influxdb2.NewPoint(
		"gRPCMetrics",
		map[string]string{"endpoint": methodName, "ip_address": ipAddress, "user_agent": userAgent},
		map[string]interface{}{
			"duration":      duration.Seconds(),
			"error":         err != nil,
			"request_size":  reqSize,
			"response_size": respSize,
			"request_count": requestCount,
			"error_rate":    errorRate,
		},
		time.Now(),
	)

	// Write the point
	writeAPI.WritePoint(point)

	// Ensure data is written
	writeAPI.Flush()
}
