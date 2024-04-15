package util

import (
	"log"
	"net/http"
)

// LoggingRoundTripper 结构封装了 http.RoundTripper
type LoggingRoundTripper struct {
	Proxied http.RoundTripper
}

// RoundTrip 执行单个HTTP事务并记录请求头
func (lrt *LoggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	log.Printf("Sending request to %s with headers: %v", r.URL, r.Header)
	return lrt.Proxied.RoundTrip(r)
}

// NewLoggingClient 创建并返回一个配置了日志记录功能的HTTP客户端
func NewLoggingClient() *http.Client {
	return &http.Client{
		Transport: &LoggingRoundTripper{
			Proxied: http.DefaultTransport,
		},
	}
}
