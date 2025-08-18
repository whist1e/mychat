# MyChat 项目第一阶段实现规划

## 阶段概述

**目标：** 实现基础功能，建立项目核心架构，为用户提供基础的聊天体验
**时间周期：** 预计4-6周
**优先级：** 高

## 功能模块详细规划

### 1. 用户注册登录系统

#### 1.1 用户注册
- **功能描述：** 新用户创建账户
- **技术要求：**
  - 用户名唯一性验证
  - 密码强度验证（至少8位，包含字母和数字）
  - 邮箱格式验证
- **GORM模型设计：**
  ```go
  // User 用户模型
  type User struct {
      ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
      Username     string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
      Email        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
      PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
      AvatarURL    string    `gorm:"type:varchar(255)" json:"avatar_url"`
      Status       string    `gorm:"type:enum('active','inactive','banned');default:'active'" json:"status"`
      CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
      UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
  }
  
  // TableName 指定表名
  func (User) TableName() string {
      return "users"
  }
  ```
- **API接口：**
  - `POST /api/auth/register` - 用户注册
  - `POST /api/auth/verify-email` - 邮箱验证

#### 1.2 用户登录
- **功能描述：** 已注册用户登录系统
- **技术要求：**
  - JWT token认证
  - 密码加密存储（bcrypt）
  - 登录失败次数限制
- **API接口：**
  - `POST /api/auth/login` - 用户登录
  - `POST /api/auth/logout` - 用户登出
  - `GET /api/auth/profile` - 获取用户信息

#### 1.3 密码重置
- **功能描述：** 忘记密码时的重置功能
- **技术要求：**
  - 邮箱验证码发送
  - 临时token生成
  - 密码重置链接有效期限制
- **API接口：**
  - `POST /api/auth/forgot-password` - 发送重置邮件
  - `POST /api/auth/reset-password` - 重置密码

### 2. 基础聊天功能

#### 2.1 一对一聊天
- **功能描述：** 两个用户之间的私聊
- **技术要求：**
  - WebSocket实时通信
  - 消息持久化存储
  - 消息状态跟踪（已发送、已读、已送达）
- **GORM模型设计：**
  ```go
  // Conversation 对话模型
  type Conversation struct {
      ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
      User1ID   uint      `gorm:"not null" json:"user1_id"`
      User2ID   uint      `gorm:"not null" json:"user2_id"`
      CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
      
      // 关联关系
      User1     User `gorm:"foreignKey:User1ID" json:"user1"`
      User2     User `gorm:"foreignKey:User2ID" json:"user2"`
      Messages  []Message `gorm:"foreignKey:ConversationID" json:"messages"`
  }
  
  // Message 消息模型
  type Message struct {
      ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
      ConversationID uint      `gorm:"not null" json:"conversation_id"`
      SenderID       uint      `gorm:"not null" json:"sender_id"`
      Content        string    `gorm:"type:text;not null" json:"content"`
      MessageType    string    `gorm:"type:enum('text','image','file');default:'text'" json:"message_type"`
      Status         string    `gorm:"type:enum('sent','delivered','read');default:'sent'" json:"status"`
      CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
      
      // 关联关系
      Conversation  Conversation `gorm:"foreignKey:ConversationID" json:"conversation"`
      Sender       User         `gorm:"foreignKey:SenderID" json:"sender"`
  }
  
  // TableName 指定表名
  func (Conversation) TableName() string {
      return "conversations"
  }
  
  func (Message) TableName() string {
      return "messages"
  }
  ```
- **API接口：**
  - `GET /api/conversations` - 获取对话列表
  - `GET /api/conversations/{id}/messages` - 获取对话消息
  - `POST /api/conversations/{id}/messages` - 发送消息
  - `PUT /api/messages/{id}/read` - 标记消息已读

#### 2.2 消息类型支持
- **文本消息：** 基础文本内容
- **图片消息：** 支持常见图片格式（JPG、PNG、GIF）
- **文件消息：** 支持文档、压缩包等文件类型
- **表情包：** 内置表情包系统

#### 2.3 消息历史
- **功能描述：** 查看历史聊天记录
- **技术要求：**
  - 分页加载
  - 消息搜索功能
  - 消息时间显示

### 3. 联系人管理

#### 3.1 好友系统
- **功能描述：** 添加、删除、管理好友关系
- **技术要求：**
  - 好友申请流程
  - 好友状态管理
  - 好友分组功能
