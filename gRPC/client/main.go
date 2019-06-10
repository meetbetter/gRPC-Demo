package main

import (
	"google.golang.org/grpc"
	"fmt"
	pb "gRPC/myproto"
	"context"
)

func main() {
	//1 连接服务器
	conn, err := grpc.Dial("127.0.0.1:12345",grpc.WithInsecure())//grpc.WithInsecure()指定后才不会报错
	if err != nil {
		fmt.Println("grpc.Dial error:", err)
		return
	}
	defer conn.Close()

	//2 创建客户端句柄
	cln := pb.NewHelloClient(conn)

	//3 调用服务器函数(RPC)
	out,err := cln.GetAdd(context.Background(), &pb.In{Num:10})
	if err != nil {
		fmt.Println("grpc.Dial error:", err)
		return
	}

	//4 打印
	fmt.Println("得到数据:", out.Size)
}
