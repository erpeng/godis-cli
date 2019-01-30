package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/erpeng/redis-go"
)

func main() {
	var args []string
	client := &redis.Client{}
	c, err := client.OpenConnection()
	if err != nil {
		fmt.Printf("connection failed:%v\n", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		args = strings.Split(scanner.Text(), " ")
		if args[0] == "exit" || args[0] == "quit" || args[0] == "q" {
			os.Exit(1)
		}
		if args[0] == "subscribe" {
			channel := make(chan string, 1)
			channel <- args[1]
			messages := make(chan redis.Message, 0)
			client.Subscribe(channel, nil, nil, nil, messages)
			// go func() {
			// 	for msg := range messages {
			// 		println("received from:", msg.Channel, " message:", string(msg.Message))
			// 	}
			// }()
			fmt.Print("> ")
		} else {
			client.SubSendCommand(c, args[0], args[1:]...)
			response, err := redis.ReadResponse(bufio.NewReader(c))
			// if strings.HasPrefix(err.Error(), "read tcp") {
			// 	return
			// }
			if err != nil {
				fmt.Printf("%v\n", err)
				fmt.Print("> ")
				goto End
			}
			// db := response.([]byte)
			// for _, b := range db {
			// 	fmt.Printf("%s ", string(b))
			// }
			switch response.(type) {
			case []uint8:
				fmt.Printf("%s", string(response.([]uint8)))
			case string:
				fmt.Printf("%s", response.(string))
			case int64:
				fmt.Printf("%d", response.(int64))
			case [][]uint8:
				for _, b := range response.([][]uint8) {
					fmt.Printf("%s ", string(b))
				}
			default:
				fmt.Printf("%T", response)
			}
			fmt.Print("\n")
			fmt.Print("> ")
		}
		// if args[0] == "publish" {
		// 	err := client.Publish(args[1], []byte(args[2]))
		// 	if err != nil {
		// 		fmt.Printf("- error %v", err)
		// 	}
		// 	fmt.Println("+ ok")
		// }
		// if args[0] == "unsubscribe" {
		// 	channel := make(chan string, 1)
		// 	channel <- args[1]
		// 	go client.Subscribe(nil, channel, nil, nil, nil)
		// 	fmt.Println("+ ok")
		// }

		// if args[0] == "ping" {
		// 	var resp string
		// 	if len(args) == 2 {
		// 		resp, _ = client.Ping(args[1])
		// 	} else {
		// 		resp, _ = client.Ping()
		// 	}
		// 	fmt.Println(resp)
		// }
	End:
	}

	if scanner.Err() != nil {
		errorExit(scanner.Err())
	}
}

func errorExit(err error) {
	fmt.Printf("%v", err)
	os.Exit(2)
}
