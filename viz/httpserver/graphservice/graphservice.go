// Package graphservice implements a gRPC service for retrieving graphs.
//
// This is called directly by the web client.
//
// Guide to writing a server: https://grpc.io/docs/tutorials/basic/go/#server
package graphservice

import (
	"context"

	"github.com/golang/protobuf/proto"

	glog "github.com/golang/glog"
	pb "github.com/gonzojive/example-ts-go-grpc-bazel/viz/httpserver/frontendpb"
)

// Service is a gRPC implementation of FrontendService.
type Service struct{}

// GetGraph returns a graph object.
func (s *Service) GetGraph(ctx context.Context, req *pb.GetGraphRequest) (*pb.GetGraphResponse, error) {
	scale := req.GetScale()
	if scale == 0 {
		scale = 1
	}
	glog.Infof("got GetGraph request %s", proto.CompactTextString(req))
	return &pb.GetGraphResponse{
		Traces: []*pb.Trace{
			&pb.Trace{
				Points: []*pb.Trace_Point{
					&pb.Trace_Point{
						X: 1,
						Y: 2 * scale,
					},
					&pb.Trace_Point{
						X: 2,
						Y: 3 * scale,
					},
					&pb.Trace_Point{
						X: 3,
						Y: -1 * scale,
					},
					&pb.Trace_Point{
						X: 7,
						Y: 2 * scale,
					},
				},
				MarkerColor: "black",
				PlotlyMode:  "lines+points",
				PlotlyType:  "scatter",
			},
		},
	}, nil
}

var _ pb.FrontendServiceServer = &Service{}
