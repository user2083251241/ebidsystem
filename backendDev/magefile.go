//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	ProjectName = "ebidsystem"
	ExeName     = ProjectName + ".exe"
	BinDir      = "bin"
	Port        = 3000
)

// 全局变量改为函数内定义
var (
	LogDir  = filepath.Join(BinDir, "logs") // 移到全局变量区，使用 var 声明
	ExePath = filepath.Join(BinDir, ExeName)
)

// 默认目标：一键启动（清理 -> 编译 -> 配置防火墙 -> 运行）
func All() error {
	Clean()
	if err := Build(); err != nil {
		return err
	}
	if err := SetupFirewall(); err != nil {
		return err
	}
	return Run()
}

// 清理旧构建文件和日志
func Clean() {
	fmt.Println("[清理] 删除旧构建文件和日志...")
	os.RemoveAll(BinDir)
}

// 编译项目
func Build() error {
	fmt.Println("[编译] 生成可执行文件...")
	cmd := exec.Command("go", "build", "-o", ExePath, "main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// 配置防火墙规则（Windows 专用）
func SetupFirewall() error {
	if runtime.GOOS != "windows" {
		fmt.Println("[跳过] 非 Windows 系统无需配置防火墙")
		return nil
	}

	fmt.Println("[防火墙] 配置入站规则...")
	exePath, _ := filepath.Abs(ExePath)

	// 使用 gsudo 提权执行 netsh 命令
	cmdDelete := exec.Command("gsudo", "netsh", "advfirewall", "firewall", "delete", "rule",
		"name=Allow "+ProjectName+" Inbound")
	cmdDelete.Stdout = os.Stdout
	cmdDelete.Stderr = os.Stderr
	_ = cmdDelete.Run()

	cmdAdd := exec.Command("gsudo", "netsh", "advfirewall", "firewall", "add", "rule",
		"name=Allow "+ProjectName+" Inbound",
		"dir=in",
		"program="+exePath,
		"action=allow")
	cmdAdd.Stdout = os.Stdout
	cmdAdd.Stderr = os.Stderr
	return cmdAdd.Run()
}

// 启动服务并记录日志
func Run() error {
	fmt.Printf("[启动] 服务运行中（端口 %d）...\n", Port)

	// 创建日志目录
	if err := os.MkdirAll(LogDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 启动服务并输出日志
	logPath := filepath.Join(LogDir, "service.log")
	logFile, _ := os.Create(logPath)
	defer logFile.Close()

	cmd := exec.Command(ExePath)
	cmd.Stdout = logFile   // 日志写入文件
	cmd.Stderr = os.Stderr // 错误输出到控制台

	// 启动服务
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动服务失败: %v", err)
	}

	fmt.Println("服务已启动，日志输出到:", logPath)
	fmt.Println("按 Ctrl+C 停止服务...")
	cmd.Wait() // 等待服务退出
	return nil
}
