
[English](/README.md)
# godis-cli
  go版本的redis客户端

## 安装
	go get github.com/erpeng/godis-cli

## 使用
	go run godis-cli.go

## sub 模式
  当执行subscribe后,godis-cli进入非阻塞的subscribe模式
  
	[sub]> subscribe foo
  
  在subscribe模式,只有subscribe/psubscribe/unsubscribe/punsubscribe/ping 这五种命令能够执行.当渠道foo上有消息推送时,Redis server推送消息到godis-cli,godis-cli会打印出消息.使用`exit`退出sub pattern,退出之后会进入正常模式。正常模式下也可以使用`exit`退出


## example
	>bogon:godis-cli didi$ go run godis-cli.go
	> set k1 v1 //正常模式执行set命令
	OK
	> get k1
	v1
	> subscribe foo //进入sub pattern
	subscribe foo 1
	[sub]>subscribe foo1
	subscribe foo1 2
	[sub]> unsubscribe foo1
	unsubscribe foo1 1
	[sub]> ping
	pong
	[sub]> get k1 //sub pattern不能执行get命令
	Redis Error: only (P)SUBSCRIBE / (P)UNSUBSCRIBE / PING / QUIT allowed in this context
	[sub]> exit //退出sub pattern
	exit sub pattern....
	>get k1//现在可以正常执行get命令
	v1
	> exit
