package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"code.olapie.com/private-repo-demo"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	private_repo_demo.Hello("hi")
	lambda.Start(handleRequest)
}

type Payload struct {
	Version string `json:"version"`

	DateTime        time.Time `json:"date_time"`
	Timestamp       int64     `json:"timestamp"`
	TimestampMillis int64     `json:"timestamp_millis"`

	Request struct {
		Body    string            `json:"body"`
		Headers map[string]string `json:"headers"`
		Method  string            `json:"method"`
		Path    string            `json:"path"`
	}
}

// Different type for requests from different sender, e.g. ApplicationLoadBalancer, APIGateway
// refer to https://pkg.go.dev/github.com/aws/aws-lambda-go/events#section-readme

func handleRequest(ctx context.Context, req *events.ALBTargetGroupRequest) (*events.ALBTargetGroupResponse, error) {
	var payload Payload
	now := time.Now()
	payload.DateTime = now
	payload.Timestamp = now.Unix()
	payload.TimestampMillis = now.UnixMilli()
	payload.Version = "v0.4"
	payload.Request.Headers = req.Headers
	payload.Request.Body = req.Body
	payload.Request.Method = req.HTTPMethod
	payload.Request.Path = req.Path
	data, err := json.Marshal(req)
	if err != nil {
		slog.Error("marshal error: " + err.Error())
	} else {
		slog.Info(string(data))
	}
	res := &events.ALBTargetGroupResponse{
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
