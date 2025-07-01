package tui

import (
	"log"
	"os"

	proto "github.com/codek7-services/codek7-tui/pkg/pb"
	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	App       *tview.Application
	Pages     *tview.Pages
	State     *AppState
	Views     *Views
	WSManager *WebSocketManager
}

func NewApp() *App {
	app := tview.NewApplication()
	pages := tview.NewPages()
	state := NewAppState()

	// Initialize gRPC client
	grpcAddr := os.Getenv("GRPC_ADDR")
	if grpcAddr == "" {
		grpcAddr = "localhost:50051" // default
	}

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to gRPC server: %v", err)
	} else {
		client := proto.NewRepoServiceClient(conn)
		state.SetGRPCClient(client)
	}

	views := NewViews(app, pages, state)

	// Create main menu
	mainMenu := tview.NewList().
		AddItem("🔑 Login", "Login to your account", 'l', views.ShowLoginView).
		AddItem("📝 Register", "Create a new account", 'r', views.ShowRegisterView).
		AddItem("📊 Dashboard", "View dashboard (login required)", 'd', func() {
			if state.IsLoggedIn() {
				views.ShowDashboardView()
			} else {
				views.showMessage("Please login first!")
			}
		}).
		AddItem("🎭 Demo Mode", "Try the app with demo data", 'm', views.EnableDemoMode).
		AddItem("❌ Quit", "Exit the application", 'q', func() {
			app.Stop()
		})

	mainMenu.SetBorder(true).SetTitle("📺 CodeK7 TUI - Main Menu").SetTitleAlign(tview.AlignCenter)

	// Create welcome text
	welcomeText := tview.NewTextView().
		SetText("Welcome to CodeK7 TUI!\n\nVideo Management System\n\n• Login or Register to get started\n• Upload and manage your videos\n• Real-time notifications via WebSocket\n• Browse your video library\n• Try Demo Mode for a quick preview").
		SetTextAlign(tview.AlignCenter).
		SetBorder(true).
		SetTitle("Welcome")

	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(welcomeText, 9, 0, false).
		AddItem(mainMenu, 0, 1, true)

	pages.AddPage("main", mainFlex, true, true)

	tuiApp := &App{
		App:   app,
		Pages: pages,
		State: state,
		Views: views,
	}

	// Set the app root
	app.SetRoot(pages, true)

	return tuiApp
}

func (a *App) Run() error {
	return a.App.Run()
}
