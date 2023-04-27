package proxy

import (
	"errors"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"math"
	"net/http"
	"os"
	"riven-gateway/config"
	"strconv"
	"time"
)

var (
	readHeaderTimeout = time.Second * 10
	readTimeout       = time.Second * 15
	writeTimeout      = time.Second * 15
	idleTimeout       = time.Second * 120
)

func init() {
	var err error
	if v := os.Getenv("PROXY_READ_HEADER_TIMEOUT"); v != "" {
		if readHeaderTimeout, err = time.ParseDuration(v); err != nil {
			panic(err)
		}
	}
	if v := os.Getenv("PROXY_READ_TIMEOUT"); v != "" {
		if readTimeout, err = time.ParseDuration(v); err != nil {
			panic(err)
		}
	}
	if v := os.Getenv("PROXY_WRITE_TIMEOUT"); v != "" {
		if writeTimeout, err = time.ParseDuration(v); err != nil {
			panic(err)
		}
	}
	if v := os.Getenv("PROXY_IDLE_TIMEOUT"); v != "" {
		if idleTimeout, err = time.ParseDuration(v); err != nil {
			panic(err)
		}
	}
}

type Server struct {
	http.Server
}

func NewProxy(handler http.Handler, bc *config.Bootstrap) *Server {
	return &Server{
		Server: http.Server{
			Addr: buildAddr(bc.Server.HostPreFix, bc.Server.Port),
			Handler: h2c.NewHandler(handler, &http2.Server{
				IdleTimeout:          idleTimeout,
				MaxConcurrentStreams: math.MaxUint32,
			}),
			ReadTimeout:       readTimeout,
			ReadHeaderTimeout: readHeaderTimeout,
			WriteTimeout:      writeTimeout,
			IdleTimeout:       idleTimeout,
		},
	}
}

func buildAddr(prefix string, port int32) string {
	if prefix != "" {
		//TODO 获取网卡信息
		return ""
	}
	return "localhost:" + strconv.Itoa(int(port))
}

// Start the server.
func (s *Server) Start(ctx context.Context) error {
	err := s.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// Stop the server.
func (s *Server) Stop(ctx context.Context) error {
	return s.Shutdown(ctx)
}
