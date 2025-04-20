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

```txt
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
