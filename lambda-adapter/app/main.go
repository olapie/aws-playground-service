package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"go.olapie.com/conv"
	"go.olapie.com/logs"
	"go.olapie.com/types"
	"go.olapie.com/x/xapp"
)

func main() {
	httpAddr := flag.String("addr", "127.0.0.1:8000", "binding address of http server")
	xapp.Initialize("echo", httpAddr, nil)
	l, err := net.Listen("tcp", *httpAddr)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	err = http.Serve(l, http.HandlerFunc(handle))
	if errors.Is(err, http.ErrServerClosed) {
		slog.Info("server closed")
	} else {
		slog.Error(err.Error())
	}
}

type echoRequest struct {
	Host       string      `json:"host"`
	Method     string      `json:"method"`
	RequestURI string      `json:"request_uri"`
	Header     http.Header `json:"header"`
	Query      url.Values  `json:"query"`
	Body       []byte      `json:"body"`
}

type echoResponse struct {
	Request                echoRequest `json:"request"`
	ServerTime             time.Time   `json:"server_time"`
	ServerTimeSeconds      int64       `json:"server_time_seconds"`
	ServerTimeMilliSeconds int64       `json:"server_time_milli_seconds"`

	Message string `json:"message"`

	Version string `json:"version"`
}

func handle(rw http.ResponseWriter, req *http.Request) {
	resp := new(echoResponse)
	resp.Request = echoRequest{
		Host:       req.Host,
		Method:     req.Method,
		RequestURI: req.RequestURI,
		Header:     req.Header,
		Query:      req.URL.Query(),
		Body:       nil,
	}
	defer req.Body.Close()
	limitBodySize := int64(32 * types.KB)
	resp.Request.Body, _ = io.ReadAll(io.LimitReader(req.Body, limitBodySize))
	resp.ServerTime = time.Now()
	resp.ServerTimeSeconds = resp.ServerTime.Unix()
	resp.ServerTimeMilliSeconds = resp.ServerTime.UnixMilli()
	resp.Message = "this is an echo response"
	resp.Version = "0.1"
	statusCode := http.StatusOK

	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "*")
	rw.Header().Set("Access-Control-Expose-Headers", "*")

	if status := req.URL.Query().Get("status"); req.Method != http.MethodOptions && status != "" {
		code, err := conv.ToInt(status)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("invalid status value"))
			return
		} else {
			resp.Message = fmt.Sprintf("this is an echo response with custom status")
			statusCode = code
		}
	}

	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("cannot serialize response", logs.Err(err), slog.String("url", req.URL.String()), slog.Any("headers", req.Header))
		rw.WriteHeader(http.StatusInternalServerError)
		_, err = rw.Write([]byte(err.Error()))
	} else {
		slog.Info(string(data))
		rw.Header().Set("Content-Type", "application/json;charset=utf-8")
		rw.WriteHeader(statusCode)
		_, err = rw.Write(data)
	}
	if err != nil {
		slog.Error("cannot write response", logs.Err(err))
	}
}
