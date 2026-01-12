# Pionex Administrative System

基于 Go 的轻量级管理后台系统，提供 REST API 服务。

## 功能特性

- **用户管理**：用户注册、登录认证、角色权限控制
- **优惠券系统**：优惠券类型管理、优惠券发放与领取
- **轻量部署**：SQLite 数据库，单文件部署，无外部依赖

## 技术栈

| 组件 | 技术 |
|------|------|
| 语言 | Go 1.24 |
| Web 框架 | Gin |
| ORM | GORM |
| 数据库 | SQLite |
| 日志 | Zap |
| 认证 | JWT |

## 快速开始

### 环境要求

- Go 1.24+

### 运行

```bash
# 克隆项目
git clone <repository-url>
cd pionex-administrative-sys

# 安装依赖
go mod tidy

# 启动服务（默认端口 8080）
go run main.go

# 自定义端口
go run main.go -p 3000

# 守护进程模式
go run main.go -d

# 启用文件日志
go run main.go -fl
```

### 构建

```bash
go build -o pionex-administrative-sys
```

### 测试

```bash
go test ./...
go test -cover ./...    # 查看覆盖率
```

## 项目结构

```
├── main.go                 # 入口，优雅关闭处理
├── server/
│   ├── server.go           # Gin 引擎配置，HTTP 服务生命周期
│   ├── handler/            # 路由处理器
│   │   ├── handler.go      # 路由总入口
│   │   ├── user/           # 用户认证、管理
│   │   ├── coupon/         # 优惠券管理
│   │   └── my_coupon/      # 我的优惠券
│   └── middleware/         # 中间件（日志、恢复、JWT 认证）
├── db/                     # 数据模型和数据访问
│   ├── db.go               # 数据库连接，自动迁移
│   ├── user.go             # 用户模型
│   ├── role.go             # 角色定义
│   ├── coupon.go           # 优惠券模型
│   ├── coupon_type.go      # 优惠券类型
│   └── errors.go           # 业务错误定义
├── utils/                  # 工具函数
│   ├── app/                # 命令行解析、守护进程
│   ├── logger/             # 日志配置
│   ├── resp.go             # 统一响应封装
│   ├── jwt.go              # JWT 工具
│   └── crypto.go           # 加密工具
└── static/                 # 静态资源
```

## API 接口

基础路径：`/api/v1`

| 模块 | 路径 | 说明 |
|------|------|------|
| 健康检查 | `GET /health` | 服务健康状态 |
| 用户 | `/api/v1/user/*` | 用户认证、管理 |
| 优惠券 | `/api/v1/coupon/*` | 优惠券管理 |
| 我的优惠券 | `/api/v1/my_coupon/*` | 用户优惠券 |

### 响应格式

```json
{
  "code": 200,
  "msg": "OK",
  "data": {}
}
```

## 配置

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

### 数据目录

| 路径 | 说明 |
|------|------|
| `~/.pas/data/` | SQLite 数据库文件 |
| `~/.pas/logs/` | 日志文件（启用 `-fl` 时） |

## 部署

### 直接运行

```bash
# 构建
go build -o pionex-administrative-sys

# 后台运行
./pionex-administrative-sys -d -p 8080
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

## 默认账户

首次启动会创建默认管理员账户：

- 账号：`admin`
- 密码：`123456`

> **安全提示**：生产环境部署后请立即修改默认密码！

## 开发

### 代码质量

```bash
go fmt ./...
go vet ./...
go mod tidy
```

### 添加新模块

1. 在 `server/handler/` 下创建新目录
2. 实现 `Register(r gin.IRouter)` 函数
3. 在 `server/handler/handler.go` 中注册

详细开发指南请参考 [CLAUDE.md](./CLAUDE.md)

## License

MIT