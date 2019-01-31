[中文](/README_CN.md)
# godis-cli
  go version redis-cli

## install
	go get github.com/erpeng/godis-cli

## usage
	go run godis-cli.go

## sub pattern
  when execute  subscribe command ,godis-cli enters subscribe pattern(non-blocking)
  
	[sub]> subscribe foo
  
  in subscribe pattern,only subscribe/psubscribe/unsubscribe/punsubscribe/ping command can execute.when someone pu   blish some message on channel foo,godis-cli can receive message and print it


