package main

import (
	"context"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handleRequest)
}

type Response struct {
	Version string `json:"version"`

	DateTime        time.Time `json:"date_time"`
	Timestamp       int64     `json:"timestamp"`
	TimestampMillis int64     `json:"timestamp_millis"`

	Request struct {
		Version string
		Headers map[string]string
		Cookies []string
	}
}

func handleRequest(ctx context.Context, req *events.APIGatewayV2HTTPRequest) (*Response, error) {
	var resp Response
	now := time.Now()
	resp.DateTime = now
	resp.Timestamp = now.Unix()
	resp.TimestampMillis = now.UnixMilli()
	resp.Version = "v0.3"
	resp.Request.Headers = req.Headers
	resp.Request.Cookies = req.Cookies
	resp.Request.Version = req.Version
	return &resp, nil
}
