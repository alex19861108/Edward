package taurusrpc

import (
	"git-pd.megvii-inc.com/liuwei02/Edward/writer"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type server struct{}

func (s *server) SearchXID(ctx context.Context, info *SearchXIDInfo) (*SearchXIDRes, error) {
	writer.TextWriter(info)
	return &SearchXIDRes{
		Xid:             "",
		StatusCode:      127, // As negative code represents meaningful circumstance,
		ChildStatusCode: 127, // Here, use 127 as the meaningless return
		Suspicious:      true,
	}, nil
}

func InitRPCServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterTaurusServer(s, &server{})
	// Register reflection service on gRPC server
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
