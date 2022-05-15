package main

import (
	"net"
	"strings"
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
	user.server.OnlineMap[user.Addr] = user
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

// 通过名字获取user
func (user *User) getUserByName(newName string) *User {
	for _, u := range user.server.OnlineMap {
		if u.Name == newName {
			return u
		}
	}
	return nil
}

// 用户处理消息的业务
func (user *User) DoMessage(msg string) {
	// 查询在线用户
	if msg == "who" {
		for _, u := range user.server.OnlineMap {
			sendMsg := "[" + u.Addr + "]" + u.Name + ":" + "online...\n"
			user.SendMsg(sendMsg)
		}
	} else if len(msg) > 7 && msg[:7] == "rename|" { // 修改用户名
		// 消息格式 rename|李四
		newName := msg[7:]

		if user.getUserByName(newName) != nil {
			user.SendMsg(newName + " has already used...\n")
		} else {
			user.Name = newName
			user.SendMsg("rename success new name: " + newName + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 消息格式 to|李四|消息内容
		splitStrs := strings.Split(msg, "|")
		if len(splitStrs) != 3 {
			user.SendMsg("pattern command is wrong，use pattern like: \"to|zhang san|hello\"\n")
			return
		}
		u := user.getUserByName(splitStrs[1])
		if u != nil {
			u.SendMsg(user.Name + ":" + splitStrs[2] + "\n")
		} else {
			user.SendMsg("no user called " + splitStrs[1] + "\n")
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
