我现在要完成登录功能
1. 前端代码路径admin-ui/apps/backend/src/views/_core/authentication
2. 后端代码路径modules/system/api/system/login.go

入参:
```
{"username":"superAdmin","password":"NUKUms5zgbOo4lX1wqZ9eQ=="}
```

出参:
```
{
    "requestId": "e1d80ca2dc259e181c2e402ef4b80c5a",
    "path": "/system/login",
    "success": true,
    "message": "OK",
    "code": 0,
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY2VuZSI6InN5c3RlbSIsImRhdGEiOnsiYXBwX2lkIjoic3lzdGVtIiwiZGVwdElkcyI6WzFdLCJpZCI6MSwicm9sZUlkcyI6WzFdLCJ1c2VybmFtZSI6InN1cGVyQWRtaW4ifSwiZXhwaXJlc0F0IjoxNzc0NTAxNzMyfQ.5PJILnjZJGZlTGGcxaYIELq7pruXmJLdbrbUUowRmy0",
        "expire": 1774501732
    },
    "takeUpTime": 100
}
```
3. 需要完成的功能：登录接口，登录页面，登录逻辑，只完成密码登录
4. 去掉不必要的功能：注册、忘记密码等，其他登录相关的功能先不处理
---------------------------------------------------------------


