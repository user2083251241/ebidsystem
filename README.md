# 在线证券交易系统

## 技术选型

### 后端

`Golang` + `Fiber` + `GORM` + `GoMage`

### 前端

`JavaScript` + `Vue.js` + `Axios`

### 测试

`Postman` + `anaconda`

## 分支管理

### `main`

*不要在`main`分支中存放开发中的项目文件夹！*

### `backend_jiongs`

由 @champNoob(jiongs) 主导的后端开发版本

### `backend_stable`

由 @champNoob(jiongs) 主导的后端稳定版本

### `frontend`

由 @user2083251241 主导的前端开发版本

## 项目结构

### 后端（`./backendDev`）

backend/
├── bin/                   # 编译输出目录
│   ├── ebidsystem.exe        # 可执行文件
│   └── logs                  # 日志目录
│       ├── service.log    # 存放 HTTP 请求、错误日志等通用日志
│       ├── error.log      # 单独记录错误级别的日志（可通过日志库的 Level 过滤）
│       └── match.log      # 存放撮合引擎业务日志
├── config/                # 配置管理
│   └── config.go             # 读取环境变量
├── controllers/           # 控制器（处理 HTTP 请求）
│   ├── auth.go               # 注册/登录/注销
│   ├── client.go             # 客户相关功能
│   ├── common.go             # 公共控制器（如对数据库/JWT/错误的处理）
│   ├── order.go              # 订单创建、查询、取消
│   ├── sales.go              # 销售相关功能（草稿、提交审批）
│   ├── seller.go             # 卖家授权管理
│   └── trader.go             # 交易员授权管理
├── middleware/            # 中间件定义
│   ├── auth.go               # 角色权限校验中间件（如 SellerOnly, SalesOnly）
│   ├── jwt.go                # JWT 认证中间件
│   └── logging.go            # 请求日志中间件
├── models/                # 数据模型定义
│   ├── user.go               # 用户模型
│   ├── order.go              # 订单模型
│   ├── stock.go              # 股票模型（暂不实现）
│   ├── trades.go             # 成交信息（撮合成功后）
│   └── authorization.go      # 卖家-销售授权模型（已定义在 ./user.go 中，暂不独立出来）
├── routes/                # 路由定义
│   └── api.go                # API 路由注册
├── services/              # 核心业务逻辑
│   ├── matching.go           # 订单撮合引擎
│   └── order.go              # 订单状态管理
├── .env                   # 环境变量（开发环境配置）
├── go.mod                 # Go 模块依赖
├── go.sum                 # 依赖校验
├── main.go                # 应用入口（初始化、启动服务）
└── magefile.go            # 自动化构建文件（取得管理员权限+授权通过防火墙+编译+运行+输出日志）
