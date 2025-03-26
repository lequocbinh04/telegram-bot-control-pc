package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"sort"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// escapeMarkdown escapes special characters for MarkdownV2
func escapeMarkdown(text string) string {
	special := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range special {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}
	return text
}

// ProcessInfo holds information about a process
type ProcessInfo struct {
	PID       int32
	Name      string
	CPUUsage  float64
	MemoryMB  float64
	MemoryPct float32
}

func getSystemStats() string {
	var stats []string

	// CPU Info
	cpuPercent, _ := cpu.Percent(0, false)
	stats = append(stats, fmt.Sprintf("*CPU Usage:* `%.2f%%`", cpuPercent[0]))

	// Memory Info
	memInfo, _ := mem.VirtualMemory()
	stats = append(stats, fmt.Sprintf("*Memory Usage:* `%.2f%%`", memInfo.UsedPercent))
	stats = append(stats, fmt.Sprintf("*Total Memory:* `%.2f GB`", float64(memInfo.Total)/1024/1024/1024))
	stats = append(stats, fmt.Sprintf("*Used Memory:* `%.2f GB`", float64(memInfo.Used)/1024/1024/1024))

	// Disk Info
	diskInfo, _ := disk.Usage("/")
	stats = append(stats, fmt.Sprintf("*Disk Usage:* `%.2f%%`", diskInfo.UsedPercent))
	stats = append(stats, fmt.Sprintf("*Total Disk:* `%.2f GB`", float64(diskInfo.Total)/1024/1024/1024))
	stats = append(stats, fmt.Sprintf("*Used Disk:* `%.2f GB`", float64(diskInfo.Used)/1024/1024/1024))

	// Host Info
	hostInfo, _ := host.Info()
	platform := escapeMarkdown(hostInfo.Platform)
	hostname := escapeMarkdown(hostInfo.Hostname)
	stats = append(stats, fmt.Sprintf("*OS:* `%s`", platform))
	stats = append(stats, fmt.Sprintf("*Hostname:* `%s`", hostname))

	// GPU Info (Windows specific)
	if runtime.GOOS == "windows" {
		gpuInfo := getGPUInfo()
		if gpuInfo != "" {
			stats = append(stats, fmt.Sprintf("*GPU:* `%s`", escapeMarkdown(gpuInfo)))
		}
	}

	// Get top processes
	topCPUProcesses, topMemProcesses := getTopProcesses(5)

	// Add top CPU processes
	stats = append(stats, "\n*Top 5 CPU Processes:*")
	for i, p := range topCPUProcesses {
		stats = append(stats, fmt.Sprintf("`%d` %s: `%.2f%%`", i+1, escapeMarkdown(p.Name), p.CPUUsage))
	}

	// Add top memory processes
	stats = append(stats, "\n*Top 5 Memory Processes:*")
	for i, p := range topMemProcesses {
		stats = append(stats, fmt.Sprintf("`%d` %s: `%.2f MB` \\(`%.2f%%`\\)", i+1, escapeMarkdown(p.Name), p.MemoryMB, p.MemoryPct))
	}

	return strings.Join(stats, "\n")
}

func getGPUInfo() string {
	cmd := exec.Command("wmic", "path", "win32_VideoController", "get", "name")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) > 1 {
		return strings.TrimSpace(lines[1])
	}
	return ""
}

// getTopProcesses returns the top N processes by CPU and memory usage
func getTopProcesses(n int) ([]ProcessInfo, []ProcessInfo) {
	var cpuProcesses []ProcessInfo
	var memProcesses []ProcessInfo
	
	// Get total memory
	memInfo, _ := mem.VirtualMemory()
	totalMem := float64(memInfo.Total)

	// Get all processes
	processes, _ := process.Processes()
	
	// Collect process information
	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}
		
		cpuPercent, err := p.CPUPercent()
		if err != nil {
			cpuPercent = 0
		}
		
		memInfo, err := p.MemoryInfo()
		if err != nil || memInfo == nil {
			continue
		}
		
		memBytes := float64(memInfo.RSS)
		memMB := memBytes / 1024 / 1024
		memPercent := float32(memBytes / totalMem * 100)
		
		processInfo := ProcessInfo{
			PID:       p.Pid,
			Name:      name,
			CPUUsage:  cpuPercent,
			MemoryMB:  memMB,
			MemoryPct: memPercent,
		}
		
		cpuProcesses = append(cpuProcesses, processInfo)
		memProcesses = append(memProcesses, processInfo)
	}
	
	// Sort by CPU usage (descending)
	sort.Slice(cpuProcesses, func(i, j int) bool {
		return cpuProcesses[i].CPUUsage > cpuProcesses[j].CPUUsage
	})
	
	// Sort by memory usage (descending)
	sort.Slice(memProcesses, func(i, j int) bool {
		return memProcesses[i].MemoryMB > memProcesses[j].MemoryMB
	})
	
	// Limit to top N
	if len(cpuProcesses) > n {
		cpuProcesses = cpuProcesses[:n]
	}
	if len(memProcesses) > n {
		memProcesses = memProcesses[:n]
	}
	
	return cpuProcesses, memProcesses
} 