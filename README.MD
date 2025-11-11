markdown
# 个人博客系统后端
基于 Go + Gin + GORM 开发的个人博客系统后端，提供完整的文章管理、用户认证和评论功能。

## 📋 项目特性

- ✅ 用户注册与登录（JWT认证）
- ✅ 博客文章CRUD操作
- ✅ 文章评论功能
- ✅ Swagger API文档
- ✅ 密码加密存储
- ✅ 权限控制（用户只能操作自己的资源）
- ✅ 完整的错误处理和日志记录

## 🛠 技术栈
- **编程语言**: Go 1.25.3
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0+
- **认证**: JWT
- **文档**: Swagger/OpenAPI 1.0.1
- **密码加密**: bcrypt


## 🚀 **快速开始**

### 环境要求

- Go 1.25.3 或更高版本
- MySQL 8.0 或更高版本
- Git

### 安装步骤

1. **克隆项目**
   ```bash
   git clone https://github.com/xxxx/xxxxx
   ```
2. **安装依赖**
   ```bash
   go mod tidy
   ```
3. **安装 Swag 工具**
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```
4. **生成 Swagger 文档**
   ```bash
   swag init
   ```
5. **数据库配置**

创建 MySQL 数据库：

```sql
CREATE DATABASE blog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```
6. **环境变量配置**

创建 config/config.yaml文件：
```yaml
server:
  port: "8080"
  mode: "debug"

database:
  host: "localhost"
  port: 3306
  user: "root"
  password: ""
  dbname: "Test"
```
或者直接修改 config/config.go 中的默认值。

# 启动应用
1. **开发模式**

```bash
go run main.go
```
2. **构建并运行**

```bash
go build -o GoBlogService
./GoBlogService
```
3. **使用 Air 热重载（开发推荐）**

```bash
# 安装 air
go install github.com/cosmtrek/air@latest

# 运行
air
```
应用启动后，访问以下地址：

应用API: http://localhost:8080

Swagger文档: http://localhost:8080/swagger/index.html

健康检查: http://localhost:8080/health

# 📚 API 文档
## 认证接口
### 1. **用户注册**
**URL**: ```POST /auth/register```

**请求体**:

```json
{
  "username": "testuser",
  "password": "password123",
  "email": "test@example.com"
}
```
**响应**:

```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```
### 2. **用户登录**
**URL**: ```POST /auth/login```

**请求体**:

```json
{
  "username": "testuser",
  "password": "password123"
}
```
**响应**:

```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

# 🗄 数据库配置
## 表结构
**users 表**

```sql
CREATE TABLE users (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  email VARCHAR(100) UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```
**posts 表**

```sql
CREATE TABLE posts (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(200) NOT NULL,
  content TEXT NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```
**comments 表**

```sql
CREATE TABLE comments (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  content TEXT NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  post_id BIGINT UNSIGNED NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);
```
# 🔒 安全特性
## JWT Token 认证

- 密码 bcrypt 加密

- SQL 注入防护（GORM 参数化查询）

- CORS 中间件支持




