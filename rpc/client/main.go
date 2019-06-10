package main

import (
	"net/rpc"
	"fmt"
)

func main() {
	//1. 连接服务器
	cln, err := rpc.DialHTTP("tcp", "127.0.0.1:12306")
	if err != nil {
		fmt.Println("rpc.Dial error:", err)
		return
	}
	defer cln.Close()

	//2. 调用服务器函数
	var data int
	err = cln.Call("Chen.GetAdd",10, &data)
	if err != nil {
		fmt.Println("cln.Call error:", err)
		return
	}
	//3. 打印输出
	fmt.Println("计算结果为:", data)
}
