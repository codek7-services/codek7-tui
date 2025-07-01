package tui

import (
	"context"
	"os"

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
		v.showMessage("Please fill in required fields (File Path and Title)")
		return
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		v.showMessage("File does not exist: " + filePath)
		return
	}

	user := v.State.GetUser()
	if user == nil {
		v.showMessage("No user logged in")
		return
	}

	client := v.State.GetGRPCClient()
	if client == nil {
		v.showMessage("gRPC client not initialized")
		return
	}

	// Show progress message
	v.showMessage("Uploading video... This may take a while.")

	go func() {
		err := internal.UploadVideo(client, filePath, title, description, user.Id)
		v.App.QueueUpdateDraw(func() {
			if err != nil {
				v.showError(err)
			} else {
				v.showMessage("Video uploaded successfully!")
				v.loadUserVideos() // Refresh video list
				v.ShowDashboardView()
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
