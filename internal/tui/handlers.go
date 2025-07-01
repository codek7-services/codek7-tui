package tui

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/codek7-services/codek7-tui/internal"
	proto "github.com/codek7-services/codek7-tui/pkg/pb"
)

// Handler methods for Views
func (v *Views) handleLogin() {
	// We'll get form data through stored references
	v.processLogin("", "") // Placeholder - will be updated when called with actual data
}

func (v *Views) processLogin(username, password string) {
	if username == "" || password == "" {
		v.showMessage("Please fill in all fields")
		return
	}

	client := v.State.GetGRPCClient()
	if client == nil {
		v.showMessage("gRPC client not initialized")
		return
	}

	// Call GetUser (which acts like login in your current setup)
	user, err := client.GetUser(context.TODO(), &proto.GetUserRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		v.showError(err)
		return
	}

	v.State.SetUser(user)
	v.loadUserVideos() // Load videos after login
	v.ShowDashboardView()
	v.showMessage("Login successful!")
}

func (v *Views) handleRegister() {
	// We'll get form data through stored references
	v.processRegister("", "", "") // Placeholder
}

func (v *Views) processRegister(username, email, password string) {
	if username == "" || email == "" || password == "" {
		v.showMessage("Please fill in all fields")
		return
	}

	client := v.State.GetGRPCClient()
	if client == nil {
		v.showMessage("gRPC client not initialized")
		return
	}

	_, err := client.CreateUser(context.TODO(), &proto.CreateUserRequest{
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		v.showError(err)
		return
	}

	v.showMessage("Registration successful! Please login.")
	v.ShowLoginView()
}

func (v *Views) handleUpload() {
	// We'll get form data through stored references
	v.processUpload("", "", "") // Placeholder
}

func (v *Views) processUpload(filePath, title, description string) {
	if filePath == "" || title == "" {
		v.showMessage("❌ Please fill in required fields (File Path and Title)")
		return
	}

	// Check if file exists and get file info
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		v.showMessage("❌ File does not exist: " + filePath)
		return
	}
	if err != nil {
		v.showMessage("❌ Cannot access file: " + err.Error())
		return
	}

	// Check file size (limit to 500MB as mentioned in server code)
	maxSize := int64(500 * 1024 * 1024) // 500MB
	if fileInfo.Size() > maxSize {
		v.showMessage(fmt.Sprintf("❌ File too large: %.2f MB (max 500MB)", float64(fileInfo.Size())/(1024*1024)))
		return
	}

	user := v.State.GetUser()
	if user == nil {
		v.showMessage("❌ No user logged in")
		return
	}

	client := v.State.GetGRPCClient()
	if client == nil {
		v.showMessage("❌ gRPC client not initialized")
		return
	}

	// Show detailed progress message
	v.showMessage(fmt.Sprintf("📤 Uploading video...\n\n"+
		"📁 File: %s\n"+
		"📏 Size: %.2f MB\n"+
		"👤 User: %s\n\n"+
		"⏳ This may take a while depending on file size.\n"+
		"The file will be processed in real-time via Kafka streams.",
		fileInfo.Name(),
		float64(fileInfo.Size())/(1024*1024),
		user.Username))

	go func() {
		err := internal.UploadVideo(client, filePath, title, description, user.Id)
		v.App.QueueUpdateDraw(func() {
			if err != nil {
				v.showError(fmt.Errorf("Upload failed: %v", err))
			} else {
				// Success message with next steps
				v.showMessage("✅ Video uploaded successfully!\n\n" +
					"🎯 Your video is now being processed.\n" +
					"📡 You'll receive real-time notifications when ready.\n" +
					"🔄 The video list will be updated automatically.")

				// Add a notification about the upload
				v.State.AddNotification(Notification{
					ID:      fmt.Sprintf("upload-%d", time.Now().Unix()),
					Type:    "upload",
					Message: fmt.Sprintf("Video '%s' uploaded successfully", title),
					Time:    time.Now().Format("15:04:05"),
				})

				// Refresh data and return to dashboard
				v.loadUserVideos()

				// Auto-return to dashboard after showing success
				time.Sleep(2 * time.Second)
				v.App.QueueUpdateDraw(func() {
					v.Pages.RemovePage("message")
					v.ShowDashboardView()
				})
			}
		})
	}()
}

func (v *Views) handleLogout() {
	v.State.Logout()
	if v.WSManager != nil {
		v.WSManager.Disconnect()
	}
	v.Pages.SwitchToPage("main")
	v.showMessage("Logged out successfully!")
}

func (v *Views) loadUserVideos() {
	user := v.State.GetUser()
	if user == nil {
		return
	}

	client := v.State.GetGRPCClient()
	if client == nil {
		return
	}

	videos, err := client.GetUserVideos(context.TODO(), &proto.GetUserVideosRequest{
		UserId: user.Id,
	})
	if err != nil {
		// Don't show error popup, just log it
		return
	}

	v.State.SetVideos(videos.Videos)
}
