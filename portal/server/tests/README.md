# 明河门户系统后端 - 单元测试

本目录包含项目的单元测试代码。

## 测试目录结构

```
tests/
├── stringx/          # 字符串工具测试
│   ├── rand_test.go
│   └── strings_test.go
├── timex/            # 时间工具测试
│   └── timex_test.go
├── encrypt/           # 加密工具测试
│   ├── aes_test.go
│   └── triple_des_test.go
├── context/           # 上下文工具测试
│   └── context_test.go
├── config/           # 配置加载测试
│   └── config_test.go
└── md5_salt_test.go  # MD5加盐哈希测试
```

## 运行测试

### 运行所有测试
```bash
go test ./tests/... -v
```

### 运行特定包的测试
```bash
# 字符串工具测试
go test ./tests/stringx/... -v

# 时间工具测试
go test ./tests/timex/... -v

# 加密工具测试
go test ./tests/encrypt/... -v

# 配置测试
go test ./tests/config/... -v
```

### 运行单个测试
```bash
go test ./tests/stringx/... -v -run TestSplit
```

### 运行测试并查看覆盖率
```bash
go test ./tests/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 测试覆盖率目标

| 模块 | 目标覆盖率 | 说明 |
|------|----------|------|
| stringx | >80% | 字符串工具函数 |
| timex | >80% | 时间工具函数 |
| encrypt | >90% | 加密解密函数（安全关键）|
| context | >85% | 上下文管理函数 |
| config | >75% | 配置加载函数 |

## 测试说明

### stringx 测试
- **Split**: 字符串分割功能
- **StrLike**: SQL LIKE 模式生成
- **GenerateRand**: 随机字符串生成（验证码功能）
- **GetTradeNo**: 交易号生成

### timex 测试
- **ParseChineseDateTime**: 中文日期时间解析
- **TimestampToTime**: 秒级时间戳转换
- **TimestampMsToTime**: 毫秒级时间戳转换

### encrypt 测试
- **AES**: AES-256-GCM 加密解密
- **3DES**: 3DES-CBC 加密解密
- **PKCS7Padding**: PKCS7 填充和去填充

### context 测试
- **SetUserToCtx**: 设置用户信息到上下文
- **DetailUserFromCtx**: 从上下文获取用户信息

### config 测试
- **LoadConfig**: 配置文件加载
- **环境变量扩展**: ${VAR} 占位符替换
- **配置提取**: LoadMysqlConfig、LoadRedisConfig

## 添加新测试

### 测试文件命名规范
测试文件应放在对应的 `tests/` 子目录下，文件名格式为 `<filename>_test.go`

### 测试函数命名规范
```go
func Test<FunctionName>(t *testing.T) {
    // 测试单个功能点
}

func Test<FunctionName><Scenario>(t *testing.T) {
    // 测试特定场景
}
```

### 测试最佳实践

1. **表驱动测试**: 使用结构体切片定义多个测试用例
2. **子测试**: 使用 `t.Run()` 创建子测试，便于定位问题
3. **清晰的错误信息**: 错误消息应包含输入和期望值
4. **边界条件**: 测试空值、最大值、最小值等边界情况
5. **并发安全**: 对于并发访问的函数，使用 `t.Parallel()`

### 测试示例

```go
func TestExample(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "normal case",
            input:    "test",
            expected: "TEST",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ToUpperCase(tt.input)
            if result != tt.expected {
                t.Errorf("ToUpperCase(%q) = %q, want %q",
                    tt.input, result, tt.expected)
            }
        })
    }
}
```

## CI/CD 集成

测试应在 CI/CD 流程中自动运行：

```yaml
- name: Run tests
  run: go test ./tests/... -v -cover

- name: Upload coverage
  run: go tool cover -html=coverage.out -o coverage.html
```

## 注意事项

1. **测试隔离**: 每个测试应该独立，不依赖其他测试的状态
2. **资源清理**: 使用 `t.Cleanup()` 或 `defer` 清理测试资源
3. **临时文件**: 使用 `t.TempDir()` 创建临时目录和文件
4. **Mock 外部依赖**: 对于网络、数据库等外部依赖，使用 mock
5. **超时控制**: 使用 `go test -timeout 30s` 限制测试运行时间
