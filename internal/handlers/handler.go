package handlers

import (
	"BalancerService/config"
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	pb "BalancerService/proto/service"
)

type BalancerHandler struct {
	config     *config.Config
	requestCtr uint64
	cache      sync.Map
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

	if cachedURL, ok := b.cache.Load(videoUrl); ok {
		return &pb.RedirectResponse{
			RedirectUrl: cachedURL.(string),
		}, nil
	}

	var builder strings.Builder
	builder.WriteString("http://")
	builder.WriteString(b.config.CDNHost)
	builder.WriteString("/")
	builder.WriteString(server)
	builder.WriteString("/")
	builder.WriteString(path)

	cdnURL := builder.String()

	b.cache.Store(videoUrl, cdnURL)
	return &pb.RedirectResponse{
		RedirectUrl: cdnURL,
	}, nil
}
