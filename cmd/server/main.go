package main

import (
	"flag"

	"github.com/rikatz/coraza-grpc/pkg/server"
	"go.uber.org/zap"
)

var (
	cfgPath *string
	port    *int
)

func main() {
	cfgPath = flag.String("cfgpath", "config", "defines the path of coraza configuration")
	port = flag.Int("port", 10000, "defines the port of gRPC server")
	flag.Parse()

	logcfg := zap.NewDevelopmentConfig()
	logcfg.DisableStacktrace = true
	logger, err := logcfg.Build()
	if err != nil {
		panic("invalid logger configuration")
	}
	cfg := server.ServerConfig{
		Port:       *port,
		Logger:     logger,
		ConfigPath: *cfgPath,
	}
	srv, err := server.NewServerWithOpts(&cfg)
	if err != nil {
		logger.Fatal("error creating server", zap.Error(err))
	}

	if err := srv.Start(); err != nil {
		logger.Fatal("error creating server", zap.Error(err))
	}
}
