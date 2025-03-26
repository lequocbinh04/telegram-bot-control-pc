package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	bot *tgbotapi.BotAPI
	activeTimers = make(map[int64]*time.Timer)
	adminChatID int64 = 2081422788
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// Parse command line flags
	registerService := flag.Bool("register-service", false, "Register as Windows service")
	unregisterService := flag.Bool("unregister-service", false, "Unregister Windows service")
	flag.Parse()

	// Handle service registration/unregistration
	if *registerService {
		fmt.Println("Registering service...")
		if err := registerAsWindowsService(); err != nil {
			log.Fatalf("Failed to register service: %v", err)
		}
		fmt.Println("Service registered successfully")
		return
	}

	if *unregisterService {
		fmt.Println("Unregistering service...")
		if err := unregisterWindowsService(); err != nil {
			log.Fatalf("Failed to unregister service: %v", err)
		}
		fmt.Println("Service unregistered successfully")
		return
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set in .env file")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Set bot commands
	commands := []tgbotapi.BotCommand{
		{Command: "shutdown", Description: "⏹️ Shutdown PC"},
		{Command: "restart", Description: "🔄 Restart PC"},
		{Command: "shutdown_timer", Description: "⏰ Set shutdown timer"},
		{Command: "restart_timer", Description: "⏰ Set restart timer"},
		{Command: "stats", Description: "📊 Get system stats"},
		{Command: "capture_screen", Description: "📸 Take screenshot"},
		{Command: "cancel_timer", Description: "❌ Cancel active timer"},
		{Command: "get_current_timer", Description: "ℹ️ Check active timer"},
	}

	config := tgbotapi.SetMyCommandsConfig{
		Commands: commands,
	}
	_, err = bot.Request(config)
	if err != nil {
		log.Printf("Error setting commands: %v", err)
	}

	// Send startup notification with markdown formatting
	startupMsg := tgbotapi.NewMessage(adminChatID, `🚀 *PC Control Bot is now running\!*

*Available commands:*

🔄 /restart
_• Restart PC immediately_

⏹️ /shutdown
_• Shutdown PC immediately_

⏰ /shutdown\_timer
_• Set a timer to shutdown PC_
_• Format: 1s, 1m, 1h_

⏰ /restart\_timer
_• Set a timer to restart PC_
_• Format: 1s, 1m, 1h_

📊 /stats
_• Get system statistics_
_• CPU, RAM, GPU, Disk usage_

📸 /capture\_screen
_• Take a screenshot of your PC_

❌ /cancel\_timer
_• Cancel any active timer_

ℹ️ /get\_current\_timer
_• Check current active timer_`)
	startupMsg.ParseMode = "MarkdownV2"
	bot.Send(startupMsg)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		// Handle callback queries for display selection
		if update.CallbackQuery != nil {
			if strings.HasPrefix(update.CallbackQuery.Data, "capture_display_") {
				handleDisplaySelection(update.CallbackQuery)
			}
			continue
		}

		// Handle regular messages
		if !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = "MarkdownV2"

		switch update.Message.Command() {
		case "shutdown":
			msg.Text = "⏹️ *Shutting down PC\\.\\.\\.*"
			bot.Send(msg)
			go shutdownPC()
		case "restart":
			msg.Text = "🔄 *Restarting PC\\.\\.\\.*"
			bot.Send(msg)
			go restartPC()
		case "shutdown_timer":
			msg.Text = "⏰ *Set Shutdown Timer*\n\n_Please specify the time in format:_\n`1s` _for seconds_\n`1m` _for minutes_\n`1h` _for hours_"
			bot.Send(msg)
		case "restart_timer":
			msg.Text = "⏰ *Set Restart Timer*\n\n_Please specify the time in format:_\n`1s` _for seconds_\n`1m` _for minutes_\n`1h` _for hours_"
			bot.Send(msg)
		case "get_current_timer":
			msg.Text = getCurrentTimer(update.Message.From.ID)
			bot.Send(msg)
		case "cancel_timer":
			msg.Text = cancelTimer(update.Message.From.ID)
			bot.Send(msg)
		case "capture_screen":
			handleScreenCapture(update.Message)
		case "stats":
			msg.Text = getSystemStats()
			bot.Send(msg)
		}
	}
}

func getCurrentTimer(userID int64) string {
	timer, exists := activeTimers[userID]
	if !exists {
		return "ℹ️ *No active timer*"
	}
	return fmt.Sprintf("⏰ *Active timer:*\n`%v`", timer)
}

func cancelTimer(userID int64) string {
	timer, exists := activeTimers[userID]
	if !exists {
		return "❌ *No active timer to cancel*"
	}
	timer.Stop()
	delete(activeTimers, userID)
	return "✅ *Timer cancelled successfully*"
}

func parseTime(input string) (time.Duration, error) {
	input = strings.TrimSpace(input)
	value := input[:len(input)-1]
	unit := input[len(input)-1:]

	var duration time.Duration
	switch unit {
	case "s":
		duration = time.Second
	case "m":
		duration = time.Minute
	case "h":
		duration = time.Hour
	default:
		return 0, fmt.Errorf("invalid time unit: %s", unit)
	}

	var amount float64
	_, err := fmt.Sscanf(value, "%f", &amount)
	if err != nil {
		return 0, err
	}

	return time.Duration(amount) * duration, nil
}