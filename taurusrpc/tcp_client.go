package taurusrpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func CSearchXID(address string, info *SearchXIDInfo) *SearchXIDRes {
	// Set up a connection to the server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}
	defer conn.Close()
	c := NewTaurusClient(conn)

	// Run API SearchXID
	r, err := c.SearchXID(context.Background(), info)
	if err != nil {
		log.Fatalf("Client error when run API SearchXID: %v", err)
	}
	return r
}
