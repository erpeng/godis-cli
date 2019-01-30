package main

import (
	"fmt"
	"os"

	"github.com/erpeng/godis-cli/redis"
)

func main() {
	client := &redis.Client{
		Addr: "127.0.0.1:6379",
	}
	c, err := client.OpenConnection()
	if err != nil {
		fmt.Printf("connection failed:%v\n", err)
		os.Exit(1)
	}
	client.LoopReader(c)
}
