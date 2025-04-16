@echo off
:: ...（保持原有管理员权限检查逻辑）

:main
cd /d "%~dp0"

:: ============== 自定义配置 ==============
set PROJECT_NAME=ebidsystem
set EXE_NAME=%PROJECT_NAME%.exe
set PORT=3000
set ENV_FILE=.env
set BIN_DIR=bin
set LOG_DIR=bin\logs

:: ============== 初始化日志目录 ==============
if not exist "%LOG_DIR%" mkdir "%LOG_DIR%"
if not exist "%LOG_DIR%\matchLog" mkdir "%LOG_DIR%\matchLog"
del /Q "%LOG_DIR%\service.log" 2>nul
del /Q "%LOG_DIR%\error.log" 2>nul
del /Q "%LOG_DIR%\matchLog\match.log" 2>nul

:: ============== 编译 & 配置 ==============
echo [%TIME%] 正在编译项目...
go build -o "%BIN_DIR%\%EXE_NAME%" main.go || goto :error

if exist "%ENV_FILE%" (
    copy "%ENV_FILE%" "%BIN_DIR%\" >nul
    echo [%TIME%] 已复制环境文件
)

:: ============== 配置防火墙 ==============
echo [%TIME%] 配置防火墙...
set EXE_PATH=%~dp0%BIN_DIR%\%EXE_NAME%
netsh advfirewall firewall delete rule name="Allow %PROJECT_NAME% Inbound" >nul 2>&1
netsh advfirewall firewall add rule name="Allow %PROJECT_NAME% Inbound" dir=in program="%EXE_PATH%" action=allow

:: ============== 启动服务 ==============
echo [%TIME%] 启动服务（端口 %PORT%）...
cd "%BIN_DIR%"
echo -------------------------------
echo 实时日志（Ctrl+C退出）：

:: 方案1：使用PowerShell双输出
powershell -Command ".\%EXE_NAME% | Tee-Object -FilePath ..\%LOG_DIR%\service.log -Append"

:: 方案2：原生批处理双输出（选其一）
:: .\%EXE_NAME% > ..\%LOG_DIR%\service.log 2>&1 & type ..\%LOG_DIR%\service.log

echo -------------------------------
pause
exit /b 0

:error
echo [错误] 编译失败，请检查代码！
pause
exit /b 1