# RuleGo Kratos Server

这是使用Kratos框架重写的RuleGo服务器，保持了原有的功能同时提供了更好的架构和可维护性。

## 项目结构

```
kratos-server/
├── api/                    # API定义 (protobuf)
│   └── rulego/v1/         # RuleGo API v1
├── cmd/
│   └── kratos-server/     # 主程序入口
├── configs/               # 配置文件
├── internal/
│   ├── biz/              # 业务逻辑层
│   ├── conf/             # 配置结构
│   ├── data/             # 数据访问层
│   ├── middleware/       # 中间件
│   ├── server/           # 服务器配置
│   └── service/          # 服务层
├── editor/               # 前端编辑器文件
└── third_party/          # 第三方proto文件
```

## 主要功能

1. **GRPC服务**: 提供完整的RuleGo API服务
2. **HTTP服务**: Kratos HTTP服务器 (端口9090)
3. **RuleGo REST服务**: 原生RuleGo REST API (端口9091)
4. **前端界面**: 规则链编辑器 (http://localhost:9091/editor/)
5. **数据库支持**: PostgreSQL数据持久化
6. **认证中间件**: JWT和API Key认证
7. **规则引擎**: 完整的RuleGo规则引擎功能

## 配置说明

配置文件位于 `configs/config.yaml`，主要配置项：

- `server`: Kratos服务器配置
- `data`: 数据库配置
- `rulego`: RuleGo特定配置，包括JWT、用户、MCP等

## 启动步骤

1. 准备PostgreSQL数据库
2. 修改配置文件中的数据库连接信息
3. 构建并启动服务：

```bash
# 生成代码
make api
make config
go generate ./...

# 构建
go build -o ./bin/ ./...

# 启动
./bin/kratos-server -conf ./configs
```

## 访问地址

- Kratos HTTP服务: http://localhost:9090
- RuleGo REST API: http://localhost:9091  
- 前端编辑器: http://localhost:9091/editor/
- GRPC服务: localhost:9000

## API服务

### GRPC服务
- UserService: 用户认证服务
- RuleService: 规则链管理服务  
- ComponentService: 组件管理服务
- LogService: 日志服务
- LocaleService: 本地化服务

### REST API
原有的RuleGo REST API都在端口9091上提供，包括：
- `/api/v1/rules`: 规则链管理
- `/api/v1/components`: 组件管理
- `/api/v1/logs`: 日志管理
- `/editor/`: 前端编辑器

## 开发指南

### 添加新的API
1. 在 `api/rulego/v1/` 中定义protobuf
2. 运行 `make api` 生成代码
3. 在 `internal/service/` 中实现服务逻辑
4. 在 `internal/server/` 中注册服务

### 数据库迁移
数据库表会在启动时自动创建和迁移，模型定义在 `internal/data/model.go`

### 中间件
认证和CORS中间件在 `internal/middleware/` 中定义

## 与原版本的区别

1. **架构改进**: 采用Kratos框架的分层架构
2. **依赖注入**: 使用Wire进行依赖注入
3. **协议支持**: 同时支持GRPC和HTTP
4. **数据持久化**: 使用GORM和PostgreSQL
5. **配置管理**: 统一的配置管理
6. **中间件支持**: 标准化的中间件机制

## 迁移说明

从原版本迁移时需要注意：

1. 数据库从文件存储改为PostgreSQL
2. 配置文件格式有所变化  
3. API保持兼容，客户端无需修改
4. 前端编辑器保持不变

## 故障排除

1. **数据库连接失败**: 检查PostgreSQL服务和配置
2. **端口冲突**: 修改配置文件中的端口设置
3. **前端文件404**: 确保editor目录存在且配置正确

更多信息请参考原项目文档和Kratos官方文档。