package redis

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

var defaultAddr = "127.0.0.1:6379"

//Client redis client
type Client struct {
	Addr     string
	Db       int
	Password string
}

type redisError string

func (err redisError) Error() string { return "Redis Error: " + string(err) }

var doesNotExist = redisError("Key does not exist ")

// reads a bulk reply (i.e $5\r\nhello)
func readBulk(reader *bufio.Reader, head string) ([]byte, error) {
	var err error
	var data []byte

	if head == "" {
		head, err = reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
	}
	switch head[0] {
	case ':':
		data = []byte(strings.TrimSpace(head[1:]))

	case '$':
		size, err := strconv.Atoi(strings.TrimSpace(head[1:]))
		if err != nil {
			return nil, err
		}
		if size == -1 {
			return nil, doesNotExist
		}
		lr := io.LimitReader(reader, int64(size))
		data, err = ioutil.ReadAll(lr)
		if err == nil {
			// read end of line
			_, err = reader.ReadString('\n')
		}
	default:
		return nil, redisError("Expecting Prefix '$' or ':'")
	}

	return data, err
}

func commandBytes(cmd string, args ...string) []byte {
	var cmdbuf bytes.Buffer
	fmt.Fprintf(&cmdbuf, "*%d\r\n$%d\r\n%s\r\n", len(args)+1, len(cmd), cmd)
	for _, s := range args {
		fmt.Fprintf(&cmdbuf, "$%d\r\n%s\r\n", len(s), s)
	}
	return cmdbuf.Bytes()
}

func readResponse(reader *bufio.Reader) (interface{}, error) {

	var line string
	var err error

	//read until the first non-whitespace line
	for {
		line, err = reader.ReadString('\n')

		if len(line) == 0 || err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			break
		}
	}

	if line[0] == '+' {
		return strings.TrimSpace(line[1:]), nil
	}

	if strings.HasPrefix(line, "-ERR ") {
		errmesg := strings.TrimSpace(line[5:])
		return nil, redisError(errmesg)
	}

	if line[0] == ':' {
		n, err := strconv.ParseInt(strings.TrimSpace(line[1:]), 10, 64)
		if err != nil {
			return nil, redisError("Int reply is not a number")
		}
		return n, nil
	}

	if line[0] == '*' {
		size, err := strconv.Atoi(strings.TrimSpace(line[1:]))
		if err != nil {
			return nil, redisError("MultiBulk reply expected a number")
		}
		if size <= 0 {
			return make([][]byte, 0), nil
		}
		res := make([][]byte, size)
		for i := 0; i < size; i++ {
			res[i], err = readBulk(reader, "")
			if err == doesNotExist {
				continue
			}
			if err != nil {
				return nil, err
			}
			// dont read end of line as might not have been bulk
		}
		return res, nil
	}
	return readBulk(reader, line)
}

func (client *Client) rawSend(c net.Conn, cmd []byte) (interface{}, error) {
	_, err := c.Write(cmd)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(c)

	data, err := readResponse(reader)
	if err != nil {
		return nil, err
	}

	return data, nil
}

//OpenConnection open a redis connection
func (client *Client) OpenConnection() (c net.Conn, err error) {

	var addr = defaultAddr

	if client.Addr != "" {
		addr = client.Addr
	}
	c, err = net.Dial("tcp", addr)
	if err != nil {
		return
	}

	//handle authentication here authored by @shxsun
	if client.Password != "" {
		cmd := fmt.Sprintf("AUTH %s\r\n", client.Password)
		_, err = client.rawSend(c, []byte(cmd))
		if err != nil {
			return
		}
	}

	if client.Db != 0 {
		cmd := fmt.Sprintf("SELECT %d\r\n", client.Db)
		_, err = client.rawSend(c, []byte(cmd))
		if err != nil {
			return
		}
	}

	return
}

//EncodeCommand encode a cmd to resp and send to c
func (client *Client) EncodeCommand(c net.Conn, cmd string, args ...string) (interface{}, error) {
	var b []byte
	b = commandBytes(cmd, args...)
	data, err := client.rawSend(c, b)
	return data, err
}

//LoopReader read from stdin and output response
func (client *Client) LoopReader(c net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		args := strings.Split(scanner.Text(), " ")
		if args[0] == "exit" || args[0] == "quit" || args[0] == "q" {
			os.Exit(0)
		}
		response, err := client.EncodeCommand(c, args[0], args[1:]...)
		parseResponse(response, err)
	}

	if scanner.Err() != nil {
		fmt.Printf("%v", scanner.Err())
		os.Exit(2)
	}
}

func parseResponse(response interface{}, err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Print("> ")
	} else {
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
}
