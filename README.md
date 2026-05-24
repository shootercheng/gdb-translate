# 翻译 GDB 官方文档

## 官方文档地址
[https://sourceware.org/gdb/documentation/](https://sourceware.org/gdb/documentation/)

## 下载源码包
1. https://ftp.gnu.org/gnu/gdb
2. 下载 https://ftp.gnu.org/gnu/gdb/gdb-17.2.tar.xz

## 编译html
1. 检查环境,解决依赖,生成Makefile
```bash
$ ./configure
```
2. 根据错误提示安装依赖
```bash
$ sudo apt update
$ sudo apt install libgmp-dev libmpfr-dev texinfo
```
3. 编译html
```bash
$ make html
```
html 文档在 ${workdir}/gdb-17.2/gdb/doc/gdb 中生成了
```bash
~/gdb-17.2/gdb/doc/gdb$ ls -lt|grep ".html"|wc -l
866
```
一共有866个html文件

## 运行翻译程序
1. 安装环境变量加载包
```bash
$ go install github.com/joho/godotenv/cmd/godotenv@latest
$ godotenv

Run a process with an env setup from a .env file

godotenv [-o] [-f ENV_FILE_PATHS] COMMAND_ARGS

ENV_FILE_PATHS: comma separated paths to .env files
COMMAND_ARGS: command and args you want to run

example
  godotenv -f /path/to/something/.env,/another/path/.env fortune
```

2. 运行测试
```bash
$ godotenv -f .env go test ./tests/...
ok  	github.com/shootercheng/gdb-translate/tests/internal/request	26.403s
```

3. 运行主程序
```bash
$ godotenv -f .env go run ./cmd/main/main.go
```