- **GORM模型设计：**
  ```go
  // Friendship 好友关系模型
  type Friendship struct {
      ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
      UserID    uint      `gorm:"not null" json:"user_id"`
      FriendID  uint      `gorm:"not null" json:"friend_id"`
      Status    string    `gorm:"type:enum('pending','accepted','rejected','blocked');default:'pending'" json:"status"`
      CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
      UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
      
      // 关联关系
      User       User `gorm:"foreignKey:UserID" json:"user"`
      Friend     User `gorm:"foreignKey:FriendID" json:"friend"`
  }
  
  // FriendGroup 好友分组模型
  type FriendGroup struct {
      ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
      UserID    uint      `gorm:"not null" json:"user_id"`
      GroupName string    `gorm:"type:varchar(50);not null" json:"group_name"`
      CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
      
      // 关联关系
      User    User                `gorm:"foreignKey:UserID" json:"user"`
      Members []FriendGroupMember `gorm:"foreignKey:GroupID" json:"members"`
  }
  
  // FriendGroupMember 好友分组成员模型
  type FriendGroupMember struct {
      ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
      GroupID   uint      `gorm:"not null" json:"group_id"`
      FriendID  uint      `gorm:"not null" json:"friend_id"`
      CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
      
      // 关联关系
      Group  FriendGroup `gorm:"foreignKey:GroupID" json:"group"`
      Friend User        `gorm:"foreignKey:FriendID" json:"friend"`
  }
  
  // TableName 指定表名
  func (Friendship) TableName() string {
      return "friendships"
  }
  
  func (FriendGroup) TableName() string {
      return "friend_groups"
  }
  
  func (FriendGroupMember) TableName() string {
      return "friend_group_members"
  }
  ```

#### 3.2 好友申请
- **功能描述：** 发送和接收好友申请
- **技术要求：**
  - 申请通知
  - 申请状态跟踪
  - 批量处理功能
- **API接口：**
  - `POST /api/friends/request` - 发送好友申请
  - `GET /api/friends/requests` - 获取好友申请列表
  - `PUT /api/friends/requests/{id}` - 处理好友申请

#### 3.3 联系人列表
- **功能描述：** 显示所有好友和在线状态
- **技术要求：**
  - 实时在线状态更新
  - 联系人搜索
  - 联系人排序（在线优先、最近聊天优先）

## 技术实现方案

### 前端技术栈
- **Vue3 Composition API** - 现代化组件开发
- **Vue Router 4** - 路由管理
- **Pinia** - 状态管理（替代Vuex）
- **Element Plus** - UI组件库
- **Axios** - HTTP客户端
- **Socket.io-client** - WebSocket客户端

### 后端技术栈
- **Go 1.21+** - 后端语言
- **Gin** - Web框架
- **GORM** - ORM框架
- **MySQL 8.0** - 主数据库
- **Redis** - 缓存和会话存储
- **JWT** - 身份认证
- **WebSocket** - 实时通信

### 数据库设计原则
- 使用InnoDB存储引擎
- 建立合适的索引
- 外键约束保证数据完整性
- 时间戳字段用于数据追踪

### 完整的GORM模型文件示例
```go
// models/models.go
package models

import (
    "time"
    "gorm.io/gorm"
)

// User 用户模型
type User struct {
    ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Username     string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
    Email        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
    PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
    AvatarURL    string    `gorm:"type:varchar(255)" json:"avatar_url"`
    Status       string    `gorm:"type:enum('active','inactive','banned');default:'active'" json:"status"`
    CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
    
    // 关联关系
    Conversations []Conversation `gorm:"foreignKey:User1ID" json:"conversations_as_user1"`
    Messages      []Message      `gorm:"foreignKey:SenderID" json:"sent_messages"`
    Friendships   []Friendship   `gorm:"foreignKey:UserID" json:"friendships"`
    FriendGroups  []FriendGroup  `gorm:"foreignKey:UserID" json:"friend_groups"`
}

// Conversation 对话模型
type Conversation struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    User1ID   uint      `gorm:"not null" json:"user1_id"`
    User2ID   uint      `gorm:"not null" json:"user2_id"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    
    // 关联关系
    User1     User      `gorm:"foreignKey:User1ID" json:"user1"`
    User2     User      `gorm:"foreignKey:User2ID" json:"user2"`
    Messages  []Message `gorm:"foreignKey:ConversationID" json:"messages"`
}

// Message 消息模型
type Message struct {
    ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    ConversationID uint      `gorm:"not null" json:"conversation_id"`
    SenderID       uint      `gorm:"not null" json:"sender_id"`
    Content        string    `gorm:"type:text;not null" json:"content"`
    MessageType    string    `gorm:"type:enum('text','image','file');default:'text'" json:"message_type"`
    Status         string    `gorm:"type:enum('sent','delivered','read');default:'sent'" json:"status"`
    CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
    
    // 关联关系
    Conversation  Conversation `gorm:"foreignKey:ConversationID" json:"conversation"`
    Sender       User         `gorm:"foreignKey:SenderID" json:"sender"`
}

// Friendship 好友关系模型
type Friendship struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID    uint      `gorm:"not null" json:"user_id"`
    FriendID  uint      `gorm:"not null" json:"friend_id"`
    Status    string    `gorm:"type:enum('pending','accepted','rejected','blocked');default:'pending'" json:"status"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
    
    // 关联关系
    User       User `gorm:"foreignKey:UserID" json:"user"`
    Friend     User `gorm:"foreignKey:FriendID" json:"friend"`
}

// FriendGroup 好友分组模型
type FriendGroup struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID    uint      `gorm:"not null" json:"user_id"`
    GroupName string    `gorm:"type:varchar(50);not null" json:"group_name"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    
    // 关联关系
    User    User                `gorm:"foreignKey:UserID" json:"user"`
    Members []FriendGroupMember `gorm:"foreignKey:GroupID" json:"members"`
}

