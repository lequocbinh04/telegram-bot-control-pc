package main

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kbinani/screenshot"
)

func handleScreenCapture(message *tgbotapi.Message) {
	numDisplays := screenshot.NumActiveDisplays()
	if numDisplays == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ No displays found")
		bot.Send(msg)
		return
	}

	if numDisplays == 1 {
		// If only one display, capture it directly
		captureAndSendScreenshot(message, 0)
		return
	}

	// For multiple displays, show selection buttons
	displays := getDisplayInfo()
	msg := tgbotapi.NewMessage(message.Chat.ID, "🖥️ *Select display to capture:*\n\n"+strings.Join(displays, "\n"))
	msg.ParseMode = "MarkdownV2"
	
	// Create inline keyboard for display selection
	var buttons [][]tgbotapi.InlineKeyboardButton
	for i := 0; i < numDisplays; i++ {
		data := fmt.Sprintf("capture_display_%d", i)
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("Display %d", i+1),
				CallbackData: &data,
			},
		})
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	bot.Send(msg)
}

func handleDisplaySelection(callback *tgbotapi.CallbackQuery) {
	// Extract display number from callback data
	displayNum, err := strconv.Atoi(strings.Split(callback.Data, "_")[2])
	if err != nil {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "❌ Invalid display selection")
		bot.Send(msg)
		return
	}

	// Capture and send screenshot for selected display
	captureAndSendScreenshot(callback.Message, displayNum)
}

func captureAndSendScreenshot(message *tgbotapi.Message, displayNum int) {
	bounds := screenshot.GetDisplayBounds(displayNum)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Failed to capture screen")
		bot.Send(msg)
		return
	}

	// Create temp directory if it doesn't exist
	tempDir := "temp"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Failed to create temporary directory")
		bot.Send(msg)
		return
	}

	// Save the screenshot
	filename := filepath.Join(tempDir, fmt.Sprintf("screenshot_%d.png", displayNum+1))
	file, err := os.Create(filename)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Failed to create screenshot file")
		bot.Send(msg)
		return
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Failed to save screenshot")
		bot.Send(msg)
		return
	}

	// Send the screenshot
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("📸 Capturing display %d...", displayNum+1))
	bot.Send(msg)
	
	photo := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FilePath(filename))
	bot.Send(photo)

	// Clean up
	os.Remove(filename)
}

func getDisplayInfo() []string {
	var displays []string
	n := screenshot.NumActiveDisplays()
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		displays = append(displays, fmt.Sprintf("🖥️ Display %d: %dx%d", i+1, bounds.Dx(), bounds.Dy()))
	}
	return displays
} 