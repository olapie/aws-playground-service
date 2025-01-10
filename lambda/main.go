package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handleRequest)
}

type Payload struct {
	Version string `json:"version"`

	DateTime        time.Time `json:"date_time"`
	Timestamp       int64     `json:"timestamp"`
	TimestampMillis int64     `json:"timestamp_millis"`

	Request struct {
		Version string
		Headers map[string]string
		Cookies []string
		RawPath string
	}
}

func handleRequest(ctx context.Context, req *events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
	var payload Payload
	now := time.Now()
	payload.DateTime = now
	payload.Timestamp = now.Unix()
	payload.TimestampMillis = now.UnixMilli()
	payload.Version = "v0.3"
	payload.Request.Headers = req.Headers
	payload.Request.Cookies = req.Cookies
	payload.Request.Version = req.Version
	payload.Request.RawPath = req.RawPath
	data, err := json.Marshal(req)
	if err != nil {
		slog.Error("marshal error: " + err.Error())
	} else {
		slog.Info(string(data))
	}
	res := &events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		slog.Error("marshal error: " + err.Error())
	} else {
		res.Body = string(body)
	}
	return res, nil
}
