# Telegram Bot for PC Control

This is a Telegram bot that allows you to control your PC remotely through Telegram commands.

## Features

- `/shutdown` - Shutdown the PC immediately
- `/restart` - Restart the PC immediately
- `/shutdown_timer` - Set a timer for shutdown (format: 1s, 1m, 1h)
- `/restart_timer` - Set a timer for restart (format: 1s, 1m, 1h)
- `/get_current_timer` - Get information about current timer
- `/cancel_timer` - Cancel any active timer
- `/capture_screen` - Capture and send screenshot (supports multiple displays)
- `/stats` - Get detailed system statistics:
  - CPU, RAM, Disk usage
  - GPU information
  - Top 5 CPU-consuming processes
  - Top 5 memory-consuming processes

## Setup

1. Create a new bot with [@BotFather](https://t.me/botfather) on Telegram
2. Get your bot token
3. Create a `.env` file in the project root with your bot token:
   ```
   TELEGRAM_BOT_TOKEN=your_bot_token_here
   ```
4. Install Go dependencies:
   ```bash
   go mod download
   ```
5. Build and run the bot:
   ```bash
   go build
   ./telegram-bot-control-pc
   ```

## Windows Service Installation

You can install the bot as a Windows service to run automatically at startup:

1. Build the application:
   ```bash
   go build
   ```

2. Run as administrator to register the service:
   ```bash
   telegram-bot-control-pc.exe -register-service
   ```

3. To unregister the service (run as administrator):
   ```bash
   telegram-bot-control-pc.exe -unregister-service
   ```

The bot will run automatically when Windows starts.

## Multi-Display Support

When using the `/capture_screen` command:
- If your PC has one display, it will capture and send the screenshot immediately
- If your PC has multiple displays, it will show selection buttons for each display with their dimensions
- Select the display you want to capture

## System Stats

The `/stats` command provides comprehensive system information:
- Basic system stats (CPU, RAM, disk usage)
- Top 5 processes consuming the most CPU
- Top 5 processes consuming the most memory
- GPU information (on Windows)

## Requirements

- Go 1.21 or higher
- Windows or Linux operating system
- Telegram Bot Token
- Required Go packages:
  - github.com/go-telegram-bot-api/telegram-bot-api/v5
  - github.com/joho/godotenv
  - github.com/kbinani/screenshot
  - github.com/shirou/gopsutil/v3

## Security Note

This bot has the ability to control your PC, so make sure to:
1. Only share the bot token with trusted users
2. Consider modifying the code to implement user authentication
3. Be careful with the commands you send
4. When running as a Windows service, it uses LocalSystem privileges

## License

MIT License 