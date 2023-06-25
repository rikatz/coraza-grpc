package server

import (
	"fmt"
	"log"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"gorm.io/gorm/logger"

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
}

func NewServerWithOpts(port int, opts ...grpc.ServerOption) (*FilterServer, error) {
	conf := coraza.NewWAFConfig().
		WithDirectives(cfg.Directives).
		WithErrorCallback(logError(logger))

	waf, err := coraza.NewWAF(conf)
	if err != nil {
		logger.Error("unable to create waf instance", zap.Error(err))
		return nil, err
	}

	server := &FilterServer{
		server: grpc.NewServer(opts...),
		waf:    waf,
		port:   port,
	}

	return server, nil
}

func (s *FilterServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	handler := &filter.GRPCHandler{
		WAF: s.waf,
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
