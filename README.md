# myblog
go语言写的个人博客，学习用

# 运行环境
go版本：go 1.24.9

# 安装依赖
go mod download

# 启动命令
go run cmd/main.go

# 测试用例

## 用户管理
### 注册用户
POST http://localhost:8080/api/user/register

请求参数：
```
{
    "userName": "user1",
    "password": "123456",
    "nickName": "测试测试1",
    "age": 30,
    "gender": 1,
    "email": "mytest@qq.com"
}
```
正常注册结果：
```
{
    "code": "0",
    "msg": "成功"
}
```
重复注册结果：
```
{
    "code": "2001",
    "msg": "用户已存在"
}
```

### 用户登录
POST http://localhost:8080/api/user/login

请求参数：
```
{
    "userName": "admin",
    "password": "123456"
}
```
登录成功：
```
{
    "code": "0",
    "msg": "成功",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjI2MDk5ODUsInVzZXJJZCI6MSwidXNlck5hbWUiOiJhZG1pbiJ9.H9_w1Cd73HWcVC294QThZVm1qlUgSV0NDMhru-mrckM"
}
```
登录失败场景：
```
{
    "code": "2002",
    "msg": "用户不存在"
}
{
    "code": "1002",
    "msg": "用户名或密码错误"
}
```

## 文章管理
### 获取文章列表
GET http://localhost:8080/api/posts/all?pageNow=1&pageSize=10

```
{
    "code": "0",
    "data": {
        "list": [
            {
                "id": 2,
                "createdAt": "2025-11-05T00:17:23.512603+08:00",
                "updatedAt": "2025-11-05T00:17:23.512603+08:00",
                "title": "这是第一篇文章",
                "content": "文章内容文章内容",
                "userId": 1
            },
            {
                "id": 1,
                "createdAt": "2025-11-04T23:44:31.552413+08:00",
                "updatedAt": "2025-11-05T00:18:31.70961+08:00",
                "title": "这是第一篇文章",
                "content": "文章内容文章内容",
                "userId": 1
            }
        ],
        "pageNow": 1,
        "pageSize": 20,
        "totalCount": 2
    },
    "msg": "成功"
}
```
### 新增文章
POST http://localhost:8080/api/posts/create

```
Header: 
token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjIzNTc0MjIsInVzZXJJZCI6MSwidXNlck5hbWUiOiJhZG1pbiJ9.zpSfa05eCPlTOCGtcS9r1Wy7yKzPEt2VOS6n3SrspjU
```

请求参数：
```
{
    "title": "这是第一篇文章",
    "content": "文章内容文章内容"
}
```
结果：
```
{
    "code": "0",
    "msg": "创建文章成功"
}
```
### 修改文章
POST http://localhost:8080/api/posts/update

```
Header: 
token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjIzNTc0MjIsInVzZXJJZCI6MSwidXNlck5hbWUiOiJhZG1pbiJ9.zpSfa05eCPlTOCGtcS9r1Wy7yKzPEt2VOS6n3SrspjU
```

请求参数：
```
{
    "id": 1,
    "title": "这是第一篇文章-修改后",
    "content": "文章内容文章内容，这里是修改后增加的内容"
}
```
结果：
```
{
    "code": "0",
    "msg": "更新文章成功"
}
```
### 删除文章
POST http://localhost:8080/api/posts/delete

```
Header: 
token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjIzNTc0MjIsInVzZXJJZCI6MSwidXNlck5hbWUiOiJhZG1pbiJ9.zpSfa05eCPlTOCGtcS9r1Wy7yKzPEt2VOS6n3SrspjU
```

请求参数：
```
{
    "id": 2
}
```
结果：
```
{
    "code": "0",
    "msg": "删除文章成功"
}
```

## 评论管理

### 获取文章的评论
GET http://localhost:8080/api/comments/all?postId=1

结果：
```
{
    "code": "0",
    "data": [
        {
            "content": "我是文章1的评论",
            "postId": 1,
            "createdAt": "2025-11-07T21:53:22.902828+08:00",
            "userId": 1,
            "user": {
                "id": 1,
                "createdAt": "0001-01-01T00:00:00Z",
                "updatedAt": "0001-01-01T00:00:00Z",
                "userName": "admin",
                "nickName": "测试测试1",
                "gender": 0,
                "age": 0,
                "email": "",
                "userStatus": 0
            }
        }
    ],
    "msg": "成功"
}
```
### 评论文章
POST http://localhost:8080/api/comments/create

```
token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjI2MTE5NTUsInVzZXJJZCI6MSwidXNlck5hbWUiOiJhZG1pbiJ9.EGIUvvV2aIFnyhWLELT9KhtFAYmVBfQeqs6EmDcthqw
```

请求参数：
```
{
    "postId": 1,
    "content": "我是文章1的评论"
}
```
结果：
```
{
    "code": "0",
    "msg": "评论文章成功"
}
```