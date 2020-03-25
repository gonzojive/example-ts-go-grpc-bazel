package graphservice

import (
	pb "github.com/gonzojive/example-ts-go-grpc-bazel/viz/httpserver/frontendpb"
)

type Service struct {}

// GetGraph returns a graph object.
func (s *Service) GetGraph(req *pb.GetGraphRequest) (*pb.GetGraphResponse, error) {
	return &pb.GetGraphResponse{
		Traces: []*pb.Trace{
			&pb.Trace{
				Points: []*pb.Trace_Point{
					&pb.Trace_Point{
						X: 1,
						Y: 2,
					},
					&pb.Trace_Point{
						X: 2,
						Y: 3,
					},
					&pb.Trace_Point{
						X: 2,
						Y: -1,
					},
				}],
				MarkerColor: "black",
				PlotlyMode: "lines+points",
				PlotlyType: "scatter",
			}
		},
	}, nil
}