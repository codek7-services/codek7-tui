package tui

import (
	"context"
	"fmt"
	"time"

	proto "github.com/codek7-services/codek7-tui/pkg/pb"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Views struct {
	App       *tview.Application
	Pages     *tview.Pages
	State     *AppState
	WSManager *WebSocketManager
}

func NewViews(app *tview.Application, pages *tview.Pages, state *AppState) *Views {
	return &Views{
		App:   app,
		Pages: pages,
		State: state,
	}
}

func (v *Views) SetWebSocketManager(wsm *WebSocketManager) {
	v.WSManager = wsm
}

// Login View
func (v *Views) ShowLoginView() {
	form := tview.NewForm()

	var username, password string

	form.AddInputField("Username", "", 30, nil, func(text string) {
		username = text
	}).
		AddPasswordField("Password", "", 30, '*', func(text string) {
			password = text
		}).
		AddButton("Login", func() {
			v.processLogin(username, password)
		}).
		AddButton("Register", v.showRegisterView).
		AddButton("Back", func() {
			v.Pages.SwitchToPage("main")
		})

	form.SetBorder(true).SetTitle("🔑 Login").SetTitleAlign(tview.AlignCenter)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(form, 60, 1, true).
			AddItem(nil, 0, 1, false), 15, 1, true).
		AddItem(nil, 0, 1, false)

	v.Pages.AddAndSwitchToPage("login", flex, true)
}

// Register View
func (v *Views) ShowRegisterView() {
	form := tview.NewForm()

	var username, email, password string

	form.AddInputField("Username", "", 30, nil, func(text string) {
		username = text
	}).
		AddInputField("Email", "", 30, nil, func(text string) {
			email = text
		}).
		AddPasswordField("Password", "", 30, '*', func(text string) {
			password = text
		}).
		AddButton("Register", func() {
			v.processRegister(username, email, password)
		}).
		AddButton("Back to Login", v.ShowLoginView).
		AddButton("Back", func() {
			v.Pages.SwitchToPage("main")
		})

	form.SetBorder(true).SetTitle("📝 Register").SetTitleAlign(tview.AlignCenter)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(form, 60, 1, true).
			AddItem(nil, 0, 1, false), 17, 1, true).
		AddItem(nil, 0, 1, false)

	v.Pages.AddAndSwitchToPage("register", flex, true)
}

func (v *Views) showRegisterView() {
	v.ShowRegisterView()
}

// Upload View
func (v *Views) ShowUploadView() {
	if !v.State.IsLoggedIn() {
		v.showMessage("Please login first!")
		return
	}

	form := tview.NewForm()

	var filePath, title, description string

	form.AddInputField("File Path", "/home/user/example.mp4", 70, nil, func(text string) {
		filePath = text
	}).
		AddInputField("Title", "", 50, nil, func(text string) {
			title = text
		}).
		AddTextArea("Description", "", 50, 4, 0, func(text string) {
			description = text
		}).
		AddButton("📤 Upload Video", func() {
			v.processUpload(filePath, title, description)
		}).
		AddButton("📁 Browse Files", func() {
			v.showFileBrowser(func(selectedPath string) {
				// Update the file path field when a file is selected
				filePath = selectedPath
				// Refresh the form to show updated path
				form.GetFormItem(0).(*tview.InputField).SetText(selectedPath)
			})
		}).
		AddButton("🏠 Back to Dashboard", func() {
			v.ShowDashboardView()
		})

	form.SetBorder(true).SetTitle("📤 Upload Video - Real-time Processing").SetTitleAlign(tview.AlignCenter)

	// Add help text
	helpText := tview.NewTextView().
		SetText("📋 Upload Instructions:\n" +
			"• Supported formats: MP4, AVI, MOV, MKV\n" +
			"• Maximum file size: 500MB\n" +
			"• Processing happens in real-time via Kafka\n" +
			"• You'll receive notifications when complete\n" +
			"• Files are chunked for efficient streaming").
		SetBorder(true).
		SetTitle("ℹ️ Help")

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(helpText, 8, 0, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(form, 80, 1, true).
			AddItem(nil, 0, 1, false), 0, 1, true)

	v.Pages.AddAndSwitchToPage("upload", flex, true)
}

