package main

import (
	"net/rpc"
	"net"
	"fmt"
	"net/rpc/jsonrpc"
)

type Chen struct {
}

//rcp方法
//func (t *T) MethodName(argType T1, replyType *T2) error
func (this *Chen) GetAdd(data int, sum *int) error {

	*sum = data + 100

	return nil
}

func main() {
	//1.对象实例化
	pd := new(Chen)
	//2. rpc注册
	rpc.Register(pd)

	//3. 监听网络
	ln, err := net.Listen("tcp", "127.0.0.1:12306")
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}

	//4. 处理客户端请求
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		go func(conn net.Conn) {
			jsonrpc.ServeConn(conn)
		}(conn)
	}
}