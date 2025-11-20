# API 路由使用指南 📖

## 🎯 路由设计原则

本项目**只使用 GET 和 POST 方法**，通过路径区分操作类型。

---

## 📋 完整 API 路由

### 健康检查（无需认证）

```
GET  /ping                    简单健康检查
GET  /health                  完整健康检查（检查 MySQL/Redis/MongoDB）
GET  /ready                   就绪检查（K8s Readiness Probe）
GET  /live                    存活检查（K8s Liveness Probe）

# 也可以通过 /api/v1 访问
GET  /api/v1/ping
GET  /api/v1/health
GET  /api/v1/ready
GET  /api/v1/live
```

### 用户认证（公开）

```
POST /api/v1/auth/register    用户注册
POST /api/v1/auth/login       用户登录
```

### 用户认证（需要 Token）

```
POST /api/v1/auth/logout      用户登出
GET  /api/v1/auth/user-info   获取当前用户信息
```

### 用户管理（需要 Token）

```
GET  /api/v1/users                获取用户列表（分页）
GET  /api/v1/users/:id            获取指定用户详情
POST /api/v1/users/update         更新当前用户信息
POST /api/v1/users/:id/delete     删除指定用户
```

### 文章管理

```
# 公开接口
GET  /api/v1/articles             获取文章列表（分页、搜索）
GET  /api/v1/articles/:id         获取文章详情

# 需要认证
POST /api/v1/articles             创建文章
POST /api/v1/articles/:id/update  更新文章
POST /api/v1/articles/:id/delete  删除文章
POST /api/v1/articles/:id/like    点赞文章
```

---

## 🎯 HTTP 方法说明

### GET 方法（查询数据）

**特点**:
- 只读操作
- 不修改服务器状态
- 可以缓存
- 幂等性（多次请求结果相同）

**使用场景**:
- 获取列表
- 获取详情
- 健康检查

### POST 方法（修改数据）

**特点**:
- 修改服务器状态
- 不缓存
- 通过路径区分操作类型

**使用场景**:
- 创建数据（如创建文章）
- 更新数据（如 /update）
- 删除数据（如 /delete）
- 其他操作（如 /like）

---

## 📝 使用示例

### 1. 用户注册和登录

```http
### 注册
POST http://localhost:8080/api/v1/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "password": "123456",
  "email": "test@example.com"
}

### 登录
POST http://localhost:8080/api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "123456"
}
```

### 2. 用户管理

```http
### 获取用户列表
GET http://localhost:8080/api/v1/users?page=1&page_size=10
X-Token: YOUR_TOKEN

### 获取用户详情
GET http://localhost:8080/api/v1/users/1
X-Token: YOUR_TOKEN

### 更新用户信息（POST）
POST http://localhost:8080/api/v1/users/update
Content-Type: application/json
X-Token: YOUR_TOKEN

{
  "email": "new@example.com",
  "phone": "13800138000"
}

### 删除用户（POST）
POST http://localhost:8080/api/v1/users/2/delete
X-Token: YOUR_TOKEN
```

### 3. 文章管理

```http
### 创建文章
POST http://localhost:8080/api/v1/articles
Content-Type: application/json
X-Token: YOUR_TOKEN

{
  "title": "我的文章",
  "content": "文章内容",
  "tags": ["Go", "技术"]
}

### 更新文章（POST）
POST http://localhost:8080/api/v1/articles/xxx/update
Content-Type: application/json
X-Token: YOUR_TOKEN

{
  "title": "更新后的标题",
  "content": "更新后的内容"
}

### 删除文章（POST）
POST http://localhost:8080/api/v1/articles/xxx/delete
X-Token: YOUR_TOKEN

### 点赞文章（POST）
POST http://localhost:8080/api/v1/articles/xxx/like
X-Token: YOUR_TOKEN
```

---

## 🔄 路径设计规则

### 查询操作 → GET

```
GET  /resource              获取列表
GET  /resource/:id          获取详情
```

### 修改操作 → POST

```
POST /resource              创建
POST /resource/update       更新（当前用户）
POST /resource/:id/update   更新（指定资源）
POST /resource/:id/delete   删除
POST /resource/:id/action   其他操作（如 /like）
```

---

## 📊 路由对比

### RESTful 风格
```
GET    /users       获取列表
GET    /users/:id   获取详情
POST   /users       创建
PUT    /users/:id   更新
DELETE /users/:id   删除
```

### 本项目风格（只用 GET/POST）
```
GET  /users              获取列表
GET  /users/:id          获取详情
POST /users              创建（如果有）
POST /users/update       更新
POST /users/:id/delete   删除
```

**优势**:
- ✅ 更简单（只有两种方法）
- ✅ 更兼容（老旧客户端、防火墙）
- ✅ 更明确（路径即操作）

---

## 🚀 快速测试

### 使用 curl

```bash
# 健康检查
curl http://localhost:8080/health

# 更新用户（注意是 POST）
curl -X POST http://localhost:8080/api/v1/users/update \
  -H "X-Token: YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"new@example.com"}'

# 删除用户（注意是 POST）
curl -X POST http://localhost:8080/api/v1/users/2/delete \
  -H "X-Token: YOUR_TOKEN"
```

### 使用 Swagger UI

访问: http://localhost:8080/swagger/index.html

现在所有接口都是 GET 或 POST 方法，更容易理解和使用。

---

## ✅ 验证清单

重启项目后，检查：

- [ ] Swagger UI 中没有 PUT 和 DELETE 方法
- [ ] 所有接口都是 GET 或 POST
- [ ] `/users/update` 可以正常更新
- [ ] `/users/:id/delete` 可以正常删除
- [ ] `/articles/:id/update` 可以正常更新
- [ ] `/articles/:id/delete` 可以正常删除

---

<div align="center">

## 🎉 路由更新完成！

**简单 · 清晰 · 实用**

查看完整列表: [API_ROUTES_UPDATED.md](API_ROUTES_UPDATED.md)

</div>