// Videos View
func (v *Views) ShowVideosView() {
	if !v.State.IsLoggedIn() {
		v.showMessage("Please login first!")
		return
	}

	// Load videos first
	v.loadUserVideos()

	table := tview.NewTable().SetBorders(true).SetSelectable(true, false)

	// Headers with bold style
	table.SetCell(0, 0, tview.NewTableCell("ID").SetTextColor(tcell.ColorYellow).SetSelectable(false).SetAlign(tview.AlignCenter))
	table.SetCell(0, 1, tview.NewTableCell("Title").SetTextColor(tcell.ColorYellow).SetSelectable(false).SetAlign(tview.AlignCenter))
	table.SetCell(0, 2, tview.NewTableCell("Description").SetTextColor(tcell.ColorYellow).SetSelectable(false).SetAlign(tview.AlignCenter))
	table.SetCell(0, 3, tview.NewTableCell("Created").SetTextColor(tcell.ColorYellow).SetSelectable(false).SetAlign(tview.AlignCenter))
	table.SetCell(0, 4, tview.NewTableCell("File").SetTextColor(tcell.ColorYellow).SetSelectable(false).SetAlign(tview.AlignCenter))

	videos := v.State.GetVideos()
	if len(videos) == 0 {
		// Show empty state message
		table.SetCell(1, 0, tview.NewTableCell("No videos found"))
		table.SetCell(1, 1, tview.NewTableCell("Upload your first video!"))
		table.SetCell(1, 2, tview.NewTableCell(""))
		table.SetCell(1, 3, tview.NewTableCell(""))
		table.SetCell(1, 4, tview.NewTableCell(""))
	} else {
		for i, video := range videos {
			row := i + 1
			table.SetCell(row, 0, tview.NewTableCell(video.Id))
			table.SetCell(row, 1, tview.NewTableCell(video.Title))

			desc := video.Description
			if len(desc) > 30 {
				desc = desc[:30] + "..."
			}
			table.SetCell(row, 2, tview.NewTableCell(desc))
			table.SetCell(row, 3, tview.NewTableCell(video.CreatedAt))
			table.SetCell(row, 4, tview.NewTableCell(video.FileName))
		}
	}

	table.Select(1, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			v.ShowDashboardView()
		}
	})

	table.SetBorder(true).SetTitle("🎞️  My Videos").SetTitleAlign(tview.AlignCenter)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(table, 0, 1, true).
		AddItem(tview.NewTextView().
			SetText("Press ESC to go back to Dashboard | Use arrow keys to navigate").
			SetTextAlign(tview.AlignCenter), 1, 0, false)

	v.Pages.AddAndSwitchToPage("videos", flex, true)
}

// Notifications View
func (v *Views) ShowNotificationsView() {
	if !v.State.IsLoggedIn() {
		v.showMessage("Please login first!")
		return
	}

	list := tview.NewList()
	notifications := v.State.GetNotifications()

	if len(notifications) == 0 {
		list.AddItem("No notifications", "Connect to WebSocket to receive notifications", 0, nil)
	} else {
		for _, notif := range notifications {
			title := fmt.Sprintf("[%s] %s", notif.Type, notif.Time)
			list.AddItem(title, notif.Message, 0, nil)
		}
	}

	list.SetBorder(true).SetTitle("📡 Notifications").SetTitleAlign(tview.AlignCenter)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(list, 0, 1, true).
		AddItem(tview.NewTextView().
			SetText("Press ESC to go back to Dashboard").
			SetTextAlign(tview.AlignCenter), 1, 0, false)

	list.SetDoneFunc(func() {
		v.ShowDashboardView()
	})

	v.Pages.AddAndSwitchToPage("notifications", flex, true)
}

