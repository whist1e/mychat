# 注册登录接口
注册时需要填写：昵称、手机号、密码、验证码，验证码由于申请资质问题跳过

登录根据手机号和密码

密码做了加密后存储，数据库看不到密码（目前忘记密码就登不了）

登录：POST方法发送JSON

{
  "telephone": "17634440689",
  "password": "123456"
}

url：http://localhost:8000/user/login

注册多发送两个字段，url：http://localhost:8000/user/register

还有更新和申请用户信息接口做了一半，可略过

项目结构：
backend/
├── api/                          # API控制器层
│   ├── controller.go             # 基础控制器
│   └── user_info_controller.go  # 用户信息控制器
├── cmd/
│   └── main.go                  # 程序入口
├── configs/
│   └── config.toml              # 配置文件
├── internal/                     # 内部包
│   ├── config/                   # 配置管理
│   ├── dao/                      # 数据访问层
│   ├── dto/                      # 数据传输对象
│   │   ├── request/              # 请求结构体
│   │   └── respond/              # 响应结构体
│   ├── https_server/             # HTTP服务器
│   └── model/                    # 数据模型
├── pkg/                          # 公共包
│   ├── constants/                # 常量定义
│   └── zaplog/                   # 日志工具
├── service/                      # 业务逻辑层
│   └── gorm/                     # GORM服务实现
├── go.mod                        # Go模块依赖
└── logs/                         # 日志文件