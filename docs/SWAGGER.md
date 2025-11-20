# Swagger API 文档使用指南

## 📚 概述

项目已集成 Swagger API 文档，可以通过 Web 界面查看和测试所有 API 接口。

## 🚀 快速开始

### 1. 生成 Swagger 文档

安装 swag CLI：
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

生成文档：
```bash
# 在项目根目录执行
swag init

# 或使用 Make 命令
make swagger
```

这会在 `docs/` 目录下生成 `docs.go`、`swagger.json`、`swagger.yaml` 文件。

### 2. 启动项目

```bash
go run main.go
```

### 3. 访问 Swagger UI

打开浏览器访问：
```
http://localhost:8080/swagger/index.html
```

你将看到漂亮的 API 文档界面！

## 📝 编写 Swagger 注解

### 基本格式

在控制器方法上添加注释：

```go
// @Summary      接口摘要
// @Description  接口详细描述
// @Tags         标签（分组）
// @Accept       json
// @Produce      json
// @Param        参数名  位置  类型  必填  "说明"
// @Success      200  {object}  ResponseType  "成功描述"
// @Failure      400  {object}  ResponseType  "失败描述"
// @Security     ApiKeyAuth
// @Router       /path [method]
func (ctrl *Controller) Method(ctx *gin.Context) {
    // 实现
}
```

### 示例：用户注册接口

```go
// @Summary      用户注册
// @Description  注册新用户账号
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest  true  "注册参数"
// @Success      200      {object}  utils.Response{data=object{user_id=int,username=string}}
// @Failure      400      {object}  utils.Response
// @Router       /auth/register [post]
func (ctrl *AuthController) Register(ctx *gin.Context) {
    // ...
}
```

### 参数位置

| 位置 | 说明 | 示例 |
|------|------|------|
| `path` | 路径参数 | `/users/{id}` |
| `query` | 查询参数 | `?page=1&size=10` |
| `body` | 请求体 | JSON Body |
| `header` | 请求头 | `X-Token` |
| `formData` | 表单数据 | `multipart/form-data` |

### 数据类型

| 类型 | 说明 |
|------|------|
| `string` | 字符串 |
| `integer` / `int` | 整数 |
| `number` | 数字 |
| `boolean` | 布尔值 |
| `array` | 数组 |
| `object` | 对象 |
| `file` | 文件 |

### 认证配置

需要认证的接口添加：
```go
// @Security     ApiKeyAuth
```

或使用 Bearer 格式：
```go
// @Security     BearerAuth
```

## 📖 完整示例

### 1. 公开接口（无需认证）

```go
// @Summary      用户登录
// @Description  使用用户名和密码登录，返回 JWT Token
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "登录参数"
// @Success      200      {object}  utils.Response{data=LoginResponse}  "登录成功"
// @Failure      400      {object}  utils.Response  "参数错误"
// @Failure      11006    {object}  utils.Response  "用户名或密码错误"
// @Failure      11008    {object}  utils.Response  "账号已被禁用"
// @Router       /auth/login [post]
func (ctrl *AuthController) Login(ctx *gin.Context) {
    // ...
}
```

### 2. 需要认证的接口

```go
// @Summary      获取用户列表
// @Description  分页获取所有用户列表（需要登录）
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        page       query     int  false  "页码"  default(1)  minimum(1)
// @Param        page_size  query     int  false  "每页数量"  default(10)  minimum(1)  maximum(100)
// @Success      200        {object}  utils.Response{data=models.PageResponse}  "获取成功"
// @Failure      10005      {object}  utils.Response  "需要登录"
// @Failure      13003      {object}  utils.Response  "查询失败"
// @Router       /users [get]
func (ctrl *UserController) GetUserList(ctx *gin.Context) {
    // ...
}
```

### 3. 路径参数

```go
// @Summary      获取指定用户
// @Description  根据用户 ID 获取用户详细信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "用户ID"  minimum(1)
// @Success      200  {object}  utils.Response{data=models.User}  "获取成功"
// @Failure      14001  {object}  utils.Response  "无效的用户 ID"
// @Failure      11002  {object}  utils.Response  "用户不存在"
// @Router       /users/{id} [get]
func (ctrl *UserController) GetUserByID(ctx *gin.Context) {
    // ...
}
```

## 🔧 Makefile 集成

在 `Makefile` 中添加：

```makefile
# 生成 Swagger 文档
swagger:
	@echo "📝 生成 Swagger 文档..."
	swag init
	@echo "✅ Swagger 文档生成完成"
	@echo "访问: http://localhost:8080/swagger/index.html"
```

使用：
```bash
make swagger
```

## 📊 Swagger 配置

主配置在 `main.go` 文件顶部：

```go
// @title           Go-Demo API
// @version         2.0
// @description     API 描述
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        X-Token
```

## 💡 高级用法

### 1. 自定义响应模型

```go
// @Success 200 {object} utils.Response{data=LoginResponse}
```

### 2. 多个响应状态

```go
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
```

### 3. 数组响应

```go
// @Success 200 {array} models.User
```

### 4. 枚举值

```go
// @Param status query string false "状态" Enums(active, inactive, deleted)
```

### 5. 示例值

```go
// @Param page query int false "页码" default(1) example(1)
```

## 🎯 最佳实践

1. **每个接口都添加注解**
2. **详细的描述和示例**
3. **正确的错误码**
4. **使用标签分组**
5. **添加认证信息**

## 🌐 访问 Swagger UI

启动项目后访问：
- Swagger UI: http://localhost:8080/swagger/index.html
- Swagger JSON: http://localhost:8080/swagger/doc.json
- Swagger YAML: http://localhost:8080/swagger/doc.yaml

## 📚 参考资料

- [Swag GitHub](https://github.com/swaggo/swag)
- [Swag 注解规范](https://github.com/swaggo/swag#declarative-comments-format)
- [Swagger Specification](https://swagger.io/specification/)

## ✅ 验证

生成文档后，检查：
- [ ] `docs/docs.go` 文件已生成
- [ ] `docs/swagger.json` 文件已生成
- [ ] 可以访问 Swagger UI
- [ ] 所有接口都显示正确
- [ ] 可以在 UI 中测试接口

---

**提示**: 每次修改注解后，需要重新运行 `swag init` 生成文档！