// Dashboard View (for logged in users)
func (v *Views) ShowDashboardView() {
	if !v.State.IsLoggedIn() {
		v.Pages.SwitchToPage("main")
		return
	}

	// Auto-load user videos when showing dashboard
	v.loadUserVideos()

	user := v.State.GetUser()
	videos := v.State.GetVideos()
	notifications := v.State.GetNotifications()

	// Handle case where user might be nil
	username := "Demo User"
	userID := "demo-123"
	if user != nil {
		username = user.Username
		userID = user.Id
	}

	// Enhanced info panel with recent activity
	recentVideoText := "No videos yet"
	if len(videos) > 0 {
		recentVideoText = fmt.Sprintf("Latest: %s", videos[0].Title)
		if len(recentVideoText) > 30 {
			recentVideoText = recentVideoText[:30] + "..."
		}
	}

	wsStatus := "❌ Disconnected"
	if v.WSManager != nil && v.WSManager.IsConnected() {
		wsStatus = "✅ Connected"
	}

	// Create info panel
	infoText := fmt.Sprintf(
		"🎉 Welcome back, %s!\n\n"+
			"👤 User ID: %s\n"+
			"🎥 Total Videos: %d\n"+
			"📺 %s\n"+
			"📡 Notifications: %d\n"+
			"🔌 WebSocket: %s\n\n"+
			"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"+
			"💡 Navigation Tips:\n"+
			"   • Use arrow keys or hotkeys\n"+
			"   • Press TAB to switch focus\n"+
			"   • ESC to go back from any view\n"+
			"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━",
		username,
		userID,
		len(videos),
		recentVideoText,
		len(notifications),
		wsStatus)

	info := tview.NewTextView().
		SetText(infoText).
		SetBorder(true).
		SetTitle("📊 Dashboard - Real-time Status")

	menu := tview.NewList().
		AddItem("📤 Upload Video", "Upload a new video file", 'u', v.ShowUploadView).
		AddItem("🎞️  My Videos", "Browse and manage your videos", 'v', v.ShowVideosView).
		AddItem("📡 Notifications", "View real-time notifications", 'n', v.ShowNotificationsView).
		AddItem("📊 Recent Videos", "View your 3 most recent videos", 's', v.ShowRecentVideosView).
		AddItem("🔄 Refresh Data", "Reload videos and notifications", 'r', v.refreshData).
		AddItem("🔌 WebSocket", "Toggle real-time connection", 'w', v.toggleWebSocket).
		AddItem("🏠 Main Menu", "Return to main menu", 'm', func() {
			v.Pages.SwitchToPage("main")
		}).
		AddItem("🚪 Logout", "End session and logout", 'l', v.handleLogout).
		SetBorder(true).
		SetTitle("🎯 Quick Actions")

	// Create main layout
	flex := tview.NewFlex().
		AddItem(info, 0, 2, false).
		AddItem(menu, 0, 1, true)

	// Enhanced keyboard shortcuts
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'u':
			v.ShowUploadView()
			return nil
		case 'v':
			v.ShowVideosView()
			return nil
		case 'n':
			v.ShowNotificationsView()
			return nil
		case 's':
			v.ShowRecentVideosView()
			return nil
		case 'r':
			v.refreshData()
			return nil
		case 'w':
			v.toggleWebSocket()
			return nil
		case 'm':
			v.Pages.SwitchToPage("main")
			return nil
		case 'l':
			v.handleLogout()
			return nil
		case 'q':
			v.App.Stop()
			return nil
		}
		return event
	})

	v.Pages.AddAndSwitchToPage("dashboard", flex, true)
	v.App.SetFocus(menu)
}

// Recent Videos View - shows last 3 videos
func (v *Views) ShowRecentVideosView() {
	if !v.State.IsLoggedIn() {
		v.showMessage("Please login first!")
		return
	}

	user := v.State.GetUser()
	client := v.State.GetGRPCClient()
	if client == nil {
		v.showMessage("gRPC client not available")
		return
	}

	// Try to get recent videos via gRPC
	recentVideos, err := client.GetLast3UserVideos(context.TODO(), &proto.GetLast3UserVideosRequest{
		UserId: user.Id,
	})

	var videos []*proto.VideoMetadataResponse
	if err != nil {
		// Fallback to local state if gRPC fails
		allVideos := v.State.GetVideos()
		if len(allVideos) > 3 {
			videos = allVideos[:3]
		} else {
			videos = allVideos
		}
	} else {
		videos = recentVideos.Videos
	}

	table := tview.NewTable().SetBorders(true).SetSelectable(true, false)

	// Headers
	table.SetCell(0, 0, tview.NewTableCell("Title").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	table.SetCell(0, 1, tview.NewTableCell("Description").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	table.SetCell(0, 2, tview.NewTableCell("Created").SetTextColor(tcell.ColorYellow).SetSelectable(false))

	if len(videos) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No recent videos"))
		table.SetCell(1, 1, tview.NewTableCell("Upload your first video!"))
		table.SetCell(1, 2, tview.NewTableCell(""))
	} else {
		for i, video := range videos {
			row := i + 1
			table.SetCell(row, 0, tview.NewTableCell(video.Title))

			desc := video.Description
			if len(desc) > 40 {
				desc = desc[:40] + "..."
			}
			table.SetCell(row, 1, tview.NewTableCell(desc))
			table.SetCell(row, 2, tview.NewTableCell(video.CreatedAt))
		}
	}

	table.Select(1, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			v.ShowDashboardView()
		}
	})

	table.SetBorder(true).SetTitle("📊 Recent Videos (Last 3)").SetTitleAlign(tview.AlignCenter)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(table, 0, 1, true).
		AddItem(tview.NewTextView().
			SetText("Press ESC to return to Dashboard | Use arrows to navigate").
			SetTextAlign(tview.AlignCenter), 1, 0, false)

	v.Pages.AddAndSwitchToPage("recent", flex, true)
}

