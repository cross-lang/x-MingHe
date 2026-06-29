# 调试模式使用指南

## 概述

MingHe 门户系统后端提供了调试模式（Debug Mode），用于开发和测试环境，提供更详细的日志输出、禁用限流、美化 JSON 响应等功能。

## 启用调试模式

### 方法 1：修改配置文件

在 `config.yaml` 或 `config_local.yaml` 中设置 `Debug` 字段为 `true`：

```yaml
ServerPort: 8088
Debug: true  # 启用调试模式
Mysql:
  Dsn: root:123456@tcp(127.0.0.1:3306)/minghe?charset=utf8mb4&parseTime=True&loc=Local
  # ... 其他配置
```

### 方法 2：通过环境变量

使用环境变量覆盖配置：

```bash
export DEBUG=true
./portal -c config.yaml
```

或者在配置文件中使用环境变量：

```yaml
Debug: ${DEBUG:false}
```

## 调试模式功能

### 1. 详细日志输出

- **SQL 查询日志**：所有数据库查询都会被记录，包括查询语句、执行时间、影响行数
- **请求/响应 Body**：记录所有请求和响应的完整内容
- **模块化日志**：日志包含模块信息（api、service、repository、database 等）
- **敏感数据脱敏**：自动脱敏手机号、邮箱、身份证、密码等敏感信息

### 2. 禁用限流

- 调试模式下，全局限流和端点限流都会被禁用
- 方便开发和测试，避免被限流器阻塞

### 3. 美化 JSON 响应

- API 响应的 JSON 输出会被格式化，便于阅读
- 自动缩进和换行，提高可读性

### 4. Gin Debug 模式

- 使用 Gin 框架的 Debug 模式
- 更详细的错误信息和堆栈跟踪

## 日志示例

### 普通模式

```json
{
  "level": "info",
  "time": "2026-04-23 17:15:30.123",
  "msg": "response",
  "method": "POST",
  "path": "/v1/users/login",
  "status": 200,
  "duration": "150ms"
}
```

### 调试模式

```json
{
  "level": "debug",
  "time": "2026-04-23 17:15:30.123",
  "msg": "SQL执行成功",
  "sql": "SELECT * FROM x_users WHERE phone_number = ? LIMIT 1",
  "rows": 1,
  "duration": "45ms",
  "module": "database",
  "trace_id": "abc-123-def-456"
},
{
  "level": "debug",
  "time": "2026-04-23 17:15:30.234",
  "msg": "request",
  "method": "POST",
  "path": "/v1/users/login",
  "query": "",
  "request_body": "{\"phone_number\":\"138****5678\",\"code\":\"123456\"}",
  "module": "api"
},
{
  "level": "debug",
  "time": "2026-04-23 17:15:30.345",
  "msg": "response",
  "status": 200,
  "duration": "150ms",
  "body": "{\"code\":0,\"msg\":\"success\",\"data\":{\"token\":\"xxx\",\"expires_at\":1234567890}}",
  "module": "api"
}
```

## JSON 响应对比

### 普通模式

```json
{"code":0,"msg":"success","data":{"token":"xxx","expires_at":1234567890}}
```

### 调试模式

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "token": "xxx",
    "expires_at": 1234567890
  }
}
```

## 健康检查端点

调试模式下，可以使用健康检查端点来监控系统状态：

### 基本健康检查

```bash
curl http://localhost:8088/health
```

响应：

```json
{
  "code": 0,
  "msg": "healthy",
  "data": {
    "status": "healthy",
    "timestamp": 1713861330
  }
}
```

### 详细健康检查

```bash
curl http://localhost:8088/health?detailed=true
```

响应：

```json
{
  "code": 0,
  "msg": "healthy",
  "data": {
    "status": "healthy",
    "timestamp": 1713861330,
    "components": {
      "mysql": {
        "name": "mysql",
        "status": "healthy",
        "message": "连接池: 空闲=8/最大=10, 使用中=2, 总数=10",
        "latency_ms": 5
      },
      "redis": {
        "name": "redis",
        "status": "healthy",
        "message": "连接池: 空闲=5, 总数=10",
        "latency_ms": 2
      }
    }
  }
}
```

## 开发工具集成

### Air 热重载

使用 Air 进行热重载开发时，配置文件中启用调试模式：

```bash
# 确保 config_local.yaml 中 Debug: true
make dev
```

### VS Code 调试

创建 `.vscode/launch.json`：

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug MingHe Portal",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/server",
      "env": {
        "DEBUG": "true"
      },
      "args": ["-c", "config_local.yaml"],
      "showLog": true
    }
  ]
}
```

