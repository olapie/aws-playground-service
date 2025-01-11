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

	Request any
}

// Different type for requests from different sender, e.g. ApplicationLoadBalancer, APIGateway
// refer to https://pkg.go.dev/github.com/aws/aws-lambda-go/events#section-readme

func handleRequest(ctx context.Context, req *events.ALBTargetGroupRequest) (*events.ALBTargetGroupResponse, error) {
	var payload Payload
	now := time.Now()
	payload.DateTime = now
	payload.Timestamp = now.Unix()
	payload.TimestampMillis = now.UnixMilli()
	payload.Version = "v0.5"
	payload.Request = req

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
