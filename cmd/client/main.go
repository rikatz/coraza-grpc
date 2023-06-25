package main

import (
	"context"
	"flag"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/rikatz/coraza-grpc/apis/nginx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	address *string
)

func main() {
	address = flag.String("grpc-address", "127.0.0.1:10000", "defines the grpc server to consume, in format of address:port")
	flag.Parse()
	conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	runWAFTest(conn)
}

func runWAFTest(conn *grpc.ClientConn) {
	wafclient := nginx.NewNginxFilterClient(conn)
	ctx := context.Background()

	request := nginx.FilterRequest{
		Id:      "blabla",
		Version: "1.1",
	}

	_, err := wafclient.Handle(ctx, &request)
	if err != nil {
		log.Printf("got error: %s", err)
	}

	requestDenied := nginx.FilterRequest{
		Id:      "bloblo",
		Version: "1.1",
		Headers: map[string]string{
			"something": "bla",
			"blo":       "ble",
			"Host":      "sometest.com",
		},
		Operation: &nginx.FilterRequest_Request{
			Request: &nginx.Request{
				Srcip:   "1.1.1.1",
				Srcport: 32000,
				Dstip:   "8.8.8.8",
				Dstport: 80,
				Method:  "TRACE",
				Query:   "id=0",
				Path:    "/blorgh/../../../etc/passwd",
			},
		},
	}

	decision, err := wafclient.Handle(ctx, &requestDenied)
	if err != nil {
		log.Printf("got error: %s", err)
	}
	spew.Dump(decision)

	requestOK := nginx.FilterRequest{
		Id:      "bloblo123",
		Version: "1.1",
		Headers: map[string]string{
			"something": "bla",
			"blo":       "ble",
			"Host":      "sometest.com",
		},
		Operation: &nginx.FilterRequest_Request{
			Request: &nginx.Request{
				Srcip:   "1.1.1.1",
				Srcport: 32000,
				Dstip:   "8.8.8.8",
				Dstport: 80,
				Method:  "GET",
				Path:    "/",
			},
		},
	}

	decision, err = wafclient.Handle(ctx, &requestOK)
	if err != nil {
		log.Printf("got error: %s", err)
	}
	spew.Dump(decision)

}
