# 在线证券交易系统

## 一、技术选型

### 后端

`Golang` + `Fiber` + `GORM` + `GoMage`

### 前端

`JavaScript` + `Vue.js` + `Axios`

### 测试

`Postman` + `anaconda`

## 二、分支管理

### `main`

*不要在`main`分支中存放开发中的项目文件夹！*

### `backend_jiongs`

由 @champNoob(jiongs) 主导的后端开发版本

### `backend_stable`

由 @champNoob(jiongs) 主导的后端稳定版本

### `frontend`

由 @user2083251241 主导的前端开发版本

## 三、项目结构

### 后端（`./backendDev`）

```txt
backend/
├── bin/                   # 编译输出目录
│   ├── ebidsystem.exe        # 可执行文件
│   └── logs                  # 日志目录
│       ├── service.log          # 存放 HTTP 请求、错误日志等通用日志
│       ├── error.log            # 单独记录错误级别的日志（可通过日志库的 Level 过滤）
│       └── match.log            # 存放撮合引擎业务日志
├── config/                # 配置管理
│   ├── config.go             # 读取环境变量
│   └── database.go           # 构建全局数据库实例并初始化数据库连接
├── controllers/           # 控制器（处理 HTTP 请求）
│   ├── auth.go               # 专注认证与授权逻辑（注册/登录/注销）
│   ├── base_controller.go    # 基础控制器（复用请求解析和校验）
│   ├── client.go             # 客户相关功能
│   ├── common.go             # 公共控制器（与业务无关的通用工具，如对数据库/JWT/错误的处理）
│   ├── sales.go              # 销售相关功能（草稿、提交审批）
│   ├── seller.go             # 卖家授权管理
│   └── trader.go             # 交易员授权管理
├── middleware/            # 中间件定义
│   ├── auth.go               # 角色权限校验中间件（如 SellerOnly, SalesOnly）
│   ├── jwt.go                # JWT 认证中间件
│   └── logging.go            # 请求日志中间件
├── models/                # 数据模型定义
│   ├── audit_log.go          # 审计日志结构体
│   ├── authorization.go      # 卖家-销售授权模型
│   ├── order.go              # 订单模型
│   ├── stock.go              # 股票模型（暂不实现）
│   ├── trades.go             # 成交信息（撮合成功后）
│   └── user.go               # 用户模型
├── routes/                # 路由定义
│   └── api.go                # API 路由注册
├── services/              # 核心业务逻辑
│   ├── matching.go           # 订单撮合引擎
│   ├── matching_test.go      # 订单撮合引擎测试（单元测试）
│   ├── order_query.go        # 订单查询管理
│   ├── order_services.go     # 订单业务管理
│   ├── requests.go           # 订单请求管理
│   └── user_services.go      # 用户管理（注册、登录、注销等）
├── static/                # 静态资源
|   └── assets/                # 前端资源
├── utils/                 # 工具函数
|   ├── ptr.go                 # 指针工具
|   ├── redis.go               # Redis 工具
|   └── validation.go          # 校验工具
├── .env                   # 环境变量（开发环境配置）
├── go.mod                 # Go 模块依赖
├── go.sum                 # 依赖校验
├── magefile.go            # 自动化构建文件（取得管理员权限+授权通过防火墙+编译+运行+输出日志）
└── main.go                # 应用入口（初始化、启动服务）
```

### 前端(`./frontendDev`)

```txt
mytext/          # 项目根目录
├── .vscode/     # VSCode 编辑器的配置文件夹
├── node_modules/ # 项目依赖的第三方库，由 npm 或 yarn 安装
├── public/       # 静态资源文件夹，存放公开访问的资源如图片、图标等
├── src/          # 源代码文件夹，存放项目的主要代码
│   ├── assets/   # 存放静态资源，如图片、样式表等
│   ├── components/ # 存放 Vue 组件
│   │   ├── icons/ # 图标组件文件夹
│   │   ├── operation/ # 操作相关的组件文件夹
│   │   │   ├── Creator.vue  # 创建操作组件
│   │   │   ├── Deletion.vue # 删除操作组件
│   │   │   ├── Modify.vue   # 修改操作组件
│   │   │   ├── Purchase.vue # 购买操作组件
│   │   ├── ClientTableComponent.vue # 客户端表格组件
│   │   ├── ConfirmModal.vue # 确认模态框组件
│   │   ├── LoginPage.vue    # 登录页面组件
│   │   ├── MainPage.vue     # 主页面组件
│   │   ├── ModifyModal.vue  # 修改模态框组件
│   │   ├── Register.vue     # 注册页面组件
│   │   ├── SetupPage.vue    # 设置页面组件
│   │   ├── Sidebar.vue      # 侧边栏组件
│   │   ├── TableComponent.vue # 表格组件
│   │   └── ...              # 其他组件
│   ├── router/     # 路由配置文件夹
│   │   ├── index.js # 路由入口文件
│   ├── utils/      # 工具函数文件夹
│   │   ├── parseToken.js # 解析 JWT token 的工具函数
│   ├── App.vue      # 主应用组件
│   └── main.js      # 入口 JavaScript 文件，初始化 Vue 实例
├── .gitignore     # Git 忽略文件配置
├── index.html     # 项目的 HTML 模板文件
├── jsconfig.json  # JavaScript 项目配置文件
├── package-lock.json # 锁定项目依赖版本的文件
├── package.json   # 项目的 package.json 文件，定义项目的依赖和脚本
└── README.md      # 项目的自述文件，通常包含项目介绍和使用说明 
```

## 四、路由（以后端为准）

```go
// 认证路由组：
authenticated := app.Group("/api", jwtMiddleware)
{
   // 卖家角色路由组：
   seller := authenticated.Group("/seller", middleware.SellerOnly)
   {
      seller.Post("/orders", controllers.CreateSellOrder)         // 创建卖出订单
      seller.Put("/orders/:id", controllers.UpdateOrder)          // 修改订单
      seller.Post("/orders/:id/cancel", controllers.CancelOrder)  // 取消订单
      seller.Get("/orders", controllers.GetSellerOrders)          // 查看卖家订单
      seller.Post("/authorize/sales", controllers.AuthorizeSales) // 授权销售
   }
   // 销售角色路由组：
   sales := authenticated.Group("/sales", middleware.SalesOnly)
   {
      sales.Get("/orders", controllers.GetAuthorizedOrders)     // 查看已授权订单
      sales.Post("/drafts", controllers.CreateDraftOrder)       // 创建订单草稿
      sales.Put("/drafts/:id", controllers.UpdateDraftOrder)    // 修改草稿
      sales.Post("/drafts/:id/submit", controllers.SubmitDraft) // 提交草稿
   }
   // 客户角色路由组：
   client := authenticated.Group("/client")
   {
      client.Post("/orders", controllers.CreateBuyOrder) // 创建买入订单
      client.Get("/orders", controllers.GetClientOrders) // 查看客户订单
   }
   // 交易员角色路由组：
   trader := authenticated.Group("/trader", middleware.TraderOnly)
   {
      trader.Get("/orders", controllers.GetAllOrders) // 查看所有订单
   }
}
```
