package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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

// 处理服务器回应的消息，直接显示到标准输出
func (client *Client) DealRespone() {
	// 一旦client.conn有数据，就直接copy到stdout标准输出上, 永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) Menu() bool {
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

//查询在线用户
func (client *Client) SelectUsers() {
	_, err := client.conn.Write([]byte("who\n"))
	if err != nil {
		fmt.Println("select users err: ", err)
	}
}

//私聊模式
func (client *Client) PrivateChat() {
	var remoteName string
	var msg string
	client.SelectUsers()
	fmt.Println("chat to someone please input remote name:(exit will exit)")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println("chat to someone please input your message:(exit will exit)")
		fmt.Scanln(&msg)

		if len(remoteName) != 0 {
			for msg != "exit" {
				if len(msg) != 0 {
					sendMsg := "to|" + remoteName + "|" + msg + "\n"
					_, err := client.conn.Write([]byte(sendMsg))
					if err != nil {
						fmt.Println("private chat err: ", err)
						break
					}
				}

				fmt.Println("chat to someone please input your message:(exit will exit)")
				msg = ""
				fmt.Scanln(&msg)
			}

			client.SelectUsers()
			fmt.Println("chat to someone please input remote name:(exit will exit)")
			remoteName = ""
			fmt.Scanln(&remoteName)
		}
	}
}

// 公聊
func (client *Client) PublicChat() {
	fmt.Println("chat to world please input your message:(exit will exit)")
	var msg string
	fmt.Scanln(&msg)

	for msg != "exit" {
		if len(msg) != 0 {
			sendMsg := msg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("public chat err: ", err)
				break
			}
		}

		fmt.Println("chat to world please input your message:(exit will exit)")
		msg = ""
		fmt.Scanln(&msg)
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println("please input your new name:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("client.conn.Write err: ", err)
		return false
	}
	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for !client.Menu() {

		}

		// 根据不同模式处理不同业务
		switch client.flag {
		case 1:
			// 公聊
			client.PublicChat()
		case 2:
			// 私聊
			client.PrivateChat()
		case 3:
			// 更新用户名
			client.UpdateName()
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

	// 开启一个goroutine监听服务器回应
	go client.DealRespone()

	fmt.Println("conn to server success")

	// 启动客户端业务
	client.Run()
}
