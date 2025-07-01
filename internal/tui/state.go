package tui

import (
	"sync"

	proto "github.com/codek7-services/codek7-tui/pkg/pb"
)

// AppState holds the global application state
type AppState struct {
	mu            sync.RWMutex
	LoggedIn      bool
	CurrentUser   *proto.UserResponse
	Token         string
	Videos        []*proto.VideoMetadataResponse
	Notifications []Notification
	GRPCClient    proto.RepoServiceClient
}

type Notification struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

func NewAppState() *AppState {
	return &AppState{
		LoggedIn:      false,
		Videos:        make([]*proto.VideoMetadataResponse, 0),
		Notifications: make([]Notification, 0),
	}
}

func (s *AppState) SetUser(user *proto.UserResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.CurrentUser = user
	s.LoggedIn = true
}

func (s *AppState) GetUser() *proto.UserResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.CurrentUser
}

func (s *AppState) IsLoggedIn() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.LoggedIn
}

func (s *AppState) SetVideos(videos []*proto.VideoMetadataResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Videos = videos
}

func (s *AppState) GetVideos() []*proto.VideoMetadataResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Videos
}

func (s *AppState) AddNotification(notif Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Notifications = append(s.Notifications, notif)
	// Keep only last 50 notifications
	if len(s.Notifications) > 50 {
		s.Notifications = s.Notifications[1:]
	}
}

func (s *AppState) GetNotifications() []Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Notifications
}

func (s *AppState) SetGRPCClient(client proto.RepoServiceClient) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.GRPCClient = client
}

func (s *AppState) GetGRPCClient() proto.RepoServiceClient {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.GRPCClient
}

func (s *AppState) Logout() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.LoggedIn = false
	s.CurrentUser = nil
	s.Token = ""
	s.Videos = make([]*proto.VideoMetadataResponse, 0)
}
