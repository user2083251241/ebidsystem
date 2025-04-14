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
:: ============== 基础配置 ==============
set "PROJECT_ROOT=%~dp0"
set "PROJECT_NAME=ebidsystem"
set "EXE_NAME=%PROJECT_NAME%.exe"
set "BIN_DIR=%PROJECT_ROOT%bin"
set "LOG_DIR=%PROJECT_ROOT%bin\logs"
set "MATCH_LOG_DIR=%PROJECT_ROOT%bin\matchLog"

:: ============== 初始化目录 ==============
echo [%TIME%] 初始化目录结构...
mkdir "%BIN_DIR%" 2>nul
mkdir "%LOG_DIR%" 2>nul
mkdir "%MATCH_LOG_DIR%" 2>nul
del "%LOG_DIR%\*.log" "%MATCH_LOG_DIR%\*.txt" 2>nul

:: ============== 编译检查 ==============
echo [%TIME%] 编译项目...
cd /d "%PROJECT_ROOT%"
go build -o "%BIN_DIR%\%EXE_NAME%" main.go || goto build_failed

:: ============== 服务启动 ==============
echo [%TIME%] 启动服务...

:: 主服务窗口（带日志重定向）
start "EBidSystem_HTTP_Service" cmd /c ^
    "cd /d "%BIN_DIR%" && ^
    %EXE_NAME% > "%LOG_DIR%\service.log" 2>&1 & ^
    type "%LOG_DIR%\service.log" & ^
    pause"

:: 撮合日志窗口（带路径验证）
if exist "%MATCH_LOG_DIR%\matchLog.txt" (
    start "EBidSystem_Matching_Log" cmd /k ^
        "cd /d "%MATCH_LOG_DIR%" && ^
        echo 正在监视撮合日志... && ^
        powershell -NoExit -Command ^
        \"Get-Content -LiteralPath '%MATCH_LOG_DIR%\matchLog.txt' -Wait -Tail 10\""
) else (
    start "EBidSystem_Matching_Log" cmd /k ^
        "cd /d "%MATCH_LOG_DIR%" && ^
        echo 等待日志文件生成... && ^
        echo 预期路径: %MATCH_LOG_DIR%\matchLog.txt && ^
        timeout /t 30"
)

:: 状态提示
echo -------------------------------
echo 服务启动完成
echo 主服务窗口: EBidSystem_HTTP_Service
echo 撮合日志窗口: EBidSystem_Matching_Log
echo -------------------------------
pause
exit /b 0

:build_failed
echo [错误] 编译失败! 可能原因:
echo 1. 未安装Go环境
echo 2. 依赖未安装(运行 go mod tidy)
echo 3. 代码存在语法错误
pause
exit /b 1