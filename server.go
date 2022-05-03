package main

import (
	"fmt"
	"net"
	"sync"
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

	user := NewUser(conn)

	// 用户上线，加入OnlineMap
	server.mapLock.Lock()
	server.OnlineMap[user.Name] = user
	server.mapLock.Unlock()

	// 广播当前用户上线消息
	server.BroadCast(user, "online")

	select {}
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
