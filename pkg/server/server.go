package server

import (
	"fmt"
	"log"
	"net"

	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	"github.com/corazawaf/coraza/v3"
	"github.com/corazawaf/coraza/v3/types"

	nginx "github.com/rikatz/coraza-grpc/apis/nginx"
	"github.com/rikatz/coraza-grpc/pkg/filter"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

type FilterServer struct {
	server *grpc.Server
	port   int
	waf    coraza.WAF
	logger *zap.Logger
}

type ServerConfig struct {
	Port       int
	Logger     *zap.Logger
	Opts       []grpc.ServerOption
	ConfigPath string
}

func NewServerWithOpts(cfg *ServerConfig) (*FilterServer, error) {
	var err error
	if cfg.Logger == nil {
		cfg.Logger, err = zap.NewDevelopmentConfig().Build()
		if err != nil {
			return nil, err
		}
	}

	// TODO: Accept the directives from file as multi slice string
	wafcfg := coraza.NewWAFConfig().
		WithDirectivesFromFile(cfg.ConfigPath + "/coraza.conf").
		WithDirectivesFromFile(cfg.ConfigPath + "/coreruleset/crs-setup.conf.example").
		WithDirectivesFromFile(cfg.ConfigPath + "/coreruleset/rules/*.conf").
		WithErrorCallback(logError(cfg.Logger))

	waf, err := coraza.NewWAF(wafcfg)
	if err != nil {
		cfg.Logger.Error("unable to create waf instance", zap.Error(err))
		return nil, err
	}

	server := &FilterServer{
		server: grpc.NewServer(cfg.Opts...),
		waf:    waf,
		port:   cfg.Port,
		logger: cfg.Logger,
	}

	return server, nil
}

func (s *FilterServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.logger.Info("starting gRPC Server", zap.Int("port", s.port))
	handler := &filter.GRPCHandler{
		WAF:    s.waf,
		Logger: *s.logger,
	}

	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(s.server, healthcheck)
	healthcheck.SetServingStatus("", healthgrpc.HealthCheckResponse_SERVING)

	// TODO: This should probably be inside a catch for Ctrl-C, interrupts instead of being here :) (need to deal with graceful shutdown)
	/*defer func() {
		healthcheck.SetServingStatus("", healthgrpc.HealthCheckResponse_NOT_SERVING)
	}()*/

	nginx.RegisterNginxFilterServer(s.server, handler)
	return s.server.Serve(lis)
}

// copied from coraza-spoa project :)
// We probably want to use the new slog from go v1.21 :D
func logError(logger *zap.Logger) func(rule types.MatchedRule) {
	return func(mr types.MatchedRule) {
		data := mr.ErrorLog()
		switch mr.Rule().Severity() {
		case types.RuleSeverityEmergency:
			logger.Error(data)
		case types.RuleSeverityAlert:
			logger.Error(data)
		case types.RuleSeverityCritical:
			logger.Error(data)
		case types.RuleSeverityError:
			logger.Error(data)
		case types.RuleSeverityWarning:
			logger.Warn(data)
		case types.RuleSeverityNotice:
			logger.Info(data)
		case types.RuleSeverityInfo:
			logger.Info(data)
		case types.RuleSeverityDebug:
			logger.Debug(data)
		}
	}
}
