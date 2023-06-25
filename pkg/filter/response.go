package filter

import (
	"context"

	"github.com/corazawaf/coraza/v3/types"
	pb "github.com/rikatz/coraza-grpc/apis/nginx"
)

func (s *GRPCHandler) handleResponse(ctx context.Context, id string, version string, headers map[string]string, body []byte, response *pb.Response) (*pb.Decision, error) {
	var tx types.Transaction

	// TODO: Implement the cache for responses using transaction id.

	for k, v := range headers {
		tx.AddResponseHeader(k, v)
	}

	if interruption := tx.ProcessResponseHeaders(int(response.Statuscode), "HTTP/"+version); interruption != nil {
		return s.processInterruption(it, hit), nil
	}

	if len(body) > 0 {
		interruption, _, err := tx.WriteResponseBody(body)
		if err != nil {
			return nil, err
		}
		if interruption != nil {
			return s.processInterruption(interruption)
		}
	}

	tx.WriteResponseBody()
	interruption, err := tx.ProcessResponseBody()
	if err != nil {
		return nil, err
	}

	if interruption != nil {
		return s.processInterruption(it, hit), nil
	}
}
