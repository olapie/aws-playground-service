package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"time"
)

func main() {
	lambda.Start(handleRequest)
}

type Response struct {
	DateTime        time.Time `json:"date_time"`
	Timestamp       int64     `json:"timestamp"`
	TimestampMillis int64     `json:"timestamp_millis"`

	ReceivedRequest map[string]any `json:"received_request"`
}

func handleRequest(ctx context.Context, event json.RawMessage) (*Response, error) {
	var resp Response
	if err := json.Unmarshal(event, &resp.ReceivedRequest); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
	now := time.Now()
	resp.DateTime = now
	resp.Timestamp = now.Unix()
	resp.TimestampMillis = now.UnixMilli()
	return &resp, nil
}
