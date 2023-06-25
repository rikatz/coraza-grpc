package filter

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/corazawaf/coraza/v3"
	"github.com/corazawaf/coraza/v3/types"
	pb "github.com/rikatz/coraza-grpc/apis/nginx"
	"go.uber.org/zap"
)

const (
	ActionAllow = iota
	ActionDeny
	ActionDrop
	ActionRedirect
)

// GRPCHandler implements a filter based on gRPC
// TODO: should create an interface for other implementations?
type GRPCHandler struct {
	pb.UnimplementedNginxFilterServer
	WAF    coraza.WAF
	Logger zap.Logger
	// TODO: Add the Cache
}

func (s *GRPCHandler) Handle(ctx context.Context, req *pb.FilterRequest) (*pb.Decision, error) {
	// Todo: Should the new cache and transaction implementation be done here? :D
	if req.GetId() == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	s.Logger.Info("received handle request", zap.String("id", req.GetId()))

	if req.GetVersion() == "" {
		return nil, fmt.Errorf("http version is mandatory")
	}

	switch operation := req.Operation.(type) {
	case *pb.FilterRequest_Request:
		return s.handleRequest(ctx, req.GetId(), req.GetVersion(), req.GetHeaders(), req.GetBody(), operation.Request)
	case *pb.FilterRequest_Response:
		return s.handleResponse(ctx, req.GetId(), req.GetVersion(), req.GetHeaders(), req.GetBody(), operation.Response)
	default:
		s.Logger.Warn("invalid operation called")
		return nil, fmt.Errorf("invalid operation called")
	}
}

func (s *GRPCHandler) generateResponse(interruption *types.Interruption) *pb.Decision {
	if interruption == nil {
		response := &pb.Decision{
			Action:     ActionAllow,
			Statuscode: http.StatusOK,
		}
		return response
	}

	var action int32
	var returncode int32
	switch interruption.Action {
	case strings.ToLower("drop"):
		action = ActionDrop
		returncode = http.StatusForbidden
	case strings.ToLower("redirect"):
		action = ActionRedirect
		returncode = http.StatusTemporaryRedirect
	default:
		action = ActionDeny
		returncode = http.StatusForbidden

	}
	return &pb.Decision{
		Action:     action,
		Decisionid: int32(interruption.RuleID),
		Statuscode: returncode,
		Message:    "request filtered by coraza waf",
	}

}
