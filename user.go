package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// 创建一个用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}

	go user.ListenMessage()

	return user
}

// 用户上线业务
func (user *User) Online() {
	// 用户上线，加入OnlineMap
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()

	// 广播当前用户上线消息
	user.server.BroadCast(user, "online")
}

// 用户下线业务
func (user *User) Offline() {
	// 用户下线，从OnlineMap中移除
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()

	// 广播当前用户下线消息
	user.server.BroadCast(user, "offline")
}

// 发送消息 测试增、删、改、查等功能
func (user *User) SendMsg(msg string) {
	user.conn.Write([]byte(msg))
}

// 用户处理消息的业务
func (user *User) DoMessage(msg string) {
	// 查询在线用户
	if msg == "who" {
		for _, u := range user.server.OnlineMap {
			sendMsg := "[" + u.Addr + "]" + u.Name + ":" + "online...\n"
			user.SendMsg(sendMsg)
		}
	} else {
		user.server.BroadCast(user, msg)
	}
}

// 监听用户 channel C，一旦有消息，就向该用户的客户端发送
func (user *User) ListenMessage() {
	for {
		msg := <-user.C

		user.conn.Write([]byte(msg + "\n"))
	}
}
