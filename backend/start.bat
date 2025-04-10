@echo off
:: 自动请求管理员权限
>nul 2>&1 "%SYSTEMROOT%\system32\cacls.exe" "%SYSTEMROOT%\system32\config\system"
if '%errorlevel%' NEQ '0' (
  echo 请求管理员权限以配置防火墙规则...
  goto UACPrompt
) else ( goto main )

:UACPrompt
echo Set UAC = CreateObject^("Shell.Application"^) > "%temp%\getadmin.vbs"
echo UAC.ShellExecute "%~s0", "", "", "runas", 1 >> "%temp%\getadmin.vbs"
"%temp%\getadmin.vbs"
exit /B

:main
cd /d "%~dp0"

:: ============== 自定义配置 ==============
set PROJECT_NAME=ebidsystem
set EXE_NAME=%PROJECT_NAME%.exe
set PORT=3000
set ENV_FILE=.env
set BIN_DIR=bin

:: ============== 编译项目 ==============
echo 正在编译项目...
if not exist "%BIN_DIR%" mkdir "%BIN_DIR%"
go build -o "%BIN_DIR%\%EXE_NAME%" main.go
if %errorlevel% neq 0 (
  echo [错误] 编译失败，请检查以下可能原因：
  echo 1. Go 代码语法错误
  echo 2. 依赖未安装（运行 go mod tidy）
  pause
  exit /b 1
)

:: ============== 复制环境文件 ==============
if exist "%ENV_FILE%" (
  copy "%ENV_FILE%" "%BIN_DIR%\%ENV_FILE%" >nul
  echo 已复制 .env 文件到输出目录
) else (
  echo [警告] 未找到 .env 文件！
)

:: ============== 配置防火墙 ==============
echo 配置防火墙规则...
set EXE_FULL_PATH=%~dp0%BIN_DIR%\%EXE_NAME%

netsh advfirewall firewall delete rule name="Allow %PROJECT_NAME% Inbound" >nul 2>&1
netsh advfirewall firewall delete rule name="Allow %PROJECT_NAME% Outbound" >nul 2>&1

netsh advfirewall firewall add rule ^
  name="Allow %PROJECT_NAME% Inbound" ^
  dir=in ^
  program="%EXE_FULL_PATH%" ^
  action=allow ^
  enable=yes

netsh advfirewall firewall add rule ^
  name="Allow %PROJECT_NAME% Outbound" ^
  dir=out ^
  program="%EXE_FULL_PATH%" ^
  action=allow ^
  enable=yes

:: ============== 启动服务 ==============
echo 启动服务（端口 %PORT%）...
cd "%BIN_DIR%"
echo 服务运行日志：
echo -------------------------------
%EXE_NAME%
cd ..

echo -------------------------------
echo 服务已退出，按任意键关闭窗口...
pause >nul