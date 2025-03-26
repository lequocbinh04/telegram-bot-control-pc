package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func registerAsWindowsService() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("service registration is only supported on Windows")
	}

	// Get the executable path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(exePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	// Create the service using sc command
	cmd := exec.Command("sc", "create", "TelegramBotControl", "binPath=", fmt.Sprintf("\"%s\"", absPath), "start=", "auto", "obj=", "LocalSystem")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create service: %v\nOutput: %s", err, string(output))
	}

	// Set the service description
	descCmd := exec.Command("sc", "description", "TelegramBotControl", "Telegram Bot for PC Control - Allows remote control of PC through Telegram")
	descOutput, err := descCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set service description: %v\nOutput: %s", err, string(descOutput))
	}

	return nil
}

func unregisterWindowsService() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("service unregistration is only supported on Windows")
	}

	// Stop the service first
	stopCmd := exec.Command("sc", "stop", "TelegramBotControl")
	stopOutput, err := stopCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop service: %v\nOutput: %s", err, string(stopOutput))
	}

	// Delete the service
	cmd := exec.Command("sc", "delete", "TelegramBotControl")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete service: %v\nOutput: %s", err, string(output))
	}

	return nil
} 