# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 提供代码库操作指南。

## 项目概述

基于 Go 的管理后台系统，提供 REST API。使用 Gin 框架、SQLite + GORM，实现用户管理和权限控制。

**核心功能**：
- 用户管理：登录认证、用户 CRUD、角色权限
- 优惠券管理：优惠券类型、优惠券发放
- 我的优惠券：用户优惠券领取与查询

**技术栈**：
- 语言：Go 1.24
- Web 框架：Gin
- ORM：GORM
- 数据库：SQLite
- 日志：Zap

## 常用命令

```bash
# 运行应用
go run main.go
go run main.go -p 3000        # 自定义端口
go run main.go -d             # 守护进程模式
go run main.go -fl            # 启用文件日志

# 构建
go build -o pionex-administrative-sys

# 测试
go test ./...
go test -run TestName ./path/to/package
go test -v ./db/...           # 详细测试某个包

# 代码质量
go fmt ./...
go vet ./...
go mod tidy
```

## 架构

### 目录结构

```
main.go                    入口，优雅关闭处理
├── server/
│   ├── server.go          Gin 引擎配置，HTTP 服务生命周期
│   ├── handler/           路由注册，按业务域划分
│   │   ├── handler.go     路由总入口，注册各模块
│   │   ├── user/          用户认证、用户管理
│   │   ├── coupon/        优惠券管理
│   │   └── my_coupon/     用户优惠券
│   └── middleware/        中间件
│       ├── auth.go        JWT 认证
│       └── logger.go      请求日志、Recovery
├── db/                    数据层
│   ├── db.go              SQLite 连接，自动迁移
│   ├── user.go            用户模型，位掩码权限
│   ├── role.go            角色模型
│   ├── coupon.go          优惠券模型
│   ├── coupon_type.go     优惠券类型模型
│   └── errors.go          业务错误定义
├── utils/
│   ├── app/               应用启动
│   │   ├── args.go        命令行参数解析
│   │   └── daemon/        守护进程模式
│   ├── logger/            Zap 日志封装
│   ├── resp.go            统一响应封装
│   ├── jwt.go             JWT 工具
│   ├── crypto.go          加密工具
│   └── env.go             环境变量读取
└── static/                静态资源
```

### 关键模式

**API 路由**：所有 API 路由在 `/api/v1` 下，按业务域组织处理器。健康检查接口：`/health`。

**用户角色**：位掩码权限定义在 `db/user.go`：
- `RoleAdmin=1`（管理员）
- `RoleLogin=2`（登录）

**响应格式**：使用 `utils.Resp(code, msg, data).Success(c)` 保持统一：

```json
{
  "code": 200,
  "msg": "OK",
  "data": {}
}
```

**错误处理**：
- 业务错误定义在 `db/errors.go`，使用 `errors.New()` 定义
- 响应失败使用 `utils.Resp(code, msg, data).Fail(c)`，返回 HTTP 401

**数据库约定**：
- SQLite 文件存储在数据目录
- 启动时自动迁移（AutoMigrate）
- 模型通过 `TableName()` 方法定义表名

## 环境与配置

### 命令行参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-p` | 服务端口 | `8080` |
| `-d` | 守护进程模式 | `false` |
| `-fl` | 启用文件日志 | `false` |

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `PAS_HOME` | 应用数据根目录 | `~/.pas/` |

### 目录结构

| 路径 | 说明 |
|------|------|
| `~/.pas/data/` | SQLite 数据库文件 |
| `~/.pas/logs/` | 日志文件（启用 `-fl` 时） |

### 默认账户

首次启动会创建默认管理员账户：
- 账号：`admin`
- 密码：`123456`

> **安全提示**：生产环境部署后请立即修改默认密码。

## 开发指南

### 添加新接口

1. **创建 Handler**：在 `server/handler/` 下新建目录，创建 `xxx_handler.go`
2. **注册路由**：在 handler 文件中实现 `Register(r gin.IRouter)` 函数
3. **挂载模块**：在 `server/handler/handler.go` 的 `Register` 函数中导入并注册

示例：
```go
// server/handler/example/example_handler.go
package example

import "github.com/gin-gonic/gin"

func Register(r gin.IRouter) {
    g := r.Group("/example")
    g.GET("/list", listHandler)
}

func listHandler(c *gin.Context) {
    // 实现逻辑
}
```

### 添加新数据模型

1. 在 `db/` 下创建模型文件，定义结构体
2. 实现 `TableName()` 方法指定表名
3. 在 `db/db.go` 的 `Init()` 中添加 AutoMigrate

### 测试规范

- 测试文件命名：`*_test.go`
- 运行单个测试：`go test -run TestFuncName ./path/to/package`
- 查看覆盖率：`go test -cover ./...`

## 部署

### 直接运行

```bash
# 构建
go build -o pionex-administrative-sys

# 后台运行
./pionex-administrative-sys -d -p 8080

# 或使用 nohup
nohup ./pionex-administrative-sys -p 8080 -fl > /dev/null 2>&1 &
```

### Systemd 服务（Linux）

```ini
# /etc/systemd/system/pas.service
[Unit]
Description=Pionex Administrative System
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/pas
ExecStart=/opt/pas/pionex-administrative-sys -p 8080 -fl
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable pas
sudo systemctl start pas
```

## 注意事项

- **JWT 密钥**：生产环境请确保 JWT 签名密钥的安全性
- **数据备份**：定期备份 `~/.pas/data/` 目录下的 SQLite 文件
- **日志清理**：启用文件日志时，注意定期清理 `~/.pas/logs/` 目录
- **权限控制**：新增接口时注意添加 `middleware.Auth()` 中间件进行认证