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

:: ============== 全局配置 ==============
set "PROJECT_NAME=ebidsystem"
set "EXE_NAME=%PROJECT_NAME%.exe"
set "PORT=3000"
set "ENV_FILE=.env"
set "ROOT_DIR=%~dp0"
set "BIN_DIR=%ROOT_DIR%bin"
set "LOG_DIR=%ROOT_DIR%bin\logs"
set "MATCH_LOG_DIR=%ROOT_DIR%bin\matchLog"

:: ============== 初始化目录结构 ==============
echo [%TIME%] 初始化目录结构...
if not exist "%BIN_DIR%" mkdir "%BIN_DIR%"
if not exist "%LOG_DIR%" mkdir "%LOG_DIR%"
if not exist "%MATCH_LOG_DIR%" mkdir "%MATCH_LOG_DIR%"
del /Q "%LOG_DIR%\*.log" "%MATCH_LOG_DIR%\*.txt" 2>nul

:: ============== 编译项目 ==============
echo [%TIME%] 正在编译项目...
go build -o "%BIN_DIR%\%EXE_NAME%" main.go
if %errorlevel% neq 0 (
    echo [%TIME%] [错误] 编译失败，请检查以下可能原因：
    echo 1. Go代码语法错误
    echo 2. 依赖未安装（运行 go mod tidy）
    pause
    exit /b 1
)

:: ============== 环境文件处理 ==============
if exist "%ENV_FILE%" (
    copy "%ENV_FILE%" "%BIN_DIR%\" >nul
    echo [%TIME%] 已复制环境文件到bin目录
) else (
    echo [%TIME%] [警告] 未找到.env文件！
)

:: ============== 防火墙配置 ==============
echo [%TIME%] 配置防火墙规则...
set "EXE_PATH=%BIN_DIR%\%EXE_NAME%"

netsh advfirewall firewall delete rule name="%PROJECT_NAME% Inbound" >nul 2>&1
netsh advfirewall firewall add rule ^
    name="%PROJECT_NAME% Inbound" ^
    dir=in ^
    action=allow ^
    program="%EXE_PATH%"

:: ============== 服务启动 ==============
echo [%TIME%] 启动服务（端口 %PORT%）...

:: 主服务窗口（HTTP日志）
start "EBidSystem Main Server" ^
    cmd /k "cd /d "%BIN_DIR%" && ^
    %EXE_NAME% > "%LOG_DIR%\service.log" 2>&1 & ^
    type "%LOG_DIR%\service.log""

:: 撮合日志窗口（独立）
set "MATCH_LOG_PATH=%MATCH_LOG_DIR%\matchLog.txt"
start "EBidSystem Matching Engine" ^
    cmd /k "cd /d "%BIN_DIR%" && ^
    powershell -NoExit -Command ^
    \"Get-Content -LiteralPath '%MATCH_LOG_PATH%' -Wait -Tail 10\""

:: 状态显示
echo -------------------------------
echo [服务状态]
echo 主服务窗口标题：EBidSystem Main Server
echo 撮合日志窗口标题：EBidSystem Matching Engine
echo 日志文件位置：
echo - 主日志：%LOG_DIR%\service.log
echo - 撮合日志：%MATCH_LOG_PATH%
echo -------------------------------
echo 按任意键关闭本窗口（服务将继续运行）...
pause >nul

exit /b 0