package tui

import (
	"fmt"

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

	form.AddInputField("File Path", "/path/to/video.mp4", 50, nil, func(text string) {
		filePath = text
	}).
		AddInputField("Title", "", 50, nil, func(text string) {
			title = text
		}).
		AddTextArea("Description", "", 50, 3, 0, func(text string) {
			description = text
		}).
		AddButton("Upload", func() {
			v.processUpload(filePath, title, description)
		}).
		AddButton("Back", func() {
			v.Pages.SwitchToPage("main")
		})

	form.SetBorder(true).SetTitle("📤 Upload Video").SetTitleAlign(tview.AlignCenter)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(form, 80, 1, true).
			AddItem(nil, 0, 1, false), 20, 1, true).
		AddItem(nil, 0, 1, false)

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

	table := tview.NewTable().SetBorders(true)

	// Headers with bold style
	table.SetCell(0, 0, tview.NewTableCell("ID").SetTextColor(tview.Styles.PrimaryTextColor).SetSelectable(false))
	table.SetCell(0, 1, tview.NewTableCell("Title").SetTextColor(tview.Styles.PrimaryTextColor).SetSelectable(false))
	table.SetCell(0, 2, tview.NewTableCell("Description").SetTextColor(tview.Styles.PrimaryTextColor).SetSelectable(false))
	table.SetCell(0, 3, tview.NewTableCell("Created").SetTextColor(tview.Styles.PrimaryTextColor).SetSelectable(false))
	table.SetCell(0, 4, tview.NewTableCell("File").SetTextColor(tview.Styles.PrimaryTextColor).SetSelectable(false))

	videos := v.State.GetVideos()
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

	table.Select(1, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			v.Pages.SwitchToPage("main")
		}
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(table, 0, 1, true).
		AddItem(tview.NewTextView().
			SetText("Press ESC to go back | Use arrow keys to navigate").
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
			SetText("Press ESC to go back").
			SetTextAlign(tview.AlignCenter), 1, 0, false)

	list.SetDoneFunc(func() {
		v.Pages.SwitchToPage("main")
	})

	v.Pages.AddAndSwitchToPage("notifications", flex, true)
}

// Dashboard View (for logged in users)
func (v *Views) ShowDashboardView() {
	if !v.State.IsLoggedIn() {
		v.Pages.SwitchToPage("main")
		return
	}

	user := v.State.GetUser()
	videos := v.State.GetVideos()
	notifications := v.State.GetNotifications()

	info := tview.NewTextView().
		SetText(fmt.Sprintf(
			"Welcome, %s!\n\nUser ID: %s\nVideos: %d\nNotifications: %d\n\nConnected to WebSocket: %v",
			user.Username,
			user.Id,
			len(videos),
			len(notifications),
			v.WSManager != nil && v.WSManager.IsConnected(),
		)).
		SetBorder(true).
		SetTitle("📊 Dashboard")

	menu := tview.NewList().
		AddItem("📤 Upload Video", "Upload a new video", 'u', v.ShowUploadView).
		AddItem("🎞️  My Videos", "View your videos", 'v', v.ShowVideosView).
		AddItem("📡 Notifications", "View notifications", 'n', v.ShowNotificationsView).
		AddItem("🔄 Refresh", "Refresh data", 'r', v.refreshData).
		AddItem("🔌 Start WebSocket", "Connect to notifications", 'w', v.startWebSocket).
		AddItem("🚪 Logout", "Logout and return to main menu", 'l', v.handleLogout).
		SetBorder(true).
		SetTitle("Actions")

	flex := tview.NewFlex().
		AddItem(info, 0, 1, false).
		AddItem(menu, 0, 1, true)

	v.Pages.AddAndSwitchToPage("dashboard", flex, true)
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

func (v *Views) refreshData() {
	if v.State.IsLoggedIn() {
		v.loadUserVideos()
		v.ShowDashboardView() // Refresh the dashboard
	}
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
