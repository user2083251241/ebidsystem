# 编译项目
build:
	go build -o bin/ebidsystem.exe main.go

# 启动服务（开发模式）
run:
	./bin/ebidsystem.exe

# 配置防火墙规则（Windows）
setup-firewall:
	netsh advfirewall firewall delete rule name="Allow ebidsystem Inbound" >nul 2>&1
	netsh advfirewall firewall add rule name="Allow ebidsystem Inbound" dir=in program="$(CURDIR)/bin/ebidsystem.exe" action=allow

# 清理构建产物
clean:
	rm -rf bin/*