## 注意事项

### ⚠️ 安全警告

**不要在生产环境启用调试模式！**

调试模式会：
- 记录所有请求和响应内容（可能包含敏感数据）
- 禁用限流（可能导致资源耗尽）
- 暴露详细的错误信息（可能泄露系统内部信息）
- 降低性能（详细的日志记录会增加开销）

### 生产环境配置

确保生产环境配置文件中 `Debug` 设置为 `false`：

```yaml
# config_prod.yaml
ServerPort: 8088
Debug: false  # 生产环境必须为 false
Logger:
  Level: info  # 生产环境使用 info 级别
  EnableRemote: true  # 建议启用远程日志
```

### 性能影响

调试模式会对性能产生以下影响：

- **日志量增加**：详细日志会显著增加 I/O 开销
- **JSON 格式化**：美化 JSON 会增加 CPU 使用
- **禁用限流**：可能导致系统资源耗尽
- **内存占用增加**：更多的日志缓冲和对象创建

建议仅在开发和测试环境使用调试模式。

## 调试技巧

### 1. 查看特定模块的日志

日志包含 `module` 字段，可以过滤查看特定模块的日志：

```bash
# 只查看数据库日志
grep '"module":"database"' logs/MingHe_portal_*.log

# 只查看 API 日志
grep '"module":"api"' logs/MingHe_portal_*.log
```

### 2. 根据 Trace ID 追踪请求

使用 X-Trace-ID 头来追踪完整的请求链路：

```bash
# 发送请求时指定 Trace ID
curl -H "X-Trace-ID: my-custom-trace-id" http://localhost:8088/v1/users/login

# 在日志中查找该 Trace ID 的所有日志
grep '"trace_id":"my-custom-trace-id"' logs/MingHe_portal_*.log
```

### 3. 慢查询分析

调试模式下，超过 200ms 的 SQL 查询会被标记为慢查询并记录警告：

```json
{
  "level": "warn",
  "msg": "慢SQL查询",
  "sql": "SELECT * FROM x_users WHERE ...",
  "rows": 100,
  "time": 250,
  "module": "database"
}
```

### 4. 使用健康检查监控

在调试过程中，定期检查系统健康状态：

```bash
# 简单健康检查
watch -n 5 'curl -s http://localhost:8088/health | jq .'

# 详细健康检查
watch -n 5 'curl -s http://localhost:8088/health?detailed=true | jq .'
```

## 常见问题

### Q: 为什么启用调试模式后性能明显下降？

A: 调试模式会增加大量日志记录和 JSON 格式化，这些都会影响性能。建议在生产环境使用普通模式。

### Q: 调试模式下如何查看完整日志？

A: 检查 `logs/` 目录下的日志文件，或使用 `tail -f` 实时查看：

```bash
tail -f logs/MingHe_portal_*.log
```

### Q: 如何在调试模式下启用限流？

A: 修改 `internal/app/app.go` 中的限流中间件注册逻辑，移除 `!conf.Debug` 的判断。

### Q: 调试模式的日志级别是什么？

A: 调试模式不改变日志级别设置，但会增加更多的日志输出。可以通过 `Logger.Level` 配置调整日志级别。

## 相关文档

- [配置说明](./CONFIG.md)
- [部署指南](./DEPLOYMENT.md)
- [API 文档](./API.md)
- [开发指南](./DEVELOPMENT.md)
