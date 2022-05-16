package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		Name:       "default name",
		flag:       999,
	}

	// 连接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = conn

	return client
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.talk to the world")
	fmt.Println("2.talk to someone")
	fmt.Println("3.change name")
	fmt.Println("0.exit")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("please input number [0 - 3]")
		return false
	}
}

func (cliet *Client) Run() {
	for cliet.flag != 0 {
		for !cliet.menu() {

		}

		// 根据不同模式处理不同业务
		switch cliet.flag {
		case 1:
			// 公聊
			fmt.Println("choose talk to the world")
		case 2:
			// 私聊
			fmt.Println("choose talk to someone")
		case 3:
			// 更新用户名
			fmt.Println("choose rename")
		}
	}
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "set server ip (default ip 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "set server port (default port 8888)")
}

func main() {
	// 命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)

	if client == nil {
		fmt.Println("conn to server fail")
		return
	}

	fmt.Println("conn to server success")

	// 启动客户端业务
	client.Run()
}
