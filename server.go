package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	// 在线客户端列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	// 消息广播chan
	Message chan string
}

// 创建server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// 监听Message channel，有需要广播的消息就发送给在线的用户
func (server *Server) ListenMessager() {
	for {
		msg := <-server.Message

		// 把msg发送给全部在线的User
		server.mapLock.Lock()
		for _, cli := range server.OnlineMap {
			cli.C <- msg
		}
		server.mapLock.Unlock()
	}
}

// 广播消息
func (server *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	fmt.Println(sendMsg)
	server.Message <- sendMsg
}

func (server *Server) Handler(conn net.Conn) {
	// 当前连接的业务
	// fmt.Println("建立连接成功..")

	user := NewUser(conn, server)

	user.Online()

	isLive := make(chan bool)

	// 接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err: ", err)
				return
			}

			// 提取用户消息（去除\n）
			msg := string(buf[:n-1])

			// 用户处理消息
			user.DoMessage(msg)

			isLive <- true
		}
	}()

	// 当前handler阻塞
	for {
		select {
		case <-isLive:
			// 当前用户是活跃的，需要重置定时器
			// 不作任何事情，为了激活select，更新下面定时器

		case <-time.After(time.Second * 10):
			// 已经超时
			// 将当前的User强制的关闭

			user.SendMsg("you are offline!!!")

			// 清理资源
			close(user.C)
			conn.Close()
			return
		}
	}
}

// 启动服务器
func (server *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("net.Listen err: ", err)
	}
	// close listen socket
	defer listener.Close()

	// 启动一个goroutine监听需要广播的消息
	go server.ListenMessager()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err: ", err)
			continue
		}
		// do handler
		go server.Handler(conn)
	}
}
