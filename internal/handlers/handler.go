package handlers

import (
	"BalancerService/config"
	"context"
	"fmt"
	"strings"
	"sync/atomic"

	pb "BalancerService/proto/service"
)

type BalancerHandler struct {
	config     *config.Config
	requestCtr uint64
	pb.UnimplementedBalancerServiceServer
}

func NewBalancerHandler(config *config.Config) *BalancerHandler {
	return &BalancerHandler{
		config: config,
	}
}

func (b *BalancerHandler) Redirect(ctx context.Context, req *pb.RedirectRequest) (*pb.RedirectResponse, error) {
	videoUrl := req.Video

	count := atomic.AddUint64(&b.requestCtr, 1)

	parts := strings.SplitN(videoUrl, "/", 4)
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid video URL format")
	}

	server := strings.Split(parts[2], ".")[0]
	path := parts[3]

	if count%10 == 0 {
		return &pb.RedirectResponse{
			RedirectUrl: videoUrl,
		}, nil
	}
	cdnURL := fmt.Sprintf("http://%s/%s/%s", b.config.CDNHost, server, path)
	return &pb.RedirectResponse{
		RedirectUrl: cdnURL,
	}, nil

}
