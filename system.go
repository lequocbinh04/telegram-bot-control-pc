package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

func shutdownPC() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("shutdown", "/s", "/t", "0")
	} else {
		cmd = exec.Command("shutdown", "-h", "now")
	}
	cmd.Run()
}

func restartPC() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("shutdown", "/r", "/t", "0")
	} else {
		cmd = exec.Command("shutdown", "-r", "now")
	}
	cmd.Run()
}

func shutdownWithTimer(duration time.Duration) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("shutdown", "/s", "/t", fmt.Sprintf("%d", int(duration.Seconds())))
	} else {
		cmd = exec.Command("shutdown", "-h", fmt.Sprintf("+%d", int(duration.Minutes())))
	}
	cmd.Run()
}

func restartWithTimer(duration time.Duration) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("shutdown", "/r", "/t", fmt.Sprintf("%d", int(duration.Seconds())))
	} else {
		cmd = exec.Command("shutdown", "-r", fmt.Sprintf("+%d", int(duration.Minutes())))
	}
	cmd.Run()
} 