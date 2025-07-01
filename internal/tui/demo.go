package tui

import (
	"context"
	"fmt"
	"time"

	proto "github.com/codek7-services/codek7-tui/pkg/pb"
)

// Demo mode for testing without real gRPC server
func (v *Views) EnableDemoMode() {
	// Create fake gRPC client responses
	demoUser := &proto.UserResponse{
		Id:        "demo-user-123",
		Username:  "demo_user",
		Password:  "", // Don't store in real apps
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	demoVideos := []*proto.VideoMetadataResponse{
		{
			Id:          "video-1",
			UserId:      "demo-user-123",
			Title:       "My First Video",
			Description: "This is a demo video showing the upload functionality",
			CreatedAt:   time.Now().Add(-2 * time.Hour).Format("2006-01-02 15:04:05"),
			FileName:    "first_video.mp4",
		},
		{
			Id:          "video-2",
			UserId:      "demo-user-123",
			Title:       "Tutorial: Getting Started",
			Description: "A comprehensive tutorial for new users",
			CreatedAt:   time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05"),
			FileName:    "tutorial.mp4",
		},
		{
			Id:          "video-3",
			UserId:      "demo-user-123",
			Title:       "Advanced Features Demo",
			Description: "Showcasing advanced features of the platform",
			CreatedAt:   time.Now().Add(-30 * time.Minute).Format("2006-01-02 15:04:05"),
			FileName:    "advanced_demo.mp4",
		},
	}

	// Set demo data
	v.State.SetUser(demoUser)
	v.State.SetVideos(demoVideos)

	// Add some demo notifications
	demoNotifications := []Notification{
		{
			ID:      "notif-1",
			Type:    "upload",
			Message: "Video 'My First Video' uploaded successfully",
			Time:    time.Now().Add(-2 * time.Hour).Format("15:04:05"),
		},
		{
			ID:      "notif-2",
			Type:    "system",
			Message: "Welcome to CodeK7! Your account is ready.",
			Time:    time.Now().Add(-1 * time.Hour).Format("15:04:05"),
		},
		{
			ID:      "notif-3",
			Type:    "upload",
			Message: "Video 'Advanced Features Demo' processing complete",
			Time:    time.Now().Add(-30 * time.Minute).Format("15:04:05"),
		},
	}

	for _, notif := range demoNotifications {
		v.State.AddNotification(notif)
	}

	v.showMessage("Demo mode enabled! You are logged in as demo_user.")
	v.ShowDashboardView()
}

// Mock gRPC client for demo purposes
type MockRepoServiceClient struct{}

func (m *MockRepoServiceClient) CreateUser(ctx context.Context, req *proto.CreateUserRequest, opts ...interface{}) (*proto.UserResponse, error) {
	return &proto.UserResponse{
		Id:        fmt.Sprintf("user-%d", time.Now().Unix()),
		Username:  req.Username,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

func (m *MockRepoServiceClient) GetUser(ctx context.Context, req *proto.GetUserRequest, opts ...interface{}) (*proto.UserResponse, error) {
	// Simple demo authentication - accept any username/password combination
	return &proto.UserResponse{
		Id:        "demo-user-123",
		Username:  req.Username,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

func (m *MockRepoServiceClient) UploadVideo(ctx context.Context, opts ...interface{}) (interface{}, error) {
	// Mock upload - just return success
	return &proto.VideoMetadataResponse{
		Id:        fmt.Sprintf("video-%d", time.Now().Unix()),
		UserId:    "demo-user-123",
		Title:     "Uploaded Video",
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

func (m *MockRepoServiceClient) GetUserVideos(ctx context.Context, req *proto.GetUserVideosRequest, opts ...interface{}) (*proto.VideoListResponse, error) {
	// Return demo videos
	videos := []*proto.VideoMetadataResponse{
		{
			Id:          "video-1",
			UserId:      req.UserId,
			Title:       "My First Video",
			Description: "This is a demo video showing the upload functionality",
			CreatedAt:   time.Now().Add(-2 * time.Hour).Format("2006-01-02 15:04:05"),
			FileName:    "first_video.mp4",
		},
		{
			Id:          "video-2",
			UserId:      req.UserId,
			Title:       "Tutorial: Getting Started",
			Description: "A comprehensive tutorial for new users",
			CreatedAt:   time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05"),
			FileName:    "tutorial.mp4",
		},
	}

	return &proto.VideoListResponse{Videos: videos}, nil
}

// Implement other required methods with mock responses
func (m *MockRepoServiceClient) GetLast3UserVideos(ctx context.Context, req *proto.GetLast3UserVideosRequest, opts ...interface{}) (*proto.Video3ListResponse, error) {
	return &proto.Video3ListResponse{}, nil
}

func (m *MockRepoServiceClient) GetVideoByID(ctx context.Context, req *proto.GetVideoRequest, opts ...interface{}) (*proto.VideoMetadataResponse, error) {
	return &proto.VideoMetadataResponse{}, nil
}

func (m *MockRepoServiceClient) DownloadVideo(ctx context.Context, req *proto.DownloadVideoRequest, opts ...interface{}) (interface{}, error) {
	return nil, nil
}

func (m *MockRepoServiceClient) RemoveVideo(ctx context.Context, req *proto.GetVideoRequest, opts ...interface{}) (interface{}, error) {
	return nil, nil
}
