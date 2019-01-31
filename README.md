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
  
  in subscribe pattern,only subscribe/psubscribe/unsubscribe/punsubscribe/ping command can execute.when someone pu   blish some message on channel foo,godis-cli can receive message and print it.we can exit sub pattern use `exit`,then we enter normal pattern.Also,we use `exit` exit normal pattern

## example
	>bogon:godis-cli didi$ go run godis-cli.go
	> set k1 v1 //normal pattern
	OK
	> get k1
	v1
	> subscribe foo //enter sub pattern
	subscribe foo 1
	[sub]>subscribe foo1
	subscribe foo1 2
	[sub]> unsubscribe foo1
	unsubscribe foo1 1
	[sub]> ping
	pong
	[sub]> get k1 //we can't execute get in sub pattern
	Redis Error: only (P)SUBSCRIBE / (P)UNSUBSCRIBE / PING / QUIT allowed in this context
	[sub]> exit //exit sub pattern
	exit sub pattern....
	>get k1//now we can execute get in normal pattern
	v1
	> exit