// Toggle WebSocket connection
func (v *Views) toggleWebSocket() {
	if v.WSManager == nil {
		v.showMessage("WebSocket manager not initialized!")
		return
	}

	if v.WSManager.IsConnected() {
		v.WSManager.Disconnect()
		v.showMessage("🔌 WebSocket disconnected")
	} else {
		v.startWebSocket()
	}
}

// File browser for selecting files
func (v *Views) showFileBrowser(onSelect func(string)) {
	// Simple file browser implementation
	// In a real implementation, you might want to use a more sophisticated file picker
	modal := tview.NewModal().
		SetText("File Browser\n\nEnter the full path to your video file:\n\nExamples:\n/home/user/videos/my_video.mp4\n/tmp/upload.mov\n./local_video.avi").
		AddButtons([]string{"📁 Select File", "❌ Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			v.Pages.RemovePage("file_browser")
			if buttonIndex == 0 && buttonLabel == "📁 Select File" {
				// Create a simple input dialog for file path
				v.showFilePathDialog(onSelect)
			}
		})
	v.Pages.AddPage("file_browser", modal, false, true)
}

// Simple file path input dialog
func (v *Views) showFilePathDialog(onSelect func(string)) {
	form := tview.NewForm()
	var selectedPath string

	form.AddInputField("File Path", "", 60, nil, func(text string) {
		selectedPath = text
	}).
		AddButton("✅ Select", func() {
			if selectedPath != "" {
				onSelect(selectedPath)
			}
			v.Pages.RemovePage("file_dialog")
		}).
		AddButton("❌ Cancel", func() {
			v.Pages.RemovePage("file_dialog")
		})

	form.SetBorder(true).SetTitle("📁 Select Video File")

	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(form, 7, 0, true).
			AddItem(nil, 0, 1, false), 60, 0, true).
		AddItem(nil, 0, 1, false)

	v.Pages.AddPage("file_dialog", modal, false, true)
}

// Enhanced refresh with real-time updates
func (v *Views) refreshData() {
	if !v.State.IsLoggedIn() {
		v.showMessage("Please login first!")
		return
	}

	// Show loading message
	v.showMessage("🔄 Refreshing data...")

	go func() {
		// Load user videos
		v.loadUserVideos()

		// Try to fetch recent videos
		user := v.State.GetUser()
		client := v.State.GetGRPCClient()

		if client != nil && user != nil {
			// Fetch recent videos asynchronously
			recentVideos, err := client.GetLast3UserVideos(context.TODO(), &proto.GetLast3UserVideosRequest{
				UserId: user.Id,
			})

			if err == nil && len(recentVideos.Videos) > 0 {
				// Update local state with fresh data
				v.State.SetVideos(recentVideos.Videos)
			}
		}

		// Update UI on main thread
		v.App.QueueUpdateDraw(func() {
			v.Pages.RemovePage("message")
			v.showMessage("✅ Data refreshed! Updated videos and notifications.")

			// Auto-refresh dashboard after a short delay
			time.Sleep(1 * time.Second)
			v.App.QueueUpdateDraw(func() {
				v.Pages.RemovePage("message")
				v.ShowDashboardView()
			})
		})
	}()
}

// Helper methods
func (v *Views) showMessage(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			v.Pages.RemovePage("message")
		})
	v.Pages.AddPage("message", modal, false, true)
}

func (v *Views) showError(err error) {
	v.showMessage(fmt.Sprintf("Error: %v", err))
}

func (v *Views) startWebSocket() {
	if !v.State.IsLoggedIn() {
		v.showMessage("Please login first!")
		return
	}

	if v.WSManager != nil && v.WSManager.IsConnected() {
		v.showMessage("WebSocket already connected!")
		return
	}

	user := v.State.GetUser()
	if user == nil {
		v.showMessage("No user information available!")
		return
	}

	if v.WSManager == nil {
		v.WSManager = NewWebSocketManager(v.State, v.App)
	}

	go v.WSManager.Connect(user.Id)
	v.showMessage("Connecting to WebSocket...")
}
