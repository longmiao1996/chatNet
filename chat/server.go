// one sever to more client chat room
//This is chat sever
package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

var ConnMap map[string]net.Conn = make(map[string]net.Conn)  //声明一个集合

//ConnMap := make(map[string]net.Conn)

func main() {
	listen_socket, err := net.Listen("tcp", "127.0.0.1:8000")  //打开监听接口
	if err != nil {
		fmt.Println("server start error")
	}

	defer listen_socket.Close()
	fmt.Println("server is wating ....")

	for {
		conn, err := listen_socket.Accept()  //收到来自客户端发来的消息
		if err != nil {
			fmt.Println("conn fail ...")
		}
		fmt.Println(conn.RemoteAddr(), "connect successed")

		go handle(conn)  //创建线程
	}
}

func handle(conn net.Conn) {
	msg_str := make([]string,0)
	for {
		data := make([]byte, 255)        //创建字节流 （此处同 一对一 通信）
		msg_read, err := conn.Read(data) //声明并将从客户端读取的消息赋给msg_read 和err
		if msg_read == 0 || err != nil {
			continue
		}

		//解析协议
		//say|nickname|message
		if msg_read == 0 {
			msg_str = strings.Split("nick|nick", "|")
		} else {
			msg_str = strings.Split(string(data[0:msg_read]), "|") //将从客户端收到的字节流分段保存到msg_str这个数组中
		}
			nickname := msg_str[1]
			tag := msg_str[0]
			switch tag {
			case "nick": //加入聊天室
				fmt.Println(conn.RemoteAddr(), "-->", nickname) //nick占在数组下标0上，客户端上写的昵称占在数组下标1上
				for k, v := range ConnMap {                     //遍历集合中存储的客户端消息
					if k != nickname {
						//v.Write([]byte("[" + nickname + "]: join..."))
						v.Write([]byte(fmt.Sprintf("[%s]:join chat room", nickname)))
					}
				}
				ConnMap[nickname] = conn
			case "say": //转发消息
				message := strings.TrimRight(msg_str[2], "\n")
				//如果消息以@开头表明是找人单独聊天
				if strings.Contains(message, "@") {
					//@xxx:hello
					new_str := strings.Split(message, ":")
					//nickname := msg_str[1]
					toNickname := strings.TrimLeft(new_str[0], "@")
					//取出单独聊天的人的客户端地址，只给这个地址发消息
					v, ok := ConnMap[toNickname]
					if ok {
						v.Write([]byte(fmt.Sprintf("[%s]:%s", nickname, new_str[1])))
					} else {
						ConnMap[nickname].Write([]byte(fmt.Sprintf("no this user %s please check it", toNickname)))
						log.Println(fmt.Sprintf("no this user --> %s", toNickname))
					}
				} else {
					for k, v := range ConnMap { //k指客户端昵称   v指客户端连接服务器端后的地址

						if k != nickname { //判断是不是给自己发，如果不是
							fmt.Println("Send "+message+" to ", k) //服务器端将消息转发给集合中的每一个客户端
							//v.Write([]byte("[" + nickname + "]: " + message)) //给除了自己的每一个客户端发送自己之前要发送的消息
							v.Write([]byte(fmt.Sprintf("[%s]:%s", nickname, message)))
						}
					}
				}
			case "quit": //退出
				for k, v := range ConnMap { //遍历集合中的客户端昵称
					if k != nickname { //如果昵称不是自己
						//v.Write([]byte("[" + nickname + "]: quit"))  //给除了自己的其他客户端昵称发送退出的消息，并使Write方法阻塞
						v.Write([]byte(fmt.Sprintf("[%s]:quit the chat room", nickname)))
					}
				}

				delete(ConnMap, nickname) //退出聊天室
			}
		}

}