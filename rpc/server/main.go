package main

import (
	"net/rpc"
	"net"
	"fmt"
	"net/http"
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
	//3. rpc网络
	rpc.HandleHTTP()
	//4. 监听网络
	ln, err := net.Listen("tcp", "127.0.0.1:12306")
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}
	//5. 等待连接
	http.Serve(ln, nil)
}