// FriendGroupMember 好友分组成员模型
type FriendGroupMember struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    GroupID   uint      `gorm:"not null" json:"group_id"`
    FriendID  uint      `gorm:"not null" json:"friend_id"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    
    // 关联关系
    Group  FriendGroup `gorm:"foreignKey:GroupID" json:"group"`
    Friend User        `gorm:"foreignKey:FriendID" json:"friend"`
}

// 表名指定
func (User) TableName() string {
    return "users"
}

func (Conversation) TableName() string {
    return "conversations"
}

func (Message) TableName() string {
    return "messages"
}

func (Friendship) TableName() string {
    return "friendships"
}

func (FriendGroup) TableName() string {
    return "friend_groups"
}

func (FriendGroupMember) TableName() string {
    return "friend_group_members"
}

// 数据库初始化函数
func InitDB(db *gorm.DB) error {
    // 自动迁移数据库表结构
    return db.AutoMigrate(
        &User{},
        &Conversation{},
        &Message{},
        &Friendship{},
        &FriendGroup{},
        &FriendGroupMember{},
    )
}
```

## 开发计划

### 第1周：项目搭建和基础架构
- [ ] 项目目录结构搭建
- [ ] 数据库设计和创建
- [ ] 基础Go项目配置
- [ ] 前端Vue3项目初始化
- [ ] 基础路由和组件结构
- [ ] GORM模型定义和数据库迁移

### 第2周：用户认证系统
- [ ] 用户注册功能
- [ ] 用户登录功能
- [ ] JWT token实现
- [ ] 密码加密和验证
- [ ] 基础用户信息管理

### 第3周：基础聊天功能
- [ ] WebSocket连接建立
- [ ] 一对一聊天实现
- [ ] 消息发送和接收
- [ ] 消息持久化存储
- [ ] 基础UI界面

### 第4周：联系人管理
- [ ] 好友申请系统
- [ ] 好友关系管理
- [ ] 联系人列表显示
- [ ] 在线状态管理
- [ ] 基础搜索功能

### 第5周：功能完善和测试
- [ ] 消息历史记录
- [ ] 文件上传功能
- [ ] 表情包系统
- [ ] 功能测试和bug修复
- [ ] 性能优化

### 第6周：部署和文档
- [ ] 系统部署
- [ ] 用户手册编写
- [ ] API文档完善
- [ ] 部署文档编写
- [ ] 第一阶段总结

## 技术难点和解决方案

### 1. WebSocket连接管理
**难点：** 大量并发连接的管理和断线重连
**解决方案：**
- 使用连接池管理WebSocket连接
- 实现心跳机制检测连接状态
- 断线重连机制和消息队列

### 2. 消息实时性
**难点：** 确保消息的实时推送
**解决方案：**
- 使用Redis pub/sub机制
- 消息队列确保消息不丢失
- 客户端消息确认机制

### 3. 数据一致性
**难点：** 多用户同时操作时的数据一致性
**解决方案：**
- 数据库事务管理
- 乐观锁机制
- 最终一致性保证

## 测试策略

### 单元测试
- 后端API接口测试
- 数据库操作测试
- 前端组件测试

### 集成测试
- WebSocket连接测试
- 用户认证流程测试
- 消息发送接收测试

### 性能测试
- 并发用户测试
- 消息处理性能测试
- 数据库查询性能测试

## 部署要求

### 开发环境
- Go 1.21+
- Node.js 18+
- MySQL 8.0
- Redis 6.0+
- Docker（可选）

### 生产环境
- 负载均衡器
- 数据库主从复制
- Redis集群
- 监控和日志系统

## 成功标准

1. **功能完整性：** 所有计划功能正常实现
2. **性能指标：** 支持100+并发用户
3. **稳定性：** 系统运行稳定，无明显bug
4. **用户体验：** 界面友好，操作流畅
5. **代码质量：** 代码规范，注释完整

## 风险评估

1. **技术风险：** WebSocket实现复杂度较高
2. **时间风险：** 6周时间可能紧张
3. **质量风险：** 快速开发可能影响代码质量

## 下一步计划

第一阶段完成后，将进入第二阶段：
- 群聊功能开发
- 文件传输优化
- 离线消息处理
- 音视频通话基础架构

## 数据库连接配置示例

### Go项目中的数据库配置
```go
// config/database.go
package config

import (
    "fmt"
    "log"
    "os"
    
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    
    "your-project/models"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() error {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    
    if err != nil {
        return fmt.Errorf("failed to connect to database: %v", err)
    }
    
    // 自动迁移数据库表结构
    if err := models.InitDB(db); err != nil {
        return fmt.Errorf("failed to migrate database: %v", err)
    }
    
    DB = db
    log.Println("Database connected successfully")
    return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
    return DB
}
```

### 环境变量配置
```bash
# .env
DB_USER=root
DB_PASSWORD=your_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=mychat
DB_CHARSET=utf8mb4
```

### 数据库连接测试
```go
// 在main.go中调用
func main() {
    if err := config.InitDatabase(); err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    
    // 启动服务器...
}
```
