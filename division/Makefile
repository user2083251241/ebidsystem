# Makefile for Windows
# 使用方式：以管理员身份运行 PowerShell，执行 `make all`
# 依赖：Go 环境、GNU Make、Netsh 命令

SHELL := cmd.exe

# 定义常量
PROJECT_NAME := ebidsystem
EXE_NAME := $(PROJECT_NAME).exe
BIN_DIR := bin
LOG_DIR := $(BIN_DIR)\logs
PORT := 3000
EXE_PATH := $(CURDIR)\$(BIN_DIR)\$(EXE_NAME)  # 绝对路径

# 默认目标：一键启动（编译 + 配置防火墙 + 运行）
.PHONY: all
all: clean build setup-firewall run

# 清理旧构建产物和日志
.PHONY: clean
clean:
	echo [清理] 删除旧构建文件和日志...
	if exist "$(BIN_DIR)\$(EXE_NAME)" del /Q "$(BIN_DIR)\$(EXE_NAME)"
	if exist "$(LOG_DIR)" rmdir /S /Q "$(LOG_DIR)"
	echo.

# 编译项目（显示详细错误）
.PHONY: build
build:
	echo [编译] 正在生成可执行文件...
	go build -o "$(BIN_DIR)\$(EXE_NAME)" main.go
	if not exist "$(BIN_DIR)\$(EXE_NAME)" (
	    echo [错误] 编译失败！
	    exit 1
	)
	echo.

# 配置防火墙规则（需管理员权限）
.PHONY: setup-firewall
setup-firewall:
	echo [防火墙] 配置入站规则...
	netsh advfirewall firewall delete rule name="Allow $(PROJECT_NAME) Inbound" >nul 2>&1
	netsh advfirewall firewall add rule name="Allow $(PROJECT_NAME) Inbound" dir=in program="$(EXE_PATH)" action=allow
	echo.

# 启动服务并输出日志（强制创建日志目录）
.PHONY: run
run:
	echo [启动] 服务运行中（端口 $(PORT)）...
	mkdir "$(LOG_DIR)" 2>nul || echo [警告] 日志目录已存在
	echo -------------------------------
	cd "$(BIN_DIR)" && (
	    powershell -Command ".\$(EXE_NAME) | Tee-Object -FilePath '..\$(LOG_DIR)\service.log' -Append"
	)
	echo -------------------------------
	pause