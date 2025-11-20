# API 路由更新说明 🔄

## ✅ 路由已更新为只使用 GET 和 POST

---

## 🎯 更新内容

### 变更说明

已将所有 **PUT** 和 **DELETE** 请求改为 **POST** 请求，通过路径区分操作类型。

---

## 📋 更新对比

### 用户管理接口

| 功能 | 旧路由 | 新路由 |
|------|--------|--------|
| 获取列表 | GET `/users` | GET `/users` ✅ 不变 |
| 获取详情 | GET `/users/:id` | GET `/users/:id` ✅ 不变 |
| 更新用户 | PUT `/users` ❌ | POST `/users/update` ✅ |
| 删除用户 | DELETE `/users/:id` ❌ | POST `/users/:id/delete` ✅ |

### 文章管理接口

| 功能 | 旧路由 | 新路由 |
|------|--------|--------|
| 获取列表 | GET `/articles` | GET `/articles` ✅ 不变 |
| 获取详情 | GET `/articles/:id` | GET `/articles/:id` ✅ 不变 |
| 创建文章 | POST `/articles` | POST `/articles` ✅ 不变 |
| 更新文章 | PUT `/articles/:id` ❌ | POST `/articles/:id/update` ✅ |
| 删除文章 | DELETE `/articles/:id` ❌ | POST `/articles/:id/delete` ✅ |
| 点赞文章 | POST `/articles/:id/like` | POST `/articles/:id/like` ✅ 不变 |

---

## 📚 完整 API 列表

### 健康检查（无需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/ping` | 简单健康检查 |
| GET | `/health` | 完整健康检查 |
| GET | `/ready` | 就绪检查 |
| GET | `/live` | 存活检查 |

### 认证接口（公开）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/auth/register` | 用户注册 |
| POST | `/api/v1/auth/login` | 用户登录 |

### 认证接口（需要 Token）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/auth/logout` | 用户登出 |
| GET | `/api/v1/auth/user-info` | 获取当前用户信息 |

### 用户管理（需要 Token）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/users` | 获取用户列表 |
| GET | `/api/v1/users/:id` | 获取指定用户 |
| POST | `/api/v1/users/update` | 更新用户信息 ✨ |
| POST | `/api/v1/users/:id/delete` | 删除用户 ✨ |

### 文章管理

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| GET | `/api/v1/articles` | 否 | 获取文章列表 |
| GET | `/api/v1/articles/:id` | 否 | 获取文章详情 |
| POST | `/api/v1/articles` | 是 | 创建文章 |
| POST | `/api/v1/articles/:id/update` | 是 | 更新文章 ✨ |
| POST | `/api/v1/articles/:id/delete` | 是 | 删除文章 ✨ |
| POST | `/api/v1/articles/:id/like` | 是 | 点赞文章 |

---

## 🎯 为什么只用 GET 和 POST？

### 优势

1. **简化客户端实现** - 有些老旧客户端不支持 PUT/DELETE
2. **防火墙友好** - 某些防火墙只允许 GET/POST
3. **更明确的语义** - 通过路径明确操作类型
4. **统一风格** - 所有修改操作都用 POST

### RESTful vs 实用主义

虽然 RESTful 推荐使用：
- GET - 查询
- POST - 创建
- PUT - 更新
- DELETE - 删除

但在实际项目中，**只用 GET 和 POST** 也是一种常见的实用主义做法。

---

## 💡 使用示例

### 更新用户信息

```bash
# 旧方式（PUT）
curl -X PUT http://localhost:8080/api/v1/users \
  -H "X-Token: YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"new@example.com"}'

# 新方式（POST）
curl -X POST http://localhost:8080/api/v1/users/update \
  -H "X-Token: YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"new@example.com"}'
```

### 删除用户

```bash
# 旧方式（DELETE）
curl -X DELETE http://localhost:8080/api/v1/users/2 \
  -H "X-Token: YOUR_TOKEN"

# 新方式（POST）
curl -X POST http://localhost:8080/api/v1/users/2/delete \
  -H "X-Token: YOUR_TOKEN"
```

### 更新文章

```bash
# 旧方式（PUT）
curl -X PUT http://localhost:8080/api/v1/articles/xxx/update \
  -H "X-Token: YOUR_TOKEN" \
  -d '{"title":"新标题"}'

# 新方式（POST）
curl -X POST http://localhost:8080/api/v1/articles/xxx/update \
  -H "X-Token: YOUR_TOKEN" \
  -d '{"title":"新标题"}'
```

---

## 📖 在 Swagger UI 中使用

访问: http://localhost:8080/swagger/index.html

现在所有接口都只显示 **GET** 和 **POST** 方法：

```
✅ GET 方法（查询）
  - /ping
  - /health
  - /users
  - /users/{id}
  - /articles
  - /articles/{id}

✅ POST 方法（创建、更新、删除）
  - /auth/register
  - /auth/login
  - /auth/logout
  - /users/update
  - /users/{id}/delete
  - /articles
  - /articles/{id}/update
  - /articles/{id}/delete
  - /articles/{id}/like
```

---

## 🔧 路由设计原则

### GET 请求
- ✅ 用于查询数据
- ✅ 不修改服务器状态
- ✅ 可以缓存
- ✅ 幂等性（多次请求结果相同）

### POST 请求
- ✅ 用于创建、更新、删除数据
- ✅ 修改服务器状态
- ✅ 不缓存
- ✅ 通过路径区分操作（/update、/delete）

---

## ✅ 已更新的文件

```
✅ routes/routes.go          路由配置
✅ controllers/user.go       Swagger 注解
✅ controllers/article.go    Swagger 注解
✅ api.http                  测试文件
✅ Swagger 文档已重新生成
```

---

## 🚀 下一步

### 1. 重启项目

```bash
# 重启应用
make run

# 或
docker-compose restart app
```

### 2. 测试新路由

使用 Swagger UI 或 `api.http` 测试更新后的接口。

### 3. 查看 Swagger 文档

访问: http://localhost:8080/swagger/index.html

所有接口都只使用 GET 和 POST 方法。

---

<div align="center">

## ✅ 路由已更新完成！

**只使用 GET 和 POST · 路径清晰 · 语义明确**

### 测试接口

```bash
# 更新用户
POST /api/v1/users/update

# 删除用户  
POST /api/v1/users/2/delete

# 更新文章
POST /api/v1/articles/xxx/update

# 删除文章
POST /api/v1/articles/xxx/delete
```

🎉 **更简洁、更友好！** 🎉

</div>

