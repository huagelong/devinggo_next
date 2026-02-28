# Encrypted Channels 使用说明

## ⚠️ 重要限制

**Encrypted Channels 仅用于客户端到客户端的端到端加密通信，不支持服务器端推送明文数据。**

---

## 🔐 Encrypted Channels 工作原理

### 客户端到客户端加密流程

1. **订阅阶段**
   ```javascript
   const channel = pusher.subscribe('private-encrypted-secure');
   // 服务器返回 shared_secret (32字节 Base64)
   ```

2. **客户端 A 发送消息**
   ```javascript
   // Pusher.js 自动使用 TweetNaCl 加密
   channel.trigger('client-my-event', {
       message: "敏感信息"
   });
   // 实际发送格式：{ciphertext: "...", nonce: "..."}
   ```

3. **服务器转发**
   - 服务器**不解密**，原封不动转发 `{ciphertext, nonce}`

4. **客户端 B 接收**
   ```javascript
   channel.bind('client-my-event', function(data) {
       // Pusher.js 自动解密，得到原始数据
       console.log(data.message); // "敏感信息"
   });
   ```

---

## ❌ 不支持的场景

### 1. 服务器端推送明文数据

```powershell
# ❌ 错误：推送到 Encrypted Channel
$body = @{
    name = "my-event"
    channels = @("private-encrypted-secure")
    data = '{"message":"Hello"}'
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8070/apps/devinggo-app-id/events" ...
```

**问题**：Pusher.js 期望接收 `{ciphertext, nonce}` 格式，会报错：
```
Unexpected format for encrypted event, expected object with `ciphertext` and `nonce` fields
```

### 2. 服务器端加密后推送

服务器端**没有 shared_secret**（这是客户端独有的），无法加密。即使服务器知道 shared_secret，也违背了端到端加密的初衷。

---

## ✅ 正确使用方法

### 场景1：客户端之间加密聊天

```javascript
// 客户端 A - 发送方
const channel = pusher.subscribe('private-encrypted-chat-room-123');

channel.bind('pusher:subscription_succeeded', () => {
    // 发送加密消息（自动加密）
    channel.trigger('client-message', {
        from: 'Alice',
        text: '这是加密消息',
        timestamp: Date.now()
    });
});

// 客户端 B - 接收方
channel.bind('client-message', (data) => {
    // 自动解密
    console.log(`${data.from}: ${data.text}`);
});
```

### 场景2：服务器端推送到 Private Channel（推荐）

```powershell
# ✅ 正确：使用 Private Channel 接收服务器推送
$body = @{
    name = "payment-notification"
    channels = @("private-user-123")
    data = '{"amount":1000,"status":"success"}'
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8070/apps/devinggo-app-id/events" `
    -Method POST -Body $body -ContentType "application/json" `
    -Headers @{...}  # 需要签名
```

客户端接收：
```javascript
const channel = pusher.subscribe('private-user-123');
channel.bind('payment-notification', (data) => {
    console.log('Payment:', data.amount);  // 正常接收
});
```

---

## 📋 功能对比

| 功能 | Private Channel | Encrypted Channel |
|------|----------------|-------------------|
| **认证** | ✅ 需要 | ✅ 需要 |
| **服务器推送** | ✅ 支持 | ❌ 不支持明文 |
| **Client Events** | ✅ 支持 | ❌ 不支持 |
| **Presence** | ✅ 支持 | ❌ 不支持 |
| **端到端加密** | ❌ 无 | ✅ TweetNaCl 加密 |
| **使用场景** | 一对一通知、权限控制 | 敏感的客户端间通信 |

---

## 🎯 推荐架构

### 混合使用策略

```javascript
// 1. Private Channel - 接收服务器推送
const privateChannel = pusher.subscribe('private-user-123');
privateChannel.bind('server-notification', (data) => {
    // 服务器推送的通知、订单更新等
});

// 2. Encrypted Channel - 客户端加密聊天
const encryptedChannel = pusher.subscribe('private-encrypted-chat-456');
encryptedChannel.bind('client-message', (data) => {
    // 用户之间的加密消息
});

// 3. Presence Channel - 在线状态
const presenceChannel = pusher.subscribe('presence-room-789');
presenceChannel.bind('pusher:member_added', (member) => {
    // 用户上线
});
```

---

## 🚀 测试脚本说明

### test_encrypted_push.ps1（已更新）

**用途**：测试服务器端推送到 **Private Channel**

```powershell
.\test_encrypted_push.ps1
```

- ✅ 推送到 `private-test` 频道
- ✅ 使用 Pusher 签名认证
- ✅ 浏览器可正常接收消息

### Encrypted Channel 测试

**仅限客户端之间**测试：

1. 打开两个浏览器标签页
2. 都订阅 `private-encrypted-secure`
3. 在一个标签页发送 Client Event：
   ```javascript
   channel.trigger('client-test', {message: 'Hello'});
   ```
4. 另一个标签页会接收到**已解密**的消息

---

## 📖 参考文档

- [Pusher Encrypted Channels 官方文档](https://pusher.com/docs/channels/using_channels/encrypted-channels/)
- [完整实现分析](PUSHER_IMPLEMENTATION_ANALYSIS.md)
- [Go 测试套件](../modules/system/pkg/websocket/encrypted_channel_test.go)

---

## 总结

1. ✅ **Encrypted Channels = 客户端到客户端的端到端加密**
2. ❌ **不要用于服务器端推送**（会报格式错误）
3. ✅ **服务器推送请使用 Private Channels**
4. 🔐 **服务器只负责路由加密数据，不解密、不加密**
