package main

import (
	"context"
	"log"
	"net"

	pb "github.com/siyopao/ipcheck/api/proto/v1"
	"github.com/siyopao/ipcheck/blocklist"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedIpCheckServer
}

func (s *server) InBlocklist(ctx context.Context, in *pb.InBlocklistRequest) (*pb.InBlocklistResponse, error) {
	ip := net.ParseIP(in.Ip)
	if ip == nil {
		return &pb.InBlocklistResponse{}, status.Errorf(codes.InvalidArgument, "%v is not a valid IP address", in.Ip)
	}

	inBlocklist, err := blocklist.InBlocklist(ip)
	if err != nil {
		return &pb.InBlocklistResponse{}, status.Errorf(codes.Internal, "%v", err)
	} else {
		return &pb.InBlocklistResponse{Ip: ip.String(), InBlocklist: inBlocklist}, nil
	}
}

func (s *server) InitBlocklists(ctx context.Context, in *pb.InitBlocklistsRequest) (*pb.InitBlocklistsResponse, error) {
	if err := blocklist.InitBlocklists(blocklist.Config); err != nil {
		log.Fatalf("error updating blocklist: %v", err)
	}
	return &pb.InitBlocklistsResponse{}, nil
}
