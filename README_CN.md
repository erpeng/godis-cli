
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
  
  在subscribe模式,只有subscribe/psubscribe/unsubscribe/punsubscribe/ping 这五种命令能够执行.当渠道foo上有消息推送时,Redis server推送消息到godis-cli,godis-cli会打印出消息
