package main

import (
	"fmt"
	pb "gRPC/myproto" //给package起别名
	"context"
	"net"
	"google.golang.org/grpc"
)

//1. 结构体
type Chen struct {

}

//2. 该结构体实现HelloServer interface的方法
func (this *Chen)GetAdd(ctx context.Context, In *pb.In)(*pb.Out,error)  {
	return &pb.Out{Size:In.Num+100},nil
}

func main() {
	fmt.Println("server runing...")

	//3. 创建网络
	ln, err := net.Listen("tcp", "127.0.0.1:12345")
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}
	defer ln.Close()

	//4. 创建gRPC句柄
	srv := grpc.NewServer()

	//5. 注册server
	pb.RegisterHelloServer(srv, &Chen{})

	//6. 等待网络连接
	err = srv.Serve(ln)
	if err != nil {
		fmt.Println("srv.Serve error:", err)
		return
	}

}