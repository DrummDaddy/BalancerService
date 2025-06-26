package handlers

import (
	"context"
	"testing"

	"BalancerService/config"
	pb "BalancerService/proto/service"

	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {
	mockConfig := &config.Config{CDNHost: "cdn.mock.com"}

	tests := []struct {
		name         string
		request      *pb.RedirectRequest
		wantURL      string
		expectError  bool
		resetCounter uint64
	}{
		{
			name: "Valid video URL - redirected to CDN",
			request: &pb.RedirectRequest{
				Video: "http://origin-cluster/video/123/abc.m3u8",
			},
			wantURL:      "http://cdn.mock.com/origin-cluster/video/123/abc.m3u8",
			expectError:  false,
			resetCounter: 0,
		},
		{
			name: "Invalid video URL format",
			request: &pb.RedirectRequest{
				Video: "invalid_format",
			},
			wantURL:      "",
			expectError:  true,
			resetCounter: 0,
		},
		{
			name: "Original URL return on every 10th request",
			request: &pb.RedirectRequest{
				Video: "http://origin-cluster/video/123/abc.m3u8",
			},
			wantURL:      "http://origin-cluster/video/123/abc.m3u8",
			expectError:  false,
			resetCounter: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			handler := NewBalancerHandler(mockConfig)

			if tt.resetCounter > 0 {
				handler.requestCtr = tt.resetCounter
			}

			response, err := handler.Redirect(context.Background(), tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tt.wantURL, response.RedirectUrl)
			}

			if tt.resetCounter > 0 {
				assert.Equal(t, tt.resetCounter+1, handler.requestCtr, "Request counter did not increase correctly")
			}
		})
	}
}

func TestRedirectEvery10thRequest(t *testing.T) {
	mockConfig := &config.Config{CDNHost: "cdn.mock.com"}
	handler := NewBalancerHandler(mockConfig)

	videoUrl := "http://origin-cluster/video/123/abc.m3u8"
	expectedCdnUrl := "http://cdn.mock.com/origin-cluster/video/123/abc.m3u8"

	// Проверяем, что каждый 10-й запрос возвращает оригинальный URL
	for i := 1; i <= 20; i++ {
		response, err := handler.Redirect(context.Background(), &pb.RedirectRequest{Video: videoUrl})
		assert.NoError(t, err)

		if i%10 == 0 {
			assert.Equal(t, videoUrl, response.RedirectUrl, "Expected original URL on 10th request")
		} else {
			assert.Equal(t, expectedCdnUrl, response.RedirectUrl, "Expected CDN URL for non-10th request")
		}
	}
}
