# 文件功能
## main.go
- 初始化，路由（针对什么路径，采取什么函数处理）注册，命令注册
## handler/
- 处理不同路径的函数实现，写在handler/目录下
## midWare/
- 中间件（如jwt，人机验证，挑战响应认证），写在middleWare/目录下
## config/
- 配置文件和数据库的读取，写在config/目录下
## routine/
- 异步执行的其他协程，写在routine/目录下
## command/
- 接受的子命令，主命令为**axoback**
  - help 查看命令帮助
  - reload 重新读取配置文件并重启
  - stop 停止服务
  - restart 重启服务
## util/
- 会在其他代码里调用的工具类函数
