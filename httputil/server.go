package httputil

import (
	"fmt"
	"github.com/alexandria-oss/core/config"
	"net/http"
	"time"
)

// DefaultServer exports an HTTP server with default configuration
func DefaultServer(cfg *config.Kernel, h http.Handler) *http.Server {
	return &http.Server{
		Addr:              cfg.Transport.HTTPHost + fmt.Sprintf(":%d", cfg.Transport.HTTPPort),
		Handler:           h,
		TLSConfig:         nil,
		ReadTimeout:       time.Second * 10,
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Second * 10,
		IdleTimeout:       time.Second * 15,
		MaxHeaderBytes:    4096,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}
}
