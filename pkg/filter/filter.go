package filter

import (
	"context"
	"fmt"

	"github.com/corazawaf/coraza/v3"
	pb "github.com/rikatz/coraza-grpc/apis/nginx"
)

// filterServer implements a new Protobuf filter server
type Handler interface {
	Handle(ctx context.Context, req *pb.FilterRequest) (*pb.Decision, error)
}

type GRPCHandler struct {
	pb.UnimplementedNginxFilterServer
	WAF coraza.WAF
	// TODO: Add the Cache
}

func (s *GRPCHandler) Handle(ctx context.Context, req *pb.FilterRequest) (*pb.Decision, error) {
	// Todo: Should the new cache and transaction implementation be done here? :D
	if req.GetId() == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	if req.GetVersion() == "" {
		return nil, fmt.Errorf("http version is mandatory")
	}

	switch operation := req.Operation.(type) {
	case *pb.FilterRequest_Request:
		return s.handleRequest(ctx, req.GetId(), req.GetVersion(), req.GetHeaders(), req.GetBody(), operation.Request)
	case *pb.FilterRequest_Response:
		return s.handleResponse(ctx, req.GetId(), req.GetVersion(), req.GetHeaders(), req.GetBody(), operation.Response)
	}

	return nil, nil
}
