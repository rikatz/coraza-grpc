package filter

import (
	"context"
	"fmt"
	"net/netip"

	pb "github.com/rikatz/coraza-grpc/apis/nginx"
	"go.uber.org/zap"
)

func (s *GRPCHandler) handleRequest(ctx context.Context, id string, version string, headers map[string]string, body []byte, request *pb.Request) (*pb.Decision, error) {
	if request == nil {
		return nil, fmt.Errorf("request information cannot be empty")
	}

	tx := s.WAF.NewTransactionWithID(id)
	s.Logger.Info("received request", zap.String("id", id))

	srcIP, err := netip.ParseAddr(request.GetSrcip())
	if err != nil {
		return nil, fmt.Errorf("invalid SRC IP address")
	}

	dstIP, err := netip.ParseAddr(request.GetDstip())
	if err != nil {
		return nil, fmt.Errorf("invalid DST IP address")
	}

	if request.Srcport == -1 || request.Dstport == -1 {
		return nil, fmt.Errorf("source and destination port should be defined")
	}

	if request.GetMethod() == "" {
		return nil, fmt.Errorf("method cannot be empty")
	}

	path := request.GetPath()
	if path == "" {
		path = "/"
	}

	query := request.GetQuery()

	for k, v := range headers {
		tx.AddRequestHeader(k, v)
	}

	if len(body) > 0 {
		interruption, _, err := tx.WriteRequestBody(body)
		if err != nil {
			return nil, err
		}
		if interruption != nil {
			return s.generateResponse(interruption), nil
		}
	}

	tx.ProcessConnection(srcIP.String(), int(request.GetSrcport()), dstIP.String(), int(request.GetDstport()))
	tx.ProcessURI(path+"?"+query, request.GetMethod(), "HTTP/"+version)

	if interruption := tx.ProcessRequestHeaders(); interruption != nil {
		return s.generateResponse(interruption), nil
	}

	it, err := tx.ProcessRequestBody()
	if err != nil {
		return nil, err
	}
	if it != nil {
		return s.generateResponse(it), nil
	}
	return s.generateResponse(nil), nil
}
