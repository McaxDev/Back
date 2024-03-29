# Axolotland API
## 用法
### GET /captcha
#### 成功时
- 响应头：
  - 键：X-Captcha-Id
  - 值：一个字符串，作为提交人机验证时区分不同用户的Id
- 响应体：
  - MIME：image/png
  - 内容：验证码图片
#### 失败时
- 响应体：
  - MIME：application/json
  - 内容：
      ```json
      {
	    "msg": "请求失败",
		"data": null
	  }
      ```
### GET /challenge
#### 成功时
- 响应体：
  - MIME：application/json
  - 内容：
      ```json
      {
	    "msg": "获取挑战值成功",
		"data": {
		  "challenge": "随机字符串，作为挑战值",
		  "data": "挑战值截止时间"
		}
	  }
      ```
#### 失败时
- 响应体：
  - MIME：application/json
  - 内容：
      ```json
      {
	    "msg": "请求失败",
		"data": null
	  }
      ```

### GET /status
#### 查询字符串参数
  - 键：srv
    - 值：互通服 main，原版服 sc，模组服 mod
#### 成功时

#### 失败时


### GET /prompt
#### 成功时

#### 失败时


### GET /variable
#### 成功时

#### 失败时


### POST /register
- 需要挑战应答验证
- 需要人机验证
#### 表单数据
  - 键：username
    - 值：用户名
  - 键：password
    - 值：sha256(挑战值 + 用户密码)
#### 成功时

#### 失败时


### POST /login
- 需要挑战应答验证
- 需要人机验证
#### 表单数据
  - 键：username
    - 值：用户名
  - 键：password
    - 值：sha256(挑战值 + 用户密码)
#### 成功时

#### 失败时


### POST /rcon
- 需要JWT认证
- 需要挑战应答认证
- 需要人机验证
#### 查询字符串参数
  - 键：server
    - 值：互通服 main，原版服 sc，模组服 mod
#### 表单数据
  - 键：command
    - 值：执行的命令
  - 键（可选）：challenge
    - 值：挑战值
  - 键（可选）：password
    - 值：sha256(挑战值 + RCON密码)
	- 用于执行管理员命令
#### 成功时

#### 失败时


### POST /gpt
- 需要JWT认证
- 需要人机验证
#### 查询字符串参数
  - 键：text
    - 值：向GPT问的问题
  - 键（可选）：model
    - 值：使用的模型，3.5 或 4，默认3.5
  - 键（可选）：temperature
    - 值：0~1之间，越低回答越固定，越高回答越随机
#### 成功时

#### 失败时


### POST /source
#### 查询字符串参数
  - 键：title
    - 值：查询的文段标题
#### 成功时

#### 失败时


## 文件功能
### main.go
- 初始化，路由（针对什么路径，采取什么函数处理）注册，命令注册
### handler/
- 处理不同路径的函数实现，写在handler/目录下
### midWare/
- 中间件（如jwt，人机验证，挑战响应认证），写在middleWare/目录下
### config/
- 配置文件和数据库的读取，写在config/目录下
### routine/
- 异步执行的其他协程，写在routine/目录下
### command/
- 接受的子命令，主命令为axoback
  - help 查看命令帮助
  - reload 重新读取配置文件并重启
  - stop 停止服务
  - restart 重启服务
### util/
- 会在其他代码里调用的工具类函数
