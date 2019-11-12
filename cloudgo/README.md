# cloudgo

## 1. 概述

这是一个web小应用，有以下几个功能：
- 支持静态文件服务
- 支持简单的js访问
- 提交表单，并输出一个表格
- 对/unknow给出开发中的提示，返回码501

## 2. 运行

- 本程序需要iris框架，golang 1.13

- 在终端或者命令行窗口转到main.go所在目录，执行以下命令
```sh
go run main.go
```

- 在浏览器地址栏输入
```sh
# 主页面
http://localhost:port/hello

# 静态资源文件服务
http://localhost:port/static/<任意文件>

# unknown
http://localhost:port/unknown
